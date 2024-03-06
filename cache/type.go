package cache

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
	ValueUnknown = -1
)

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

type ErrCacheMiss struct{}

func (m *ErrCacheMiss) Error() string {
	return "cache: key not found"
}

type ErrInvalidValueType struct{}

func (m *ErrInvalidValueType) Error() string {
	return "cache: invalid value type"
}

type ErrTypeMismatch struct{}

func (m *ErrTypeMismatch) Error() string {
	return "cache: type mismatch"
}
