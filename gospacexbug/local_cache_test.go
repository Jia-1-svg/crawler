package gospacexbug

import (
	"sync"
	"testing"
	"time"
)

func TestLocalCache_ConcurrentAccess(t *testing.T) {
	cache := NewLocalCache()
	var wg sync.WaitGroup

	// 并发写入
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			key := string(rune('a' + n%26))
			cache.Set(key, n, time.Hour)
		}(i)
	}

	// 并发读取
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			key := string(rune('a' + n%26))
			cache.Get(key)
		}(i)
	}

	// 并发删除
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			key := string(rune('a' + n%26))
			cache.Delete(key)
		}(i)
	}

	wg.Wait()
}

func TestLocalCache_RaceDetection(t *testing.T) {
	cache := NewLocalCache()
	var wg sync.WaitGroup

	// 并发 Set
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			cache.Set("key", "value", time.Hour)
		}()
	}

	// 并发 Get
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			cache.Get("key")
		}()
	}

	// 并发 Delete
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			cache.Delete("key")
		}()
	}

	wg.Wait()
}

func TestLocalCache_Increment(t *testing.T) {
	cache := NewLocalCache()
	var wg sync.WaitGroup

	// 并发递增 100 次
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			cache.Increment("counter")
		}()
	}

	wg.Wait()

	// 验证结果
	val, ok := cache.Get("counter")
	if !ok {
		t.Fatal("counter not found")
	}
	if val != 100 {
		t.Errorf("Expected counter = 100, got %v", val)
	}
}

func TestLocalCache_GetOrSet(t *testing.T) {
	cache := NewLocalCache()
	var wg sync.WaitGroup

	// 并发 GetOrSet
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			cache.GetOrSet("key", "value", time.Hour)
		}()
	}

	wg.Wait()

	val, ok := cache.Get("key")
	if !ok {
		t.Fatal("key not found")
	}
	if val != "value" {
		t.Errorf("Expected value, got %v", val)
	}
}

func TestLocalCache_MultiOperations(t *testing.T) {
	cache := NewLocalCache()
	var wg sync.WaitGroup

	// 并发 SetMulti
	wg.Add(1)
	go func() {
		defer wg.Done()
		items := map[string]interface{}{
			"a": 1,
			"b": 2,
			"c": 3,
		}
		cache.SetMulti(items, time.Hour)
	}()

	// 并发 GetMulti
	wg.Add(1)
	go func() {
		defer wg.Done()
		cache.GetMulti([]string{"a", "b", "c"})
	}()

	// 并发 DeleteMulti
	wg.Add(1)
	go func() {
		defer wg.Done()
		cache.DeleteMulti([]string{"d", "e", "f"})
	}()

	wg.Wait()
}

func TestLocalCache_Size(t *testing.T) {
	cache := NewLocalCache()
	var wg sync.WaitGroup

	// 并发写入
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			key := string(rune('a' + n%26))
			cache.Set(key, n, time.Hour)
		}(i)
	}

	// 并发读取 Size
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			cache.Size()
		}()
	}

	wg.Wait()
}
