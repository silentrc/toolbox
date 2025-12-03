package toolbox

import (
	"log"
	"sync"
	"sync/atomic"
)

// SafeMap 是一个线程安全的 map 实现，基于 sync.Map
// 提供了计数功能和类型安全的 string key 操作
type SafeMap struct {
	data  sync.Map
	count int64
}

// NewSafeMap 创建一个新的 SafeMap 实例
func NewSafeMap() *SafeMap {
	return &SafeMap{}
}

// Store 存储键值对，如果键已存在则覆盖
func (m *SafeMap) Store(key string, value any) {
	if _, loaded := m.data.Swap(key, value); !loaded {
		atomic.AddInt64(&m.count, 1)
	}
}

// Load 根据键获取值
func (m *SafeMap) Load(key string) (any, bool) {
	return m.data.Load(key)
}

// LoadOrStore 获取键对应的值，如果不存在则存储新值
// 返回实际的值和是否是加载的（true）还是存储的（false）
func (m *SafeMap) LoadOrStore(key string, value any) (actual any, loaded bool) {
	actual, loaded = m.data.LoadOrStore(key, value)
	if !loaded {
		atomic.AddInt64(&m.count, 1)
	}
	return actual, loaded
}

// LoadAndDelete 加载并删除键值对
func (m *SafeMap) LoadAndDelete(key string) (value any, loaded bool) {
	value, loaded = m.data.LoadAndDelete(key)
	if loaded {
		atomic.AddInt64(&m.count, -1)
	}
	return value, loaded
}

// Delete 删除指定键的键值对
func (m *SafeMap) Delete(key string) {
	if _, loaded := m.data.LoadAndDelete(key); loaded {
		atomic.AddInt64(&m.count, -1)
	}
}

// CompareAndSwap 原子性地比较并交换值
func (m *SafeMap) CompareAndSwap(key string, oldValue any, newValue any) bool {
	return m.data.CompareAndSwap(key, oldValue, newValue)
}

// Len 返回 map 中键值对的数量
func (m *SafeMap) Len() int {
	return int(atomic.LoadInt64(&m.count))
}

// Clear 清空所有键值对
func (m *SafeMap) Clear() {
	m.data.Range(func(key, value interface{}) bool {
		m.data.Delete(key)
		return true
	})
	atomic.StoreInt64(&m.count, 0)
}

// Range 遍历所有键值对
// fn 函数返回 false 时停止遍历
func (m *SafeMap) Range(fn func(key string, value any) bool) {
	m.data.Range(func(rawKey, rawValue interface{}) bool {
		// 类型断言确保 key 是 string 类型
		key, ok := rawKey.(string)
		if !ok {
			// 遇到非 string 类型的 key，记录警告并跳过
			log.Printf("WARNING: unexpected key type %T, expected string", rawKey)
			return true // 继续遍历
		}
		// rawValue 本身就是 any 类型，无需类型断言
		return fn(key, rawValue)
	})
}

// Keys 返回所有键的切片
func (m *SafeMap) Keys() []string {
	var keys []string
	m.Range(func(key string, value any) bool {
		keys = append(keys, key)
		return true
	})
	return keys
}

// Values 返回所有值的切片
func (m *SafeMap) Values() []any {
	var values []any
	m.Range(func(key string, value any) bool {
		values = append(values, value)
		return true
	})
	return values
}

// Has 检查键是否存在
func (m *SafeMap) Has(key string) bool {
	_, ok := m.data.Load(key)
	return ok
}
