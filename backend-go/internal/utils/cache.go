package utils

import (
	"sync"
	"time"
)

type CacheItem struct {
	Value     interface{}
	ExpiresAt time.Time
	CreatedAt time.Time
}

type Cache struct {
	items map[string]*CacheItem
	mutex sync.RWMutex
}

func NewCache() *Cache {
	cache := &Cache{
		items: make(map[string]*CacheItem),
	}

	go cache.cleanup()
	return cache
}

func (c *Cache) Set(key string, value interface{}, ttl time.Duration) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.items[key] = &CacheItem{
		Value:     value,
		ExpiresAt: time.Now().Add(ttl),
		CreatedAt: time.Now(),
	}
}

func (c *Cache) Get(key string) (interface{}, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	item, exists := c.items[key]
	if !exists {
		return nil, false
	}

	if time.Now().After(item.ExpiresAt) {
		delete(c.items, key)
		return nil, false
	}

	return item.Value, true
}

func (c *Cache) Delete(key string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	delete(c.items, key)
}

func (c *Cache) Clear() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.items = make(map[string]*CacheItem)
}

func (c *Cache) Keys() []string {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	keys := make([]string, 0, len(c.items))
	for key := range c.items {
		keys = append(keys, key)
	}
	return keys
}

func (c *Cache) Size() int {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	return len(c.items)
}

func (c *Cache) cleanup() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		c.mutex.Lock()
		now := time.Now()
		for key, item := range c.items {
			if now.After(item.ExpiresAt) {
				delete(c.items, key)
			}
		}
		c.mutex.Unlock()
	}
}

type CacheStats struct {
	Size        int
	HitRate     float64
	MissRate    float64
	TotalHits   int64
	TotalMisses int64
	Hits        int64
	Misses      int64
}

type StatsCache struct {
	*Cache
	hits   int64
	misses int64
	mutex  sync.RWMutex
}

func NewStatsCache() *StatsCache {
	return &StatsCache{
		Cache: NewCache(),
	}
}

func (sc *StatsCache) Get(key string) (interface{}, bool) {
	value, exists := sc.Cache.Get(key)

	sc.mutex.Lock()
	if exists {
		sc.hits++
	} else {
		sc.misses++
	}
	sc.mutex.Unlock()

	return value, exists
}

func (sc *StatsCache) GetStats() CacheStats {
	sc.mutex.RLock()
	defer sc.mutex.RUnlock()

	total := sc.hits + sc.misses
	var hitRate, missRate float64

	if total > 0 {
		hitRate = float64(sc.hits) / float64(total) * 100
		missRate = float64(sc.misses) / float64(total) * 100
	}

	return CacheStats{
		Size:        sc.Size(),
		HitRate:     hitRate,
		MissRate:    missRate,
		TotalHits:   sc.hits,
		TotalMisses: sc.misses,
		Hits:        sc.hits,
		Misses:      sc.misses,
	}
}

func (sc *StatsCache) ResetStats() {
	sc.mutex.Lock()
	defer sc.mutex.Unlock()

	sc.hits = 0
	sc.misses = 0
}

type CacheManager struct {
	caches map[string]*StatsCache
	mutex  sync.RWMutex
}

func NewCacheManager() *CacheManager {
	return &CacheManager{
		caches: make(map[string]*StatsCache),
	}
}

func (cm *CacheManager) GetCache(name string) *StatsCache {
	cm.mutex.RLock()
	cache, exists := cm.caches[name]
	cm.mutex.RUnlock()

	if !exists {
		cm.mutex.Lock()
		cache = NewStatsCache()
		cm.caches[name] = cache
		cm.mutex.Unlock()
	}

	return cache
}

func (cm *CacheManager) DeleteCache(name string) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	delete(cm.caches, name)
}

func (cm *CacheManager) ListCaches() []string {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()

	names := make([]string, 0, len(cm.caches))
	for name := range cm.caches {
		names = append(names, name)
	}
	return names
}

func (cm *CacheManager) GetStats() map[string]CacheStats {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()

	stats := make(map[string]CacheStats)
	for name, cache := range cm.caches {
		stats[name] = cache.GetStats()
	}
	return stats
}

func (cm *CacheManager) ClearAll() {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	for _, cache := range cm.caches {
		cache.Clear()
	}
}

type CacheDecorator struct {
	cache *StatsCache
	ttl   time.Duration
}

func NewCacheDecorator(cache *StatsCache, ttl time.Duration) *CacheDecorator {
	return &CacheDecorator{
		cache: cache,
		ttl:   ttl,
	}
}

func (cd *CacheDecorator) GetOrSet(key string, fn func() (interface{}, error)) (interface{}, error) {
	if value, exists := cd.cache.Get(key); exists {
		return value, nil
	}

	value, err := fn()
	if err != nil {
		return nil, err
	}

	cd.cache.Set(key, value, cd.ttl)
	return value, nil
}

func (cd *CacheDecorator) Invalidate(key string) {
	cd.cache.Delete(key)
}

func (cd *CacheDecorator) InvalidatePattern(pattern string) {
	keys := cd.cache.Keys()
	for _, key := range keys {
		if matchPattern(key, pattern) {
			cd.cache.Delete(key)
		}
	}
}

func matchPattern(str, pattern string) bool {
	if pattern == "*" {
		return true
	}

	if len(pattern) == 0 {
		return len(str) == 0
	}

	if pattern[0] == '*' {
		return matchPattern(str, pattern[1:]) ||
			(len(str) > 0 && matchPattern(str[1:], pattern))
	}

	if len(str) == 0 {
		return false
	}

	if pattern[0] == '?' || pattern[0] == str[0] {
		return matchPattern(str[1:], pattern[1:])
	}

	return false
}

var globalCacheManager = NewCacheManager()

func GetCache(name string) *StatsCache {
	return globalCacheManager.GetCache(name)
}

func CacheGetOrSet(cacheName, key string, ttl time.Duration, fn func() (interface{}, error)) (interface{}, error) {
	cache := GetCache(cacheName)
	decorator := NewCacheDecorator(cache, ttl)
	return decorator.GetOrSet(key, fn)
}

func CacheInvalidate(cacheName, key string) {
	cache := GetCache(cacheName)
	cache.Delete(key)
}

func CacheInvalidatePattern(cacheName, pattern string) {
	cache := GetCache(cacheName)
	decorator := NewCacheDecorator(cache, 0)
	decorator.InvalidatePattern(pattern)
}

func GetCacheStats() map[string]CacheStats {
	return globalCacheManager.GetStats()
}
