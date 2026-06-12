package main

import "time"

// ProbeTarget 探测目标配置
type ProbeTarget struct {
	Name       string `json:"name"`
	Address    string `json:"address"`
	RetryCount int    `json:"retry_count,omitempty"` // 重试次数，0表示不重试，最大3
}

// Config 配置文件结构
type Config struct {
	Targets []ProbeTarget `json:"targets"`
}

// ProbeResult 单个探测结果
type ProbeResult struct {
	Name       string        `json:"name"`
	Address    string        `json:"address"`
	Success    bool          `json:"success"`
	Latency    time.Duration `json:"latency_ms"`
	ErrorMsg   string        `json:"error_msg,omitempty"`
	RetryCount int           `json:"retry_count"`
	Attempts   int           `json:"attempts"` // 实际尝试次数
}

// Report 最终报告
type Report struct {
	TotalCount      int           `json:"total_count"`
	SuccessCount    int           `json:"success_count"`
	FailureCount    int           `json:"failure_count"`
	SuccessRate     float64       `json:"success_rate"`
	AvgLatency      time.Duration `json:"avg_latency_ms"`
	MaxLatency      time.Duration `json:"max_latency_ms"`
	MinLatency      time.Duration `json:"min_latency_ms"`
	SlowestServices []ProbeResult `json:"slowest_services"`
	AllResults      []ProbeResult `json:"all_results"`
}
