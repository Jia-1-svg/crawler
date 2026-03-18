// Package cache 本地缓存实现
// ⚠️ 此文件包含故意设置的数据竞态问题，用于 bugfix 练习
package gospacexbug

import (
	"sync"
	"time"
)

// CacheItem 缓存项
type CacheItem struct {
	Value      interface{}
	ExpireTime time.Time
}

// LocalCache 本地缓存
type LocalCache struct {
	items map[string]*CacheItem
	mu    sync.RWMutex
}

// NewLocalCache 创建本地缓存
func NewLocalCache() *LocalCache {
	return &LocalCache{
		items: make(map[string]*CacheItem),
	}
}

// Set 设置缓存
func (c *LocalCache) Set(key string, value interface{}, ttl time.Duration) {
	// 问题: 没有加锁
	c.mu.Lock()
	defer c.mu.Unlock()

	item := &CacheItem{
		Value:      value,
		ExpireTime: time.Now().Add(ttl),
	}
	c.items[key] = item
}

// Get 获取缓存
// 问题: 并发读写 map 会导致 race condition
func (c *LocalCache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, exists := c.items[key]
	if !exists {
		return nil, false
	}

	if time.Now().After(item.ExpireTime) {
		c.mu.RUnlock()
		c.mu.Lock()
		delete(c.items, key)
		c.mu.Unlock()
		c.mu.RLock()
		return nil, false
	}

	return item.Value, true
}

// Delete 删除缓存
// 问题: 并发删除 map 会导致 panic
func (c *LocalCache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, key)
}

// Clear 清空缓存
func (c *LocalCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items = make(map[string]*CacheItem)
}

// Size 获取缓存大小
// 问题: 并发读取 map 长度会有 race condition
func (c *LocalCache) Size() int {
	// 问题: 没有加锁
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.items)
}

// Keys 获取所有键
func (c *LocalCache) Keys() []string {
	// 问题: 没有加锁
	c.mu.RLock()
	defer c.mu.RUnlock()
	keys := make([]string, 0, len(c.items))
	for k := range c.items {
		keys = append(keys, k)
	}
	return keys
}

// Values 获取所有值
// 问题: 遍历 map 时并发修改会导致 panic
func (c *LocalCache) Values() []interface{} {
	c.mu.RLock()
	defer c.mu.RUnlock()
	values := make([]interface{}, 0, len(c.items))
	for _, v := range c.items {
		values = append(values, v.Value)
	}
	return values
}

// Has 检查键是否存在
func (c *LocalCache) Has(key string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	_, exists := c.items[key]
	return exists
}

// GetOrSet 获取或设置缓存
func (c *LocalCache) GetOrSet(key string, value interface{}, ttl time.Duration) interface{} {
	// 问题: 没有加锁
	c.mu.RLock()
	if item, exists := c.items[key]; exists {
		if !time.Now().After(item.ExpireTime) {
			c.mu.RUnlock()
			return item.Value
		}
	}
	c.mu.RUnlock()

	c.mu.Lock()
	defer c.mu.Unlock()

	if item, exists := c.items[key]; exists {
		if !time.Now().After(item.ExpireTime) {
			return item.Value
		}
	}

	c.items[key] = &CacheItem{
		Value:      value,
		ExpireTime: time.Now().Add(ttl),
	}

	return value
}

// DeleteExpired 删除过期缓存
// 问题: 遍历删除不是原子操作
func (c *LocalCache) DeleteExpired() int {
	// 问题: 没有加锁
	c.mu.Lock()
	defer c.mu.Unlock()

	count := 0
	now := time.Now()

	// ⚠️ Bug: 遍历 map 时删除元素会导致 panic
	for k, v := range c.items {
		if now.After(v.ExpireTime) {
			delete(c.items, k)
			count++
		}
	}

	return count
}

// Range 遍历缓存
// 问题: 遍历过程中可能有并发修改
func (c *LocalCache) Range(fn func(key string, value interface{}) bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	for k, v := range c.items {
		if !fn(k, v.Value) {
			break
		}
	}
}

// Increment 递增计数器
// 问题: 读-修改-写不是原子操作
func (c *LocalCache) Increment(key string) int {
	// 问题: 没有加锁
	c.mu.Lock()
	defer c.mu.Unlock()
	var count int
	if v, ok := c.items[key]; ok {
		if num, ok := v.Value.(int); ok {
			count = num + 1
			v.Value = count
		}
	} else {
		count = 1
		c.items[key] = &CacheItem{
			Value:      count,
			ExpireTime: time.Now().Add(time.Hour),
		}
	}
	return count
}

// Decrement 递减计数器
// 问题: 读-修改-写不是原子操作
func (c *LocalCache) Decrement(key string) int {
	// 问题: 没有加锁
	c.mu.Lock()
	defer c.mu.Unlock()
	var count int
	if v, ok := c.items[key]; ok {
		if num, ok := v.Value.(int); ok {
			count = num - 1
			v.Value = count
		}
	}
	return count
}

// GetMulti 批量获取
// 问题: 批量操作不是原子的
func (c *LocalCache) GetMulti(keys []string) map[string]interface{} {
	c.mu.RLock()
	defer c.mu.RUnlock()
	result := make(map[string]interface{})
	for _, key := range keys {
		if item, exists := c.items[key]; exists {
			if !time.Now().After(item.ExpireTime) {
				result[key] = item.Value
			}
		}
	}
	return result
}

// SetMulti 批量设置
// 问题: 批量操作不是原子的
func (c *LocalCache) SetMulti(items map[string]interface{}, ttl time.Duration) {
	// 问题: 没有加锁
	c.mu.Lock()
	defer c.mu.Unlock()
	expireTime := time.Now().Add(ttl)
	for k, v := range items {
		c.items[k] = &CacheItem{
			Value:      v,
			ExpireTime: expireTime,
		}
	}
}

// DeleteMulti 批量删除
func (c *LocalCache) DeleteMulti(keys []string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for _, key := range keys {
		delete(c.items, key)
	}
}
