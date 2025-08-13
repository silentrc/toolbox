package toolbox

import (
	"log"
	"sync"
	"sync/atomic"
)

type SafeMap struct {
	data  sync.Map
	count int64
}

// 新增或覆盖修改
func (m *SafeMap) Store(key string, value any) {
	if _, loaded := m.data.Swap(key, value); !loaded {
		atomic.AddInt64(&m.count, 1)
	}
}

// 条件更新（Compare-and-Swap）
func (m *SafeMap) CompareAndSwap(key string, oldValue any, newValue any) bool {
	return m.data.CompareAndSwap(key, oldValue, newValue)
}

func (m *SafeMap) Update(key string, value any) {
	m.data.Store(key, value) // 直接覆盖，无论键是否存在
}

// 删除
func (m *SafeMap) Delete(key string) {
	if _, loaded := m.data.LoadAndDelete(key); loaded {
		atomic.AddInt64(&m.count, -1)
	}
}

// 获取数量
func (m *SafeMap) Len() int {
	return int(atomic.LoadInt64(&m.count))
}

// 查询
func (m *SafeMap) Load(key string) (any, bool) {
	value, ok := m.data.Load(key)
	if !ok {
		return nil, false
	}
	return value, true
}

func (m *SafeMap) Range(fn func(key string, value any) bool) {
	m.data.Range(func(rawKey, rawValue interface{}) bool {
		// 类型断言确保 key 和 value 的预期类型
		key, ok := rawKey.(string)
		if !ok {
			// 遇到非 string 类型的 key，跳过并记录日志（根据需求处理）
			log.Printf("WARNING: unexpected key type %T, expected string", rawKey)
			return true // 继续遍历
		}
		value, ok := rawValue.(any)
		if !ok {
			// 遇到非 *SOLTokenData 类型的 value，跳过并记录日志
			log.Printf("WARNING: unexpected value type %T for key %s", rawValue, key)
			return true // 继续遍历
		}
		// 调用用户传入的函数
		return fn(key, value)
	})
}
