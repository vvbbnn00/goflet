// Package model provides the model for the cache
package model

// ValueType is the type of the value
type ValueType int

const (
	// ValueInt is the integer value type
	ValueInt ValueType = iota
	// ValueString is the string value type
	ValueString
	// ValueFloat is the float value type
	ValueFloat
	// ValueBool is the boolean value type
	ValueBool
	// ValueArray is the array value type
	ValueArray
	// ValueMap is the map value type
	ValueMap
	// ValueUnknown is the unknown value type
	ValueUnknown = -1
)

// Cache is the interface for the cache
type Cache interface {
	// GetInt returns the integer value from the cache
	GetInt(key string) (int, error)
	// GetString returns the string value from the cache
	GetString(key string) (string, error)
	// GetFloat returns the float value from the cache
	GetFloat(key string) (float64, error)
	// GetBool returns the boolean value from the cache
	GetBool(key string) (bool, error)
	// GetArray returns the array value from the cache
	GetArray(key string) ([]any, error)
	// GetMap returns the map value from the cache
	GetMap(key string) (map[string]any, error)

	// Set sets the value to the cache
	Set(key string, value any) error
	SetEx(key string, value any, ttl int) error

	// Del deletes the value from the cache
	Del(key string) error

	// Exists checks whether the key exists in the cache
	Exists(key string) (bool, error)

	// Clear clears the cache
	Clear() error

	// RefreshTTL refreshes the time to live of the key
	RefreshTTL(key string, ttl int) error
}

// ErrCacheMiss is the error for cache miss
type ErrCacheMiss struct{}

// Error returns the error message
func (m *ErrCacheMiss) Error() string {
	return "cache: key not found"
}

// ErrInvalidValueType is the error for invalid value type
type ErrInvalidValueType struct{}

// Error returns the error message
func (m *ErrInvalidValueType) Error() string {
	return "cache: invalid value type"
}

// ErrTypeMismatch is the error for type mismatch
type ErrTypeMismatch struct{}

// Error returns the error message
func (m *ErrTypeMismatch) Error() string {
	return "cache: type mismatch"
}
