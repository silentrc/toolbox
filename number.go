package toolbox

import (
	"fmt"
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
	_, err := strconv.Atoi(q)
	if err != nil {
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
	if len(q) <= d {
		sb.WriteString("0.")
		sb.WriteString(strings.Repeat("0", d-len(q))) // 补零
		sb.WriteString(q)
	} else {
		sb.WriteString(q[:len(q)-d])
		sb.WriteByte('.')
		sb.WriteString(q[len(q)-d:])
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
