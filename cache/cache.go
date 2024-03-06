package cache

import "goflet/config"

var gCache Cache

// initConfig initializes the cache configuration
func initConfig() {
	cacheType := config.GofletCfg.CacheConfig.CacheType
	switch cacheType {
	case config.CacheTypeMemory:
		gCache = NewMemoryCache(config.GofletCfg.CacheConfig.MemoryCache.MaxEntries, config.GofletCfg.CacheConfig.MemoryCache.DefaultTTL)
	default:
		panic("Unknown cache type, please check the configuration")
	}
}

// GetCache returns the cache instance (which should be a singleton)
func GetCache() Cache {
	if gCache == nil {
		initConfig()
	}
	return gCache
}
