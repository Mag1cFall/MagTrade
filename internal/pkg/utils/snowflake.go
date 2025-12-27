// Snowflake ID 生成器
//
// 本檔案實現 Twitter Snowflake 分散式 ID 演算法
// 結構：41位時間戳 + 10位節點ID + 12位序列號
// 每毫秒每節點可生成 4096 個唯一 ID
package utils

import (
	"fmt"
	"sync"
	"time"
)

// Snowflake ID 生成器
type Snowflake struct {
	mu       sync.Mutex
	epoch    int64 // 起始時間戳
	nodeID   int64 // 節點 ID
	sequence int64 // 序列號
	lastTime int64 // 上次生成時間
}

// 位元配置
const (
	nodeBits     = 10                        // 節點 ID 位元數
	sequenceBits = 12                        // 序列號位元數
	nodeMax      = -1 ^ (-1 << nodeBits)     // 最大節點 ID（1023）
	sequenceMax  = -1 ^ (-1 << sequenceBits) // 最大序列號（4095）
	timeShift    = nodeBits + sequenceBits   // 時間戳左移位元數
	nodeShift    = sequenceBits              // 節點 ID 左移位元數
)

var defaultSnowflake *Snowflake

// InitSnowflake 初始化 Snowflake 生成器
func InitSnowflake(nodeID int64) error {
	if nodeID < 0 || nodeID > nodeMax {
		return fmt.Errorf("node ID must be between 0 and %d", nodeMax)
	}

	defaultSnowflake = &Snowflake{
		epoch:  time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC).UnixMilli(), // 自訂紀元
		nodeID: nodeID,
	}

	return nil
}

// Generate 生成唯一 ID
func (s *Snowflake) Generate() int64 {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now().UnixMilli()

	if now == s.lastTime {
		s.sequence = (s.sequence + 1) & sequenceMax
		if s.sequence == 0 { // 序列號溢出，等待下一毫秒
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

// GenerateID 使用預設生成器生成 ID
func GenerateID() int64 {
	if defaultSnowflake == nil {
		_ = InitSnowflake(1)
	}
	return defaultSnowflake.Generate()
}

// GenerateOrderNo 生成訂單號（FS 前綴）
func GenerateOrderNo() string {
	id := GenerateID()
	return fmt.Sprintf("FS%d", id)
}

// GenerateTicket 生成搶購憑證（TK 前綴）
func GenerateTicket() string {
	id := GenerateID()
	return fmt.Sprintf("TK%d", id)
}
