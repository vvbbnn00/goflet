package memory_cache

import (
	"errors"
	"goflet/cache/model"
	"reflect"
	"sync"
	"time"
)

const garbageCollectionInterval = 10 * time.Second

type valueMemory struct {
	Type       model.ValueType
	Value      interface{}
	Expiration int       // <= 0 means no expiration
	LastAccess time.Time // The last access time, used to determine whether the value is expired
}

type MemoryCache struct {
	MaxEntries int // The maximum number of entries to be stored in the cache
	DefaultTTL int // The default time to live for the cache

	memoryCacheMap    map[string]valueMemory // The memory cache map
	memoryCacheRWLock *sync.RWMutex          // The read-write lock for the memory cache
	stopGC            chan struct{}          // Channel to stop the garbage collection goroutine
}

var errCacheMiss = errors.New("cache miss")
var errInvalidValueType = errors.New("invalid value type")
var errTypeMismatch = errors.New("type mismatch")

// getValueType returns the value type
func getValueType(value interface{}) model.ValueType {
	if value == nil {
		return model.ValueUnknown
	}

	switch reflect.TypeOf(value).Kind() {
	case reflect.Int:
		return model.ValueInt
	case reflect.String:
		return model.ValueString
	case reflect.Float64:
		return model.ValueFloat
	case reflect.Bool:
		return model.ValueBool
	// case reflect.Slice, reflect.Array:
	// 	return ValueArray
	// case reflect.Map:
	// 	return ValueMap
	default:
		return model.ValueUnknown
	}
}

// NewMemoryCache creates a new memory cache
func NewMemoryCache(maxEntries, defaultTTL int) *MemoryCache {
	instance := &MemoryCache{
		MaxEntries:        maxEntries,
		DefaultTTL:        defaultTTL,
		memoryCacheMap:    make(map[string]valueMemory),
		memoryCacheRWLock: new(sync.RWMutex),
		stopGC:            make(chan struct{}),
	}

	// Start the garbage collection goroutine
	go instance.garbageCollection()

	return instance
}

// get returns the value from the memory cache
func (c *MemoryCache) get(key string) (valueMemory, bool) {
	c.memoryCacheRWLock.RLock()
	defer c.memoryCacheRWLock.RUnlock()

	value, ok := c.memoryCacheMap[key]
	if ok && value.Expiration > 0 && time.Now().After(value.LastAccess.Add(time.Duration(value.Expiration)*time.Second)) {
		// The value is expired
		return valueMemory{}, false
	}
	return value, ok
}

// GetInt returns the integer value from the memory cache
func (c *MemoryCache) GetInt(key string) (int, error) {
	value, ok := c.get(key)
	if !ok {
		return 0, errCacheMiss
	}
	if value.Type != model.ValueInt {
		return 0, errTypeMismatch
	}
	return value.Value.(int), nil
}

// GetString returns the string value from the memory cache
func (c *MemoryCache) GetString(key string) (string, error) {
	value, ok := c.get(key)
	if !ok {
		return "", errCacheMiss
	}
	if value.Type != model.ValueString {
		return "", errTypeMismatch
	}
	return value.Value.(string), nil
}

// GetFloat returns the float value from the memory cache
func (c *MemoryCache) GetFloat(key string) (float64, error) {
	value, ok := c.get(key)
	if !ok {
		return 0, errCacheMiss
	}
	if value.Type != model.ValueFloat {
		return 0, errTypeMismatch
	}
	return value.Value.(float64), nil
}

// GetBool returns the boolean value from the memory cache
func (c *MemoryCache) GetBool(key string) (bool, error) {
	value, ok := c.get(key)
	if !ok {
		return false, errCacheMiss
	}
	if value.Type != model.ValueBool {
		return false, errTypeMismatch
	}
	return value.Value.(bool), nil
}

// GetArray returns the array value from the memory cache
func (c *MemoryCache) GetArray(key string) ([]interface{}, error) {
	return nil, errTypeMismatch // Not implemented
}

// GetMap returns the map value from the memory cache
func (c *MemoryCache) GetMap(key string) (map[string]interface{}, error) {
	return nil, errTypeMismatch // Not implemented
}

// Set sets the value to the memory cache
func (c *MemoryCache) Set(key string, value interface{}) error {
	return c.SetEx(key, value, c.DefaultTTL)
}

// SetEx sets the value to the memory cache with a specific TTL
func (c *MemoryCache) SetEx(key string, value interface{}, ttl int) error {
	valueType := getValueType(value)
	if valueType == model.ValueUnknown {
		return errInvalidValueType
	}

	c.memoryCacheRWLock.Lock()
	defer c.memoryCacheRWLock.Unlock()

	// If the maximum number of entries is reached, remove the oldest entry
	if c.MaxEntries > 0 && len(c.memoryCacheMap) >= c.MaxEntries {
		// Remove the oldest entry
		var oldestKey string
		var oldestAccessTime time.Time
		for k, v := range c.memoryCacheMap {
			if oldestAccessTime.IsZero() || v.LastAccess.Before(oldestAccessTime) {
				oldestKey = k
				oldestAccessTime = v.LastAccess
			}
		}
		delete(c.memoryCacheMap, oldestKey)
	}

	c.memoryCacheMap[key] = valueMemory{
		Type:       valueType,
		Value:      value,
		Expiration: ttl,
		LastAccess: time.Now(),
	}
	return nil
}

// Del deletes the value from the memory cache
func (c *MemoryCache) Del(key string) error {
	c.memoryCacheRWLock.Lock()
	defer c.memoryCacheRWLock.Unlock()

	delete(c.memoryCacheMap, key)
	return nil
}

// Exists checks whether the key exists in the memory cache
func (c *MemoryCache) Exists(key string) (bool, error) {
	c.memoryCacheRWLock.RLock()
	defer c.memoryCacheRWLock.RUnlock()

	_, ok := c.memoryCacheMap[key]
	return ok, nil
}

// Clear clears the memory cache
func (c *MemoryCache) Clear() error {
	c.memoryCacheRWLock.Lock()
	defer c.memoryCacheRWLock.Unlock()

	c.memoryCacheMap = make(map[string]valueMemory)
	return nil
}

// RefreshTTL refreshes the time to live of the key
func (c *MemoryCache) RefreshTTL(key string, ttl int) error {
	c.memoryCacheRWLock.Lock()
	defer c.memoryCacheRWLock.Unlock()

	// If the TTL is <= 0, delete the key
	if ttl <= 0 {
		delete(c.memoryCacheMap, key)
		return nil
	}

	value, ok := c.memoryCacheMap[key]
	if ok {
		value.Expiration = ttl
		c.memoryCacheMap[key] = value
	}
	return nil
}

// garbageCollection performs the garbage collection, checking the expiration of the values
func (c *MemoryCache) garbageCollection() {
	ticker := time.NewTicker(garbageCollectionInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.memoryCacheRWLock.Lock()
			for key, value := range c.memoryCacheMap {
				if value.Expiration > 0 && time.Now().After(value.LastAccess.Add(time.Duration(value.Expiration)*time.Second)) {
					delete(c.memoryCacheMap, key)
				}
			}
			c.memoryCacheRWLock.Unlock()
		case <-c.stopGC:
			return
		}
	}
}

// Close stops the garbage collection goroutine and clears the memory cache
func (c *MemoryCache) Close() {
	close(c.stopGC)
	_ = c.Clear()
}
