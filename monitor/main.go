package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	// 命令行参数解析
	configPath := flag.String("config", "config.json", "配置文件路径")
	timeoutSec := flag.Int("timeout", 3, "单次探测超时时间(秒)")
	verbose := flag.Bool("v", false, "开启详细模式")
	flag.Parse()

	// 加载配置
	config, err := loadConfig(*configPath)
	if err != nil {
		fmt.Printf("加载配置失败: %v\n", err)
		os.Exit(1)
	}

	if len(config.Targets) == 0 {
		fmt.Println("配置文件中没有探测目标")
		os.Exit(1)
	}

	fmt.Printf("开始探测 %d 个目标...\n", len(config.Targets))
	fmt.Printf("超时时间: %d秒\n", *timeoutSec)
	if *verbose {
		fmt.Println("详细模式: 开启")
	} else {
		fmt.Println("详细模式: 关闭 (使用 -v 开启)")
	}
	fmt.Println()

	// 创建探测引擎并运行
	engine := NewProbeEngine(*timeoutSec, *verbose)
	results := engine.Run(config.Targets)

	// 生成报告
	report := generateReport(results, *configPath, *timeoutSec)

	// 打印报告到控制台
	printReport(report)

	// 保存报告到文件
	filename, err := saveReportToFile(report, *configPath, *timeoutSec)
	if err != nil {
		fmt.Printf("保存报告文件失败: %v\n", err)
	} else {
		fmt.Printf("\n报告已保存到: %s\n", filename)
	}
}
