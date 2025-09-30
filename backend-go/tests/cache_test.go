package tests

import (
	"errors"
	"testing"
	"time"

	"ecommerce-backend/internal/utils"
)

func TestNewCache(t *testing.T) {
	cache := utils.NewCache()
	if cache == nil {
		t.Error("NewCache should return a non-nil cache")
	}
}

func TestCacheSetAndGet(t *testing.T) {
	cache := utils.NewCache()

	cache.Set("key1", "value1", time.Minute)

	value, exists := cache.Get("key1")
	if !exists {
		t.Error("Key should exist after setting")
	}

	if value != "value1" {
		t.Errorf("Expected 'value1', got %v", value)
	}
}

func TestCacheGetNonExistent(t *testing.T) {
	cache := utils.NewCache()

	_, exists := cache.Get("nonexistent")
	if exists {
		t.Error("Non-existent key should not exist")
	}
}

func TestCacheDelete(t *testing.T) {
	cache := utils.NewCache()

	cache.Set("key1", "value1", time.Minute)
	cache.Delete("key1")

	_, exists := cache.Get("key1")
	if exists {
		t.Error("Key should not exist after deletion")
	}
}

func TestCacheClear(t *testing.T) {
	cache := utils.NewCache()

	cache.Set("key1", "value1", time.Minute)
	cache.Set("key2", "value2", time.Minute)
	cache.Clear()

	_, exists1 := cache.Get("key1")
	_, exists2 := cache.Get("key2")

	if exists1 || exists2 {
		t.Error("Cache should be empty after clear")
	}
}

func TestCacheExpiration(t *testing.T) {
	cache := utils.NewCache()

	cache.Set("key1", "value1", 100*time.Millisecond)

	time.Sleep(150 * time.Millisecond)

	_, exists := cache.Get("key1")
	if exists {
		t.Error("Key should be expired")
	}
}

func TestCacheKeys(t *testing.T) {
	cache := utils.NewCache()

	cache.Set("key1", "value1", time.Minute)
	cache.Set("key2", "value2", time.Minute)

	keys := cache.Keys()
	if len(keys) != 2 {
		t.Errorf("Expected 2 keys, got %d", len(keys))
	}
}

func TestCacheSize(t *testing.T) {
	cache := utils.NewCache()

	if cache.Size() != 0 {
		t.Error("Empty cache should have size 0")
	}

	cache.Set("key1", "value1", time.Minute)
	if cache.Size() != 1 {
		t.Error("Cache should have size 1 after adding one item")
	}
}

func TestStatsCache(t *testing.T) {
	cache := utils.NewStatsCache()

	cache.Set("key1", "value1", time.Minute)
	cache.Get("key1")
	cache.Get("key1")
	cache.Get("nonexistent")

	stats := cache.GetStats()
	if stats.Hits != 2 {
		t.Errorf("Expected 2 hits, got %d", stats.Hits)
	}
	if stats.Misses != 1 {
		t.Errorf("Expected 1 miss, got %d", stats.Misses)
	}
}

func TestCacheManager(t *testing.T) {
	manager := utils.NewCacheManager()

	cache1 := manager.GetCache("cache1")
	cache2 := manager.GetCache("cache2")

	if cache1 == cache2 {
		t.Error("Different cache names should return different cache instances")
	}

	cache1Again := manager.GetCache("cache1")
	if cache1 != cache1Again {
		t.Error("Same cache name should return same cache instance")
	}
}

func TestCacheDecorator(t *testing.T) {
	cache := utils.NewStatsCache()
	decorator := utils.NewCacheDecorator(cache, time.Minute)

	value, err := decorator.GetOrSet("key1", func() (interface{}, error) {
		return "value1", nil
	})

	if err != nil {
		t.Errorf("GetOrSet should not return error: %v", err)
	}

	if value != "value1" {
		t.Errorf("Expected 'value1', got %v", value)
	}

	cachedValue, exists := cache.Get("key1")
	if !exists {
		t.Error("Value should be cached")
	}

	if cachedValue != "value1" {
		t.Errorf("Cached value should be 'value1', got %v", cachedValue)
	}
}

func TestCacheDecoratorWithError(t *testing.T) {
	cache := utils.NewStatsCache()
	decorator := utils.NewCacheDecorator(cache, time.Minute)

	_, err := decorator.GetOrSet("key1", func() (interface{}, error) {
		return nil, errors.New("test error")
	})

	if err == nil {
		t.Error("GetOrSet should return error when function fails")
	}

	_, exists := cache.Get("key1")
	if exists {
		t.Error("Cache should not be set when function returns error")
	}
}

func TestCacheInvalidate(t *testing.T) {
	cache := utils.NewStatsCache()
	decorator := utils.NewCacheDecorator(cache, time.Minute)

	decorator.GetOrSet("key1", func() (interface{}, error) {
		return "value1", nil
	})

	decorator.Invalidate("key1")

	_, exists := cache.Get("key1")
	if exists {
		t.Error("Key should not exist after invalidation")
	}
}

func TestCacheInvalidatePattern(t *testing.T) {
	cache := utils.NewStatsCache()
	decorator := utils.NewCacheDecorator(cache, time.Minute)

	decorator.GetOrSet("user:1", func() (interface{}, error) {
		return "user1", nil
	})
	decorator.GetOrSet("user:2", func() (interface{}, error) {
		return "user2", nil
	})
	decorator.GetOrSet("product:1", func() (interface{}, error) {
		return "product1", nil
	})

	decorator.InvalidatePattern("user:*")

	_, exists1 := cache.Get("user:1")
	_, exists2 := cache.Get("user:2")
	_, exists3 := cache.Get("product:1")

	if exists1 || exists2 {
		t.Error("User keys should be invalidated")
	}

	if !exists3 {
		t.Error("Product key should not be invalidated")
	}
}
