package main

import (
	"testing"
)

// 测试用例结构体
type testCase struct {
	name        string  // 测试名称
	electricity float64 // 用电量
	timeStr     string  // 用电时段
	expected    float64 // 期望最终电费
}

// TestCalculateFinalCost 测试最终电费计算
func TestCalculateFinalCost(t *testing.T) {
	testCases := []testCase{
		// 高峰时段测试
		{"高峰时段-第一档", 100.0, "14:00", 100.0 * 0.5 * 1.1},
		{"高峰时段-第一档边界", 200.0, "10:00", 200.0 * 0.5 * 1.1},
		{"高峰时段-第二档", 300.0, "12:00", (200.0*0.5 + 100.0*0.8) * 1.1},
		{"高峰时段-第二档边界", 400.0, "08:00", (200.0*0.5 + 200.0*0.8) * 1.1},
		{"高峰时段-第三档", 500.0, "21:00", (200.0*0.5 + 200.0*0.8 + 100.0*1.2) * 1.1},

		// 低谷时段测试
		{"低谷时段-第一档", 100.0, "23:00", 100.0 * 0.5 * 0.8},
		{"低谷时段-第一档边界", 200.0, "22:00", 200.0 * 0.5 * 0.8},
		{"低谷时段-第二档", 300.0, "05:00", (200.0*0.5 + 100.0*0.8) * 0.8},
		{"低谷时段-第二档边界", 400.0, "07:00", (200.0*0.5 + 200.0*0.8) * 0.8},
		{"低谷时段-第三档", 500.0, "02:00", (200.0*0.5 + 200.0*0.8 + 100.0*1.2) * 0.8},

		// 边界测试
		{"边界-零点低谷", 150.0, "00:00", 150.0 * 0.5 * 0.8},
		{"边界-八点高峰", 250.0, "08:00", (200.0*0.5 + 50.0*0.8) * 1.1},
		{"边界-22点低谷", 350.0, "22:00", (200.0*0.5 + 150.0*0.8) * 0.8},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := calculateFinalCost(tc.electricity, tc.timeStr)
			// 由于浮点数精度问题，允许0.001的误差
			if result < tc.expected-0.001 || result > tc.expected+0.001 {
				t.Errorf("用电量：%.2f，时段：%s，期望：%.4f，实际：%.4f",
					tc.electricity, tc.timeStr, tc.expected, result)
			}
		})
	}
}

// TestCalculateTierCost 单独测试阶梯电费计算
func TestCalculateTierCost(t *testing.T) {
	tests := []struct {
		name        string
		electricity float64
		expected    float64
	}{
		{"第一档-50度", 50.0, 25.0},
		{"第一档-200度", 200.0, 100.0},
		{"第二档-250度", 250.0, 200.0*0.5 + 50.0*0.8},
		{"第二档-400度", 400.0, 200.0*0.5 + 200.0*0.8},
		{"第三档-450度", 450.0, 200.0*0.5 + 200.0*0.8 + 50.0*1.2},
		{"第三档-500度", 500.0, 200.0*0.5 + 200.0*0.8 + 100.0*1.2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calculateTierCost(tt.electricity)
			if result != tt.expected {
				t.Errorf("用电量：%.2f，期望：%.4f，实际：%.4f",
					tt.electricity, tt.expected, result)
			}
		})
	}
}

// TestGetPeakValleyFactor 单独测试峰谷因子获取
func TestGetPeakValleyFactor(t *testing.T) {
	tests := []struct {
		name     string
		timeStr  string
		expected float64
	}{
		{"高峰-8点", "08:00", PeakFactor},
		{"高峰-12点", "12:00", PeakFactor},
		{"高峰-21点", "21:00", PeakFactor},
		{"低谷-22点", "22:00", ValleyFactor},
		{"低谷-23点", "23:00", ValleyFactor},
		{"低谷-0点", "00:00", ValleyFactor},
		{"低谷-7点", "07:00", ValleyFactor},
		{"错误格式-默认高峰", "25:00", PeakFactor},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getPeakValleyFactor(tt.timeStr)
			if result != tt.expected {
				t.Errorf("时段：%s，期望：%.1f，实际：%.1f",
					tt.timeStr, tt.expected, result)
			}
		})
	}
}
