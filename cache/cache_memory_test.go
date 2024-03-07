package cache

import (
	"goflet/config"
	"log"
	"testing"
	"time"
)

// init initializes the configuration
func init() {
	log.Printf("Initializing the configuration...")
	config.GofletCfg.CacheConfig.CacheType = config.CacheTypeMemory
	config.GofletCfg.CacheConfig.MemoryCache.MaxEntries = 10
	config.GofletCfg.CacheConfig.MemoryCache.DefaultTTL = 60
	log.Printf("Configuration initialized.")
}

func TestGetCache(t *testing.T) {
	c := GetCache()
	if c == nil {
		t.Errorf("The cache instance should not be nil.")
	}
}

func TestSet(t *testing.T) {
	c := GetCache()
	// Set integer
	err := c.Set("integer", 1)
	if err != nil {
		t.Errorf("The error should be nil, but got %v.", err)
	}
	// Set float
	err = c.Set("float", 1.1)
	if err != nil {
		t.Errorf("The error should be nil, but got %v.", err)
	}
	// Set boolean
	err = c.Set("boolean", true)
	if err != nil {
		t.Errorf("The error should be nil, but got %v.", err)
	}
	// Set string
	err = c.Set("string", "value")
	if err != nil {
		t.Errorf("The error should be nil, but got %v.", err)
	}
}

func TestGetInt(t *testing.T) {
	c := GetCache()
	actual, err := c.GetInt("integer")
	if err != nil {
		t.Errorf("The error should be nil, but got %v.", err)
	}
	if actual != 1 {
		t.Errorf("The actual value should be [1], but got [%d].", actual)
	}
}

func TestGetFloat(t *testing.T) {
	c := GetCache()
	actual, err := c.GetFloat("float")
	if err != nil {
		t.Errorf("The error should be nil, but got %v.", err)
	}
	if actual != 1.1 {
		t.Errorf("The actual value should be [1.1], but got [%f].", actual)
	}
}

func TestGetBool(t *testing.T) {
	c := GetCache()
	actual, err := c.GetBool("boolean")
	if err != nil {
		t.Errorf("The error should be nil, but got %v.", err)
	}
	if actual != true {
		t.Errorf("The actual value should be [true], but got [%t].", actual)
	}
}

func TestGetString(t *testing.T) {
	c := GetCache()
	actual, err := c.GetString("string")
	if err != nil {
		t.Errorf("The error should be nil, but got %v.", err)
	}
	if actual != "value" {
		t.Errorf("The actual value should be [value], but got [%s].", actual)
	}
}

func TestExists(t *testing.T) {
	c := GetCache()
	exists, err := c.Exists("integer")
	if err != nil {
		t.Errorf("The error should be nil, but got %v.", err)
	}
	if !exists {
		t.Errorf("The key [integer] should exist.")
	}

	_, err = c.Exists("nonexistent")
	if err != nil {
		t.Errorf("The error should be nil, but got %v.", err)
	}
}

func TestDel(t *testing.T) {
	c := GetCache()
	err := c.Del("integer")
	if err != nil {
		t.Errorf("The error should be nil, but got %v.", err)
	}
	_, err = c.GetInt("integer")
	if err == nil {
		t.Errorf("The error should not be nil.")
	}
}

func TestClear(t *testing.T) {
	c := GetCache()
	err := c.Clear()
	if err != nil {
		t.Errorf("The error should be nil, but got %v.", err)
	}
	_, err = c.GetInt("float")
	if err == nil {
		t.Errorf("The error should not be nil.")
	}
}

func TestRefreshTTL(t *testing.T) {
	c := GetCache()
	err := c.SetEx("integer", 1, 1)
	if err != nil {
		t.Errorf("The error should be nil, but got %v.", err)
	}
	time.Sleep(1100 * time.Millisecond) // Sleep for 1.1 seconds
	err = c.RefreshTTL("integer", 2)
	if err != nil {
		t.Errorf("The error should be nil, but got %v.", err)
	}
	_, err = c.GetInt("integer")
	if err != nil {
		t.Errorf("The error should be nil, but got %v.", err)
	}
}

func TestSetEx(t *testing.T) {
	c := GetCache()
	err := c.SetEx("integer", 1, 1)
	if err != nil {
		t.Errorf("The error should be nil, but got %v.", err)
	}
	time.Sleep(2 * time.Second)
	_, err = c.GetInt("integer")
	if err == nil {
		t.Errorf("The error should not be nil.")
	}
}

func TestGarbageCollection(t *testing.T) {
	c := GetCache()
	err := c.SetEx("integer", 1, 1)
	if err != nil {
		t.Errorf("The error should be nil, but got %v.", err)
	}
	time.Sleep(10 * time.Second)
	_, err = c.GetInt("integer")
	if err == nil {
		t.Errorf("The error should not be nil.")
	}
}

func TestManyEntries(t *testing.T) {
	c := GetCache()
	for i := 0; i < 20; i++ {
		err := c.SetEx(string(rune(i)), i, 0)
		if err != nil {
			t.Errorf("The error should be nil, but got %v.", err)
		}
	}
	_, err := c.GetInt("0")
	if err == nil {
		t.Errorf("The error should not be nil.")
	}
}
