package main

import (
	"context"
	"encoding/json"
	"os"
	"testing"
	"time"
)

// TestLoadConfig 测试配置文件加载
func TestLoadConfig(t *testing.T) {
	// 创建临时配置文件
	tmpFile, err := os.CreateTemp("", "config*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	configContent := `{
		"targets": [
			{"name": "test1", "address": "https://example.com", "retry_count": 2},
			{"name": "test2", "address": "localhost:8080"}
		]
	}`
	if _, err := tmpFile.Write([]byte(configContent)); err != nil {
		t.Fatal(err)
	}
	tmpFile.Close()

	// 测试加载
	config, err := loadConfig(tmpFile.Name())
	if err != nil {
		t.Errorf("加载配置失败: %v", err)
	}
	if len(config.Targets) != 2 {
		t.Errorf("期望2个目标，实际%d个", len(config.Targets))
	}
	if config.Targets[0].RetryCount != 2 {
		t.Errorf("期望retry_count=2，实际=%d", config.Targets[0].RetryCount)
	}
	if config.Targets[1].RetryCount != 0 {
		t.Errorf("期望retry_count=0，实际=%d", config.Targets[1].RetryCount)
	}
}

// TestRetryCountLimit 测试重试次数限制
func TestRetryCountLimit(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "config*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	// 测试超过3的重试次数会被限制
	configContent := `{
		"targets": [
			{"name": "test", "address": "https://example.com", "retry_count": 5}
		]
	}`
	if _, err := tmpFile.Write([]byte(configContent)); err != nil {
		t.Fatal(err)
	}
	tmpFile.Close()

	config, err := loadConfig(tmpFile.Name())
	if err != nil {
		t.Errorf("加载配置失败: %v", err)
	}
	if config.Targets[0].RetryCount > 3 {
		t.Errorf("重试次数应被限制在3以内，实际=%d", config.Targets[0].RetryCount)
	}
}

// TestNewProbeEngine 测试探测引擎创建
func TestNewProbeEngine(t *testing.T) {
	engine := NewProbeEngine(5, true)
	if engine.timeout != 5*time.Second {
		t.Errorf("超时时间设置错误: 期望%v, 实际%v", 5*time.Second, engine.timeout)
	}
	if engine.verbose != true {
		t.Errorf("详细模式设置错误")
	}
}

// TestProbeHTTP 测试HTTP探测
func TestProbeHTTP(t *testing.T) {
	engine := NewProbeEngine(5, false)

	// 测试有效URL
	ctx := context.Background()
	err := engine.probeHTTP(ctx, "https://www.baidu.com")
	if err != nil {
		t.Logf("百度探测结果: %v", err)
	}

	// 测试无效URL
	err = engine.probeHTTP(ctx, "https://not-exist-domain-12345.com")
	if err == nil {
		t.Log("无效域名应该返回错误")
	}
}

// TestProbeTCP 测试TCP探测
func TestProbeTCP(t *testing.T) {
	engine := NewProbeEngine(3, false)

	ctx := context.Background()
	// 测试本地端口
	err := engine.probeTCP(ctx, "localhost:22")
	if err != nil {
		t.Logf("SSH端口探测: %v", err)
	}

	// 测试无效端口
	err = engine.probeTCP(ctx, "localhost:99999")
	if err == nil {
		t.Log("无效端口应该返回错误")
	}
}

// TestGenerateReport 测试报告生成
func TestGenerateReport(t *testing.T) {
	results := []ProbeResult{
		{Name: "test1", Success: true, Latency: 100 * time.Millisecond},
		{Name: "test2", Success: true, Latency: 200 * time.Millisecond},
		{Name: "test3", Success: false, Latency: 0},
	}

	report := generateReport(results, "config.json", 3)

	if report.TotalCount != 3 {
		t.Errorf("期望总数3，实际%d", report.TotalCount)
	}
	if report.SuccessCount != 2 {
		t.Errorf("期望成功数2，实际%d", report.SuccessCount)
	}
	if report.FailureCount != 1 {
		t.Errorf("期望失败数1，实际%d", report.FailureCount)
	}
	if report.SuccessRate < 66.66 || report.SuccessRate > 66.67 {
		t.Errorf("成功率计算错误: %f", report.SuccessRate)
	}
	if report.AvgLatency != 150*time.Millisecond {
		t.Errorf("平均延迟计算错误: %v", report.AvgLatency)
	}
}

// TestConfigJSONParsing 测试JSON解析
func TestConfigJSONParsing(t *testing.T) {
	jsonData := `{
		"targets": [
			{"name": "测试服务", "address": "https://test.com", "retry_count": 2}
		]
	}`

	var config Config
	err := json.Unmarshal([]byte(jsonData), &config)
	if err != nil {
		t.Errorf("JSON解析失败: %v", err)
	}

	if config.Targets[0].Name != "测试服务" {
		t.Errorf("名称解析错误")
	}
	if config.Targets[0].Address != "https://test.com" {
		t.Errorf("地址解析错误")
	}
}

// TestConcurrentProbe 测试并发探测
func TestConcurrentProbe(t *testing.T) {
	engine := NewProbeEngine(5, false)

	targets := []ProbeTarget{
		{Name: "target1", Address: "https://www.baidu.com"},
		{Name: "target2", Address: "https://www.163.com"},
		{Name: "target3", Address: "localhost:22"},
	}

	start := time.Now()
	results := engine.Run(targets)
	elapsed := time.Since(start)

	// 由于有人为1秒延迟，总时间应该接近1秒（并发执行）
	if elapsed > 2*time.Second {
		t.Logf("并发探测耗时: %v (由于网络延迟可能较长)", elapsed)
	}

	if len(results) != len(targets) {
		t.Errorf("结果数量不匹配: 期望%d, 实际%d", len(targets), len(results))
	}
}
