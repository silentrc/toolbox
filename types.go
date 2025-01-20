package toolbox

import (
	"bytes"
	"encoding/json"
	"strconv"
)

type typesUtils struct {
}

// 数据类型类
func (u *utils) NewTypesUtils() *typesUtils {
	return &typesUtils{}
}

// json去掉特殊字符
func (t *typesUtils) TransHtmlJson(data []byte) []byte {
	data = bytes.Replace(data, []byte("\\u0026"), []byte("&"), -1)
	data = bytes.Replace(data, []byte("\\u003c"), []byte("<"), -1)
	data = bytes.Replace(data, []byte("\\u003e"), []byte(">"), -1)
	return data
}

// 统一转为string
func (t *typesUtils) InterfaceToString(value interface{}) string {
	var key string
	if value == nil {
		return key
	}

	switch ft := value.(type) {
	case float64:
		key = strconv.FormatFloat(ft, 'f', -1, 64)
	case float32:
		key = strconv.FormatFloat(float64(ft), 'f', -1, 64)
	case int:
		key = strconv.Itoa(ft)
	case uint, int8, uint8, int16, uint16, int32, uint32, int64, uint64:
		key = strconv.FormatUint(ft.(uint64), 10)
	case string:
		key = value.(string)
	case []byte:
		key = string(value.([]byte))
	default:
		newValue, _ := json.Marshal(value)
		key = string(newValue)
	}
	return key
}
