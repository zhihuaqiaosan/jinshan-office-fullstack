package main

import (
	"fmt"
	"os"
	"sort"
	"time"
)

// generateReport 生成报告
func generateReport(results []ProbeResult, configPath string, timeoutSec int) *Report {
	if len(results) == 0 {
		return &Report{}
	}

	report := &Report{
		TotalCount: len(results),
		AllResults: results,
		MinLatency: time.Duration(1<<63 - 1),
	}

	var totalLatency time.Duration
	var successResults []ProbeResult

	for _, r := range results {
		if r.Success {
			report.SuccessCount++
			totalLatency += r.Latency
			successResults = append(successResults, r)

			if r.Latency < report.MinLatency {
				report.MinLatency = r.Latency
			}
			if r.Latency > report.MaxLatency {
				report.MaxLatency = r.Latency
			}
		} else {
			report.FailureCount++
		}
	}

	if report.SuccessCount > 0 {
		report.SuccessRate = float64(report.SuccessCount) / float64(report.TotalCount) * 100
		report.AvgLatency = totalLatency / time.Duration(report.SuccessCount)
	} else {
		report.SuccessRate = 0
		report.AvgLatency = 0
		report.MinLatency = 0
	}

	// 找出最慢的3个服务
	sort.Slice(successResults, func(i, j int) bool {
		return successResults[i].Latency > successResults[j].Latency
	})

	maxSlow := 3
	if len(successResults) < maxSlow {
		maxSlow = len(successResults)
	}
	report.SlowestServices = successResults[:maxSlow]

	return report
}

// printReport 打印报告到控制台
func printReport(report *Report) {
	fmt.Println("\n" + stringRepeat("=", 60))
	fmt.Println("                    服务健康探测报告")
	fmt.Println(stringRepeat("=", 60))
	fmt.Printf("总探测数: %d\n", report.TotalCount)
	fmt.Printf("成功数: %d\n", report.SuccessCount)
	fmt.Printf("失败数: %d\n", report.FailureCount)
	fmt.Printf("成功率: %.2f%%\n", report.SuccessRate)

	if report.SuccessCount > 0 {
		fmt.Printf("平均延迟: %v\n", report.AvgLatency)
		fmt.Printf("最大延迟: %v\n", report.MaxLatency)
		fmt.Printf("最小延迟: %v\n", report.MinLatency)
	}

	fmt.Println("\n" + stringRepeat("-", 60))
	fmt.Println("最慢服务 TOP 3:")
	fmt.Println(stringRepeat("-", 60))

	for i, s := range report.SlowestServices {
		fmt.Printf("%d. %s (%s) - %v\n", i+1, s.Name, s.Address, s.Latency)
	}

	fmt.Println("\n" + stringRepeat("-", 60))
	fmt.Println("详细结果:")
	fmt.Println(stringRepeat("-", 60))

	for _, r := range report.AllResults {
		status := "✓ 成功"
		if !r.Success {
			status = "✗ 失败"
		}
		fmt.Printf("[%s] %s (%s)", status, r.Name, r.Address)
		if r.Success {
			fmt.Printf(" - %v", r.Latency)
		} else {
			fmt.Printf(" - %s", r.ErrorMsg)
		}
		if r.RetryCount > 0 {
			fmt.Printf(" (重试%d次, 实际尝试%d次)", r.RetryCount, r.Attempts)
		}
		fmt.Println()
	}

	fmt.Println(stringRepeat("=", 60))
}

// saveReportToFile 保存报告到文件
func saveReportToFile(report *Report, configPath string, timeoutSec int) (string, error) {
	// 生成文件名: monitor-log-20260325164040.log
	filename := fmt.Sprintf("monitor-log-%s.log", time.Now().Format("20060102150405"))

	file, err := os.Create(filename)
	if err != nil {
		return "", fmt.Errorf("创建文件失败: %w", err)
	}
	defer file.Close()

	// 写入报告内容
	fmt.Fprintf(file, "服务健康探测报告\n")
	fmt.Fprintf(file, "生成时间: %s\n", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Fprintf(file, "配置文件: %s\n", configPath)
	fmt.Fprintf(file, "超时设置: %d秒\n", timeoutSec)
	fmt.Fprintf(file, stringRepeat("=", 60)+"\n")
	fmt.Fprintf(file, "总探测数: %d\n", report.TotalCount)
	fmt.Fprintf(file, "成功数: %d\n", report.SuccessCount)
	fmt.Fprintf(file, "失败数: %d\n", report.FailureCount)
	fmt.Fprintf(file, "成功率: %.2f%%\n", report.SuccessRate)

	if report.SuccessCount > 0 {
		fmt.Fprintf(file, "平均延迟: %v\n", report.AvgLatency)
		fmt.Fprintf(file, "最大延迟: %v\n", report.MaxLatency)
		fmt.Fprintf(file, "最小延迟: %v\n", report.MinLatency)
	}

	fmt.Fprintf(file, "\n最慢服务 TOP 3:\n")
	fmt.Fprintf(file, stringRepeat("-", 60)+"\n")
	for i, s := range report.SlowestServices {
		fmt.Fprintf(file, "%d. %s (%s) - %v\n", i+1, s.Name, s.Address, s.Latency)
	}

	fmt.Fprintf(file, "\n详细结果:\n")
	fmt.Fprintf(file, stringRepeat("-", 60)+"\n")
	for _, r := range report.AllResults {
		status := "成功"
		if !r.Success {
			status = "失败"
		}
		fmt.Fprintf(file, "[%s] %s (%s)", status, r.Name, r.Address)
		if r.Success {
			fmt.Fprintf(file, " - %v", r.Latency)
		} else {
			fmt.Fprintf(file, " - %s", r.ErrorMsg)
		}
		if r.RetryCount > 0 {
			fmt.Fprintf(file, " (重试%d次, 实际尝试%d次)", r.RetryCount, r.Attempts)
		}
		fmt.Fprintln(file)
	}

	fmt.Fprintf(file, stringRepeat("=", 60)+"\n")

	return filename, nil
}

// stringRepeat 重复字符串
func stringRepeat(s string, count int) string {
	result := ""
	for i := 0; i < count; i++ {
		result += s
	}
	return result
}
