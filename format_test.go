package caolog

import (
	"strconv"
	"testing"
)

// TestFormatUint 测试FormatUint函数
func TestFormatUint(t *testing.T) {
	tests := []struct {
		name     string
		input    uint64
		expected string
	}{
		{"Zero", 0, "0"},
		{"One", 1, "1"},
		{"Ten", 10, "10"},
		{"Hundred", 100, "100"},
		{"Thousand", 1000, "1000"},
		{"TenThousand", 10000, "10000"},
		{"HundredThousand", 100000, "100000"},
		{"Million", 1000000, "1000000"},
		{"MaxUint32", 4294967295, "4294967295"},
		{"MaxUint64", 18446744073709551615, "18446744073709551615"},
		{"Random1", 123456789, "123456789"},
		{"Random2", 987654321, "987654321"},
		{"Random3", 12345, "12345"},
		{"Random4", 54321, "54321"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatUint(tt.input)
			if result != tt.expected {
				t.Errorf("FormatUint(%d) = %v, want %v", tt.input, result, tt.expected)
			}
			// 验证结果与标准库一致
			expected := strconv.FormatUint(tt.input, 10)
			if result != expected {
				t.Errorf("FormatUint(%d) = %v, but strconv.FormatUint gives %v", tt.input, result, expected)
			}
		})
	}
}

// TestFormatInt 测试FormatInt函数
func TestFormatInt(t *testing.T) {
	tests := []struct {
		name     string
		input    int64
		expected string
	}{
		{"Zero", 0, "0"},
		{"One", 1, "1"},
		{"NegativeOne", -1, "-1"},
		{"Ten", 10, "10"},
		{"NegativeTen", -10, "-10"},
		{"Hundred", 100, "100"},
		{"NegativeHundred", -100, "-100"},
		{"Thousand", 1000, "1000"},
		{"NegativeThousand", -1000, "-1000"},
		{"MaxInt32", 2147483647, "2147483647"},
		{"MinInt32", -2147483648, "-2147483648"},
		{"MaxInt64", 9223372036854775807, "9223372036854775807"},
		{"MinInt64", -9223372036854775808, "-9223372036854775808"},
		{"Random1", 123456789, "123456789"},
		{"Random2", -123456789, "-123456789"},
		{"Random3", 12345, "12345"},
		{"Random4", -54321, "-54321"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatInt(tt.input)
			if result != tt.expected {
				t.Errorf("FormatInt(%d) = %v, want %v", tt.input, result, tt.expected)
			}
			// 验证结果与标准库一致
			expected := strconv.FormatInt(tt.input, 10)
			if result != expected {
				t.Errorf("FormatInt(%d) = %v, but strconv.FormatInt gives %v", tt.input, result, expected)
			}
		})
	}
}

// TestFormatUintEdgeCases 测试FormatUint边界情况
func TestFormatUintEdgeCases(t *testing.T) {
	// 测试0到99的所有值
	for i := uint64(0); i <= 99; i++ {
		result := FormatUint(i)
		expected := strconv.FormatUint(i, 10)
		if result != expected {
			t.Errorf("FormatUint(%d) = %v, want %v", i, result, expected)
		}
	}

	// 测试100的倍数
	for i := uint64(100); i <= 1000; i += 100 {
		result := FormatUint(i)
		expected := strconv.FormatUint(i, 10)
		if result != expected {
			t.Errorf("FormatUint(%d) = %v, want %v", i, result, expected)
		}
	}
}

// TestFormatIntEdgeCases 测试FormatInt边界情况
func TestFormatIntEdgeCases(t *testing.T) {
	// 测试-99到99的所有值
	for i := int64(-99); i <= 99; i++ {
		result := FormatInt(i)
		expected := strconv.FormatInt(i, 10)
		if result != expected {
			t.Errorf("FormatInt(%d) = %v, want %v", i, result, expected)
		}
	}

	// 测试100的倍数
	for i := int64(-1000); i <= 1000; i += 100 {
		result := FormatInt(i)
		expected := strconv.FormatInt(i, 10)
		if result != expected {
			t.Errorf("FormatInt(%d) = %v, want %v", i, result, expected)
		}
	}
}

// TestFormatUintPerformance 测试FormatUint性能（与标准库比较）
func TestFormatUintPerformance(t *testing.T) {
	testValues := []uint64{0, 1, 10, 100, 1000, 10000, 100000, 1000000, 4294967295, 18446744073709551615}

	for _, val := range testValues {
		result := FormatUint(val)
		expected := strconv.FormatUint(val, 10)
		if result != expected {
			t.Errorf("FormatUint(%d) = %v, want %v", val, result, expected)
		}
	}
}

// TestFormatIntPerformance 测试FormatInt性能（与标准库比较）
func TestFormatIntPerformance(t *testing.T) {
	testValues := []int64{0, 1, -1, 10, -10, 100, -100, 1000, -1000, 2147483647, -2147483648, 9223372036854775807, -9223372036854775808}

	for _, val := range testValues {
		result := FormatInt(val)
		expected := strconv.FormatInt(val, 10)
		if result != expected {
			t.Errorf("FormatInt(%d) = %v, want %v", val, result, expected)
		}
	}
}
