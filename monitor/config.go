package main

import (
	"encoding/json"
	"fmt"
	"os"
)

// loadConfig 加载配置文件
func loadConfig(configPath string) (*Config, error) {
	// 如果配置文件路径为空，使用默认路径
	if configPath == "" {
		configPath = "config.json"
	}

	// 读取文件
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	// 解析JSON
	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %w", err)
	}

	// 验证重试次数不超过3
	for i, target := range config.Targets {
		if target.RetryCount > 3 {
			config.Targets[i].RetryCount = 3
		}
		if target.RetryCount < 0 {
			config.Targets[i].RetryCount = 0
		}
	}

	return &config, nil
}
