package cache

import (
	"log"
	"reflect"
	"sync"
	"time"
)

const garbageCollectionInterval = 10 * time.Second

// valueMemory is used in memory cache
type valueMemory struct {
	Type       ValueType
	Value      any
	Expiration int       // le 0 means no expiration
	LastAccess time.Time // The last access time, used to determine whether the value is expired
}

type MemoryCache struct {
	MaxEntries int // The maximum number of entries to be stored in the cache
	DefaultTTL int // The default time to live for the cache

	MemoryCacheMap    map[string]valueMemory // The memory cache map
	MemoryCacheRWLock *sync.RWMutex          // The read-write lock for the memory cache
	stopGC            chan struct{}          // Channel to stop the garbage collection goroutine
}

// getValueType returns the value type
func getValueType(value interface{}) ValueType {
	if value == nil {
		return ValueUnknown
	}

	switch reflect.TypeOf(value).Kind() {
	case reflect.Int:
		return ValueInt
	case reflect.String:
		return ValueString
	case reflect.Float64:
		return ValueFloat
	case reflect.Bool:
		return ValueBool
	case reflect.Slice, reflect.Array:
		return ValueArray
	case reflect.Map:
		return ValueMap
	default:
		return ValueUnknown
	}
}

// NewMemoryCache creates a new memory cache
func NewMemoryCache(maxEntries, defaultTTL int) *MemoryCache {
	instance := &MemoryCache{
		MaxEntries:        maxEntries,
		DefaultTTL:        defaultTTL,
		MemoryCacheMap:    make(map[string]valueMemory),
		MemoryCacheRWLock: new(sync.RWMutex),
		stopGC:            make(chan struct{}),
	}

	// Start the garbage collection goroutine
	go instance.garbageCollection()

	return instance
}

// get returns the value from the memory cache
func (c *MemoryCache) get(key string) (valueMemory, bool) {
	c.MemoryCacheRWLock.RLock()
	defer c.MemoryCacheRWLock.RUnlock()

	value, ok := c.MemoryCacheMap[key]
	if ok {
		if value.Expiration > 0 && time.Now().After(value.LastAccess.Add(time.Duration(value.Expiration)*time.Second)) {
			// The value is expired
			return valueMemory{}, false
		}
	}
	return value, ok
}

// GetInt returns the integer value from the memory cache
func (c *MemoryCache) GetInt(key string) (int, error) {
	value, ok := c.get(key)
	if !ok {
		return 0, new(ErrCacheMiss)
	}
	if value.Type != ValueInt {
		log.Printf("[Warning] The value type of the key [%s] is not integer.", key)
	}
	return value.Value.(int), nil
}

// GetString returns the string value from the memory cache
func (c *MemoryCache) GetString(key string) (string, error) {
	value, ok := c.get(key)
	if !ok {
		return "", new(ErrCacheMiss)
	}
	if value.Type != ValueString {
		log.Printf("[Warning] The value type of the key [%s] is not string.", key)
	}
	return value.Value.(string), nil
}

// GetFloat returns the float value from the memory cache
func (c *MemoryCache) GetFloat(key string) (float64, error) {
	value, ok := c.get(key)
	if !ok {
		return 0, new(ErrCacheMiss)
	}
	if value.Type != ValueFloat {
		log.Printf("[Warning] The value type of the key [%s] is not float.", key)
	}
	return value.Value.(float64), nil
}

// GetBool returns the boolean value from the memory cache
func (c *MemoryCache) GetBool(key string) (bool, error) {
	value, ok := c.get(key)
	if !ok {
		return false, new(ErrCacheMiss)
	}
	if value.Type != ValueBool {
		log.Printf("[Warning] The value type of the key [%s] is not boolean.", key)
	}
	return value.Value.(bool), nil
}

// GetArray returns the array value from the memory cache
func (c *MemoryCache) GetArray(key string) ([]any, error) {
	value, ok := c.get(key)
	if !ok {
		return nil, new(ErrCacheMiss)
	}
	if value.Type != ValueArray {
		log.Printf("[Warning] The value type of the key [%s] is not array.", key)
	}
	return value.Value.([]any), nil
}

// GetMap returns the map value from the memory cache
func (c *MemoryCache) GetMap(key string) (map[string]any, error) {
	value, ok := c.get(key)
	if !ok {
		return nil, new(ErrCacheMiss)
	}
	if value.Type != ValueMap {
		log.Printf("[Warning] The value type of the key [%s] is not map.", key)
	}
	return value.Value.(map[string]any), nil
}

// Set sets the value to the memory cache
func (c *MemoryCache) Set(key string, value any) error {
	return c.SetEx(key, value, c.DefaultTTL)
}

// SetEx sets the value to the memory cache with a specific TTL
func (c *MemoryCache) SetEx(key string, value any, ttl int) error {
	valueType := getValueType(value, 0)
	if valueType == ValueUnknown {
		return new(ErrInvalidValueType)
	}

	c.MemoryCacheRWLock.Lock()
	defer c.MemoryCacheRWLock.Unlock()

	// If the maximum number of entries is reached, remove the oldest entry
	if c.MaxEntries > 0 && len(c.MemoryCacheMap) >= c.MaxEntries {
		// Remove the oldest entry
		var oldestKey string
		var oldestAccessTime time.Time
		for k, v := range c.MemoryCacheMap {
			if oldestAccessTime.IsZero() || v.LastAccess.Before(oldestAccessTime) {
				oldestKey = k
				oldestAccessTime = v.LastAccess
			}
		}
		delete(c.MemoryCacheMap, oldestKey)
	}

	c.MemoryCacheMap[key] = valueMemory{
		Type:       valueType,
		Value:      value,
		Expiration: ttl,
		LastAccess: time.Now(),
	}
	return nil
}

// Del deletes the value from the memory cache
func (c *MemoryCache) Del(key string) error {
	c.MemoryCacheRWLock.Lock()
	defer c.MemoryCacheRWLock.Unlock()

	delete(c.MemoryCacheMap, key)
	return nil
}

// Exists checks whether the key exists in the memory cache
func (c *MemoryCache) Exists(key string) (bool, error) {
	c.MemoryCacheRWLock.RLock()
	defer c.MemoryCacheRWLock.RUnlock()

	_, ok := c.MemoryCacheMap[key]
	return ok, nil
}

// Clear clears the memory cache
func (c *MemoryCache) Clear() error {
	c.MemoryCacheRWLock.Lock()
	defer c.MemoryCacheRWLock.Unlock()

	c.MemoryCacheMap = make(map[string]valueMemory)
	return nil
}

// RefreshTTL refreshes the time to live of the key
func (c *MemoryCache) RefreshTTL(key string, ttl int) error {
	c.MemoryCacheRWLock.Lock()
	defer c.MemoryCacheRWLock.Unlock()

	// If the TTL le 0, delete the key
	if ttl <= 0 {
		delete(c.MemoryCacheMap, key)
		return nil
	}

	value, ok := c.MemoryCacheMap[key]
	if ok {
		value.Expiration = ttl
		c.MemoryCacheMap[key] = value
	}
	return nil
}

// garbageCollection performs the garbage collection, checking the expiration of the values
func (c *MemoryCache) garbageCollection() {
	c.MemoryCacheRWLock.Lock()
	defer c.MemoryCacheRWLock.Unlock()

	for key, value := range c.MemoryCacheMap {
		if value.Expiration > 0 && time.Now().After(value.LastAccess.Add(time.Duration(value.Expiration)*time.Second)) {
			delete(c.MemoryCacheMap, key)
		}
	}
}
