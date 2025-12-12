package util

import (
	"fmt"
	"sync"

	"github.com/bwmarrin/snowflake"
)

var (
	node *snowflake.Node
	once sync.Once
)

// InitSnowflake 初始化 Snowflake ID 生成器
func InitSnowflake(datacenterID, machineID int64) error {
	var err error
	once.Do(func() {
		// 组合 datacenter 和 machine ID
		nodeID := (datacenterID << 5) | machineID
		node, err = snowflake.NewNode(nodeID)
	})
	return err
}

// GenerateID 生成 Snowflake ID
func GenerateID() (int64, error) {
	if node == nil {
		return 0, fmt.Errorf("snowflake not initialized")
	}
	return node.Generate().Int64(), nil
}