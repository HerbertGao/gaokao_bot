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
		id, err := GenerateID()
		if err != nil {
			t.Errorf("GenerateID() error = %v, want nil", err)
		}
		if id <= 0 {
			t.Errorf("GenerateID() = %v, want positive number", id)
		}
	})

	t.Run("Generate multiple unique IDs", func(t *testing.T) {
		const count = 1000
		ids := make(map[int64]bool)

		for i := 0; i < count; i++ {
			id, err := GenerateID()
			if err != nil {
				t.Errorf("GenerateID() error = %v, want nil", err)
			}
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
					id, err := GenerateID()
					if err != nil {
						t.Errorf("GenerateID() error = %v, want nil", err)
						return
					}
					idsChan <- id
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

func TestGenerateID_Uninitialized(t *testing.T) {
	// 注意：由于 sync.Once 的特性，在同一进程中无法重置 snowflake 初始化状态
	// 此测试仅作文档说明：如果 snowflake 未初始化，GenerateID 应返回错误
	// 实际测试需要在独立的进程中运行
	t.Skip("Cannot test uninitialized scenario in same process due to sync.Once limitation")

	// 如果能重置，预期行为应该是：
	// _, err := GenerateID()
	// if err == nil {
	//     t.Error("GenerateID() should return error when snowflake is not initialized")
	// }
}
