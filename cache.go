package static

import (
	"sync"
	"time"
)

var hashCacheData = &hashCache{store: make(map[string]*cacheFileInfo)}

type hashCache struct {
	sync.RWMutex
	store map[string]*cacheFileInfo
}

type cacheFileInfo struct {
	Hash    string
	ModTime time.Time
}

func (c *hashCache) Get(path string) *cacheFileInfo {
	c.RLock()
	defer c.RUnlock()
	return c.store[path]
}

func (c *hashCache) Set(path, hash string, mod time.Time) {
	c.Lock()
	defer c.Unlock()
	c.store[path] = &cacheFileInfo{hash, mod}
}
