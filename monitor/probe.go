package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"
)

// ProbeEngine 探测引擎
type ProbeEngine struct {
	timeout   time.Duration
	verbose   bool
	client    *http.Client
	waitGroup sync.WaitGroup
	results   chan ProbeResult
}

// NewProbeEngine 创建探测引擎
func NewProbeEngine(timeoutSec int, verbose bool) *ProbeEngine {
	timeout := time.Duration(timeoutSec) * time.Second
	return &ProbeEngine{
		timeout: timeout,
		verbose: verbose,
		client: &http.Client{
			Timeout: timeout,
			Transport: &http.Transport{
				DialContext: (&net.Dialer{
					Timeout:   timeout,
					KeepAlive: 0,
				}).DialContext,
				DisableKeepAlives: true,
			},
		},
		results: make(chan ProbeResult, 100),
	}
}

// Run 运行探测引擎
func (e *ProbeEngine) Run(targets []ProbeTarget) []ProbeResult {
	// 为每个任务添加人为延迟1秒以展示并发性能
	time.Sleep(1 * time.Second)

	for _, target := range targets {
		e.waitGroup.Add(1)
		go e.probeTarget(target)
	}

	// 等待所有探测完成
	go func() {
		e.waitGroup.Wait()
		close(e.results)
	}()

	// 收集结果
	var results []ProbeResult
	for result := range e.results {
		results = append(results, result)
	}

	return results
}

// probeTarget 探测单个目标
func (e *ProbeEngine) probeTarget(target ProbeTarget) {
	defer e.waitGroup.Done()

	maxRetries := target.RetryCount
	if maxRetries < 0 {
		maxRetries = 0
	}
	if maxRetries > 3 {
		maxRetries = 3
	}

	var result ProbeResult

	for attempt := 0; attempt <= maxRetries; attempt++ {
		// 创建带超时的context
		ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
		defer cancel()

		startTime := time.Now()

		// 判断协议类型并探测
		var err error
		if len(target.Address) >= 4 && target.Address[:4] == "http" {
			err = e.probeHTTP(ctx, target.Address)
		} else {
			err = e.probeTCP(ctx, target.Address)
		}

		latency := time.Since(startTime)

		if err == nil {
			result = ProbeResult{
				Name:       target.Name,
				Address:    target.Address,
				Success:    true,
				Latency:    latency,
				RetryCount: maxRetries,
				Attempts:   attempt + 1,
			}
			break
		}

		// 最后一次尝试失败
		if attempt == maxRetries {
			result = ProbeResult{
				Name:       target.Name,
				Address:    target.Address,
				Success:    false,
				Latency:    latency,
				ErrorMsg:   err.Error(),
				RetryCount: maxRetries,
				Attempts:   attempt + 1,
			}
		}
	}

	// 详细模式打印结果
	if e.verbose {
		if result.Success {
			fmt.Printf("[✓] %s (%s) - 成功, 延迟: %v, 尝试次数: %d\n",
				result.Name, result.Address, result.Latency, result.Attempts)
		} else {
			fmt.Printf("[✗] %s (%s) - 失败: %s, 尝试次数: %d\n",
				result.Name, result.Address, result.ErrorMsg, result.Attempts)
		}
	}

	e.results <- result
}

// probeHTTP HTTP探测
func (e *ProbeEngine) probeHTTP(ctx context.Context, url string) error {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return err
	}

	resp, err := e.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 只要状态码在200-399之间都算成功
	if resp.StatusCode >= 200 && resp.StatusCode < 400 {
		return nil
	}
	return fmt.Errorf("HTTP状态码: %d", resp.StatusCode)
}

// probeTCP TCP探测
func (e *ProbeEngine) probeTCP(ctx context.Context, address string) error {
	var dialer net.Dialer
	conn, err := dialer.DialContext(ctx, "tcp", address)
	if err != nil {
		return err
	}
	conn.Close()
	return nil
}
