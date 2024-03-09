package cache

import (
	"goflet/cache/memory_cache"
	"goflet/cache/model"
	"goflet/config"
)

var gCache model.Cache

// initConfig initializes the cache configuration
func initConfig() {
	cacheType := config.GofletCfg.CacheConfig.CacheType
	switch cacheType {
	case config.CacheTypeMemory:
		gCache = memory_cache.NewMemoryCache(config.GofletCfg.CacheConfig.MemoryCache.MaxEntries, config.GofletCfg.CacheConfig.MemoryCache.DefaultTTL)
	default:
		panic("Unknown cache type, please check the configuration")
	}
}

// GetCache returns the cache instance (which should be a singleton)
func GetCache() model.Cache {
	if gCache == nil {
		initConfig()
	}
	return gCache
}
