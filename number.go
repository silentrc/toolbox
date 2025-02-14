package toolbox

import (
	"fmt"
	"math/big"
	"strconv"
	"strings"
)

type numberUtils struct {
}

func NewNumberUtils() *numberUtils {
	return &numberUtils{}
}

func (n *numberUtils) DecimalString(q string, d int) string {
	if d < 0 {
		return ""
	}
	q = strings.TrimSpace(q)
	if q == "" {
		return ""
	}

	// 使用 big.Int 验证是否为有效整数
	if _, ok := new(big.Int).SetString(q, 10); !ok {
		return ""
	}

	if d == 0 {
		return q
	}

	var sb strings.Builder
	if q[0] == '-' {
		sb.WriteByte('-')
		q = q[1:]
	}

	qLen := len(q)
	if qLen <= d {
		sb.WriteString("0.")
		sb.WriteString(strings.Repeat("0", d-qLen))
		sb.WriteString(q)
	} else {
		// 处理整数部分和小数部分
		integerPart := q[:qLen-d]
		decimalPart := q[qLen-d:]

		// 移除整数部分前导零（可选）
		if integerPart == "" {
			integerPart = "0"
		} else {
			integerPart = strings.TrimLeft(integerPart, "0")
			if integerPart == "" {
				integerPart = "0"
			}
		}

		sb.WriteString(integerPart)
		sb.WriteByte('.')
		sb.WriteString(decimalPart)
	}

	return sb.String()
}

func (n *numberUtils) DecimalFloat(q string, d int) (string, error) {
	f, err := strconv.ParseFloat(q, 64)
	if err != nil {
		return "", fmt.Errorf("invalid number format in q: %w", err)
	}
	result := f / float64(pow(10, d))
	format := "%." + strconv.Itoa(d) + "f"
	return fmt.Sprintf(format, result), nil
}

func pow(x, n int) int {
	ret := 1 // 结果
	for n != 0 {
		if n%2 != 0 {
			ret = ret * x
		}
		n /= 2
		x = x * x
	}
	return ret
}
