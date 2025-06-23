package memory

import (
	"sync"
	"time"
)

// cacheValue wraps the cached value and its expiry time.
type cacheValue struct {
	expire time.Time
	value  string
}

// MemoryCache is a simple in-memory implementation of URLCache using sync.Map.
// maxEntries controls the maximum number of keys stored; zero means unlimited.
type MemoryCache struct {
	store      sync.Map
	mu         sync.Mutex
	maxEntries int
	count      int
	ttl        time.Duration
}

func NewMemoryCache(limit int, ttl time.Duration) *MemoryCache {
	return &MemoryCache{maxEntries: limit, ttl: ttl}
}

func (m *MemoryCache) Get(code string) (string, bool) {
	if v, ok := m.store.Load(code); ok {
		cv := v.(cacheValue)
		if cv.expire.IsZero() || time.Now().Before(cv.expire) {
			return cv.value, true
		}
		m.Delete(code)
	}
	return "", false
}

func (m *MemoryCache) Set(key, value string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, loaded := m.store.Load(key); !loaded {
		if m.maxEntries > 0 && m.count >= m.maxEntries {
			var delKey any
			m.store.Range(func(k, _ any) bool {
				delKey = k
				return false
			})
			if delKey != nil {
				m.store.Delete(delKey)
				m.count--
			}
		}
		m.count++
	}

	cv := cacheValue{value: value}
	if m.ttl > 0 {
		cv.expire = time.Now().Add(m.ttl)
	}
	m.store.Store(key, cv)
}

func (m *MemoryCache) Delete(key string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, loaded := m.store.LoadAndDelete(key); loaded {
		if m.count > 0 {
			m.count--
		}
	}
}
