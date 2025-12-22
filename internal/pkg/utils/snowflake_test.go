package utils

import (
	"sync"
	"testing"
)

func TestInitSnowflake(t *testing.T) {
	err := InitSnowflake(1)
	if err != nil {
		t.Fatalf("InitSnowflake() error = %v", err)
	}
}

func TestGenerateID(t *testing.T) {
	InitSnowflake(1)

	id := GenerateID()
	if id == 0 {
		t.Error("GenerateID() returned 0")
	}

	id2 := GenerateID()
	if id == id2 {
		t.Error("GenerateID() returned duplicate IDs")
	}
}

func TestGenerateID_Uniqueness(t *testing.T) {
	InitSnowflake(1)

	ids := make(map[int64]bool)
	count := 10000

	for i := 0; i < count; i++ {
		id := GenerateID()
		if ids[id] {
			t.Fatalf("GenerateID() returned duplicate ID at iteration %d", i)
		}
		ids[id] = true
	}
}

func TestGenerateID_Concurrent(t *testing.T) {
	InitSnowflake(1)

	var wg sync.WaitGroup
	ids := make(chan int64, 10000)
	goroutines := 100
	idsPerGoroutine := 100

	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < idsPerGoroutine; j++ {
				ids <- GenerateID()
			}
		}()
	}

	wg.Wait()
	close(ids)

	seen := make(map[int64]bool)
	for id := range ids {
		if seen[id] {
			t.Fatal("GenerateID() returned duplicate ID in concurrent test")
		}
		seen[id] = true
	}

	if len(seen) != goroutines*idsPerGoroutine {
		t.Errorf("Expected %d unique IDs, got %d", goroutines*idsPerGoroutine, len(seen))
	}
}

func TestGenerateOrderNo(t *testing.T) {
	InitSnowflake(1)

	orderNo := GenerateOrderNo()
	if orderNo == "" {
		t.Error("GenerateOrderNo() returned empty string")
	}

	if len(orderNo) < 10 {
		t.Errorf("GenerateOrderNo() returned too short: %s", orderNo)
	}

	if orderNo[:2] != "FS" {
		t.Errorf("GenerateOrderNo() should start with 'FS', got: %s", orderNo)
	}
}

func TestGenerateTicket(t *testing.T) {
	InitSnowflake(1)

	ticket := GenerateTicket()
	if ticket == "" {
		t.Error("GenerateTicket() returned empty string")
	}

	if ticket[:2] != "TK" {
		t.Errorf("GenerateTicket() should start with 'TK', got: %s", ticket)
	}
}

func BenchmarkGenerateID(b *testing.B) {
	InitSnowflake(1)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		GenerateID()
	}
}

func BenchmarkGenerateOrderNo(b *testing.B) {
	InitSnowflake(1)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		GenerateOrderNo()
	}
}
