package util

import (
	"sync"
	"testing"
)

func TestInitSnowflake(t *testing.T) {
	// 注意：由于 sync.Once 的特性，InitSnowflake 在同一进程中只会执行一次
	// 因此只测试一次调用即可
	err := InitSnowflake(0, 1)
	if err != nil {
		t.Errorf("InitSnowflake() error = %v, want nil", err)
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
