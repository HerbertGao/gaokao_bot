package util

import (
	"sync"
	"testing"
)

func TestInitSnowflake(t *testing.T) {
	tests := []struct {
		name         string
		datacenterID int64
		machineID    int64
		wantErr      bool
	}{
		{
			name:         "Valid initialization with datacenter 0 machine 0",
			datacenterID: 0,
			machineID:    0,
			wantErr:      false,
		},
		{
			name:         "Valid initialization with datacenter 1 machine 1",
			datacenterID: 1,
			machineID:    1,
			wantErr:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 重置 node 和 once 以便每个测试独立
			// 注意：由于 once 的特性，实际上只会初始化一次
			err := InitSnowflake(tt.datacenterID, tt.machineID)
			if (err != nil) != tt.wantErr {
				t.Errorf("InitSnowflake() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGenerateID(t *testing.T) {
	// 确保 snowflake 已初始化
	err := InitSnowflake(0, 1)
	if err != nil {
		t.Fatalf("Failed to initialize snowflake: %v", err)
	}

	t.Run("Generate single ID", func(t *testing.T) {
		id := GenerateID()
		if id <= 0 {
			t.Errorf("GenerateID() = %v, want positive number", id)
		}
	})

	t.Run("Generate multiple unique IDs", func(t *testing.T) {
		const count = 1000
		ids := make(map[int64]bool)

		for i := 0; i < count; i++ {
			id := GenerateID()
			if ids[id] {
				t.Errorf("GenerateID() generated duplicate ID: %v", id)
			}
			ids[id] = true
		}

		if len(ids) != count {
			t.Errorf("Expected %d unique IDs, got %d", count, len(ids))
		}
	})

	t.Run("Concurrent ID generation", func(t *testing.T) {
		const goroutines = 10
		const idsPerGoroutine = 100
		totalIDs := goroutines * idsPerGoroutine

		idsChan := make(chan int64, totalIDs)
		var wg sync.WaitGroup

		// 并发生成 ID
		for i := 0; i < goroutines; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for j := 0; j < idsPerGoroutine; j++ {
					idsChan <- GenerateID()
				}
			}()
		}

		wg.Wait()
		close(idsChan)

		// 检查所有 ID 是否唯一
		ids := make(map[int64]bool)
		for id := range idsChan {
			if ids[id] {
				t.Errorf("Concurrent GenerateID() generated duplicate ID: %v", id)
			}
			ids[id] = true
		}

		if len(ids) != totalIDs {
			t.Errorf("Expected %d unique IDs in concurrent test, got %d", totalIDs, len(ids))
		}
	})
}

func TestGenerateID_Panic(t *testing.T) {
	// 这个测试无法直接运行，因为我们无法重置 once
	// 仅作为文档说明：如果未初始化 snowflake，GenerateID 会 panic
	t.Skip("Cannot test panic scenario due to sync.Once limitation")
}
