package main

import (
	"fmt"
	"time"
)

// 常量定义：阶梯阈值和单价
const (
	// 阶梯阈值
	Threshold1 = 200.0 // 第一档上限
	Threshold2 = 400.0 // 第二档上限

	// 阶梯单价
	Price1 = 0.5 // 第一档单价：0.5元/度
	Price2 = 0.8 // 第二档单价：0.8元/度
	Price3 = 1.2 // 第三档单价：1.2元/度

	// 峰谷调节因子
	PeakFactor   = 1.1 // 高峰时段：增加10%
	ValleyFactor = 0.8 // 低谷时段：减少20%
)

// 计费规则版本号
const Version = "v1.0.0"

// init 函数：系统初始化
func init() {
	fmt.Printf("计费规则版本号：%s\n", Version)
	fmt.Printf("系统初始化时间：%s\n", time.Now().Format("2006-01-02 15:04:05"))
}

// main 函数：程序入口
func main() {
	runBilling()
}

// runBilling 函数：引导用户输入并计算电费
func runBilling() {
	var electricity float64
	var timeStr string

	// 引导用户输入用电量
	fmt.Print("请输入用电量（度）：")
	fmt.Scanln(&electricity)

	// 引导用户输入用电时段
	fmt.Print("请输入用电时段（格式：15:04）：")
	fmt.Scanln(&timeStr)

	// 计算最终电费
	finalCost := calculateFinalCost(electricity, timeStr)

	// 打印账单明细
	fmt.Println("\n--- 账单明细 ---")
	fmt.Printf("当前用电：%.2f 度\n", electricity)
	fmt.Printf("当前时段：%s 点\n", timeStr)
	fmt.Printf("最终电费：%.2f 元\n", finalCost)
}

// calculateFinalCost 函数：计算最终电费（阶梯电费 × 峰谷调节因子）
func calculateFinalCost(electricity float64, timeStr string) float64 {
	// 1. 计算阶梯电费
	tierCost := calculateTierCost(electricity)

	// 2. 获取峰谷调节因子
	factor := getPeakValleyFactor(timeStr)

	// 3. 返回最终电费
	return tierCost * factor
}

// calculateTierCost 函数：根据阶梯规则计算电费
func calculateTierCost(electricity float64) float64 {
	if electricity <= Threshold1 {
		// 第一档：0-200度
		return electricity * Price1
	} else if electricity <= Threshold2 {
		// 第二档：200-400度（超过200度部分按0.8元）
		firstTierCost := Threshold1 * Price1
		secondTierUsage := electricity - Threshold1
		secondTierCost := secondTierUsage * Price2
		return firstTierCost + secondTierCost
	} else {
		// 第三档：400度以上
		firstTierCost := Threshold1 * Price1
		secondTierCost := (Threshold2 - Threshold1) * Price2
		thirdTierUsage := electricity - Threshold2
		thirdTierCost := thirdTierUsage * Price3
		return firstTierCost + secondTierCost + thirdTierCost
	}
}

// getPeakValleyFactor 函数：根据时间字符串返回峰谷调节因子
func getPeakValleyFactor(timeStr string) float64 {
	// 解析时间字符串
	t, err := time.Parse("15:04", timeStr)
	if err != nil {
		// 如果解析失败，默认按高峰时段处理
		fmt.Println("时段格式错误，默认按高峰时段计费")
		return PeakFactor
	}

	hour := t.Hour()

	// 高峰时段：8:00-22:00（包含8点和22点）
	if hour >= 8 && hour < 22 {
		return PeakFactor
	}
	// 低谷时段：22:00-次日8:00
	return ValleyFactor
}
