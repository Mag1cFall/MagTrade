package utils

import (
	"fmt"
	"sync"
	"time"
)

type Snowflake struct {
	mu       sync.Mutex
	epoch    int64
	nodeID   int64
	sequence int64
	lastTime int64
}

const (
	nodeBits     = 10
	sequenceBits = 12
	nodeMax      = -1 ^ (-1 << nodeBits)
	sequenceMax  = -1 ^ (-1 << sequenceBits)
	timeShift    = nodeBits + sequenceBits
	nodeShift    = sequenceBits
)

var defaultSnowflake *Snowflake

func InitSnowflake(nodeID int64) error {
	if nodeID < 0 || nodeID > nodeMax {
		return fmt.Errorf("node ID must be between 0 and %d", nodeMax)
	}

	defaultSnowflake = &Snowflake{
		epoch:  time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC).UnixMilli(),
		nodeID: nodeID,
	}

	return nil
}

func (s *Snowflake) Generate() int64 {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now().UnixMilli()

	if now == s.lastTime {
		s.sequence = (s.sequence + 1) & sequenceMax
		if s.sequence == 0 {
			for now <= s.lastTime {
				now = time.Now().UnixMilli()
			}
		}
	} else {
		s.sequence = 0
	}

	s.lastTime = now

	return ((now - s.epoch) << timeShift) | (s.nodeID << nodeShift) | s.sequence
}

func GenerateID() int64 {
	if defaultSnowflake == nil {
		_ = InitSnowflake(1)
	}
	return defaultSnowflake.Generate()
}

func GenerateOrderNo() string {
	id := GenerateID()
	return fmt.Sprintf("FS%d", id)
}

func GenerateTicket() string {
	id := GenerateID()
	return fmt.Sprintf("TK%d", id)
}
