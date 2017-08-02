package static

import (
	"sync"
	"time"
)

var hashCacheData = &hashCache{files: make(map[string]*hashCacheFile)}

type hashCache struct {
	sync.RWMutex
	files map[string]*hashCacheFile
}

type hashCacheFile struct {
	Hash    string
	ModTime time.Time
}

func (c *hashCache) Del(path string) {
	c.Lock()
	defer c.Unlock()
	delete(c.files, path)
}

func (c *hashCache) Get(path string) *hashCacheFile {
	c.RLock()
	defer c.RUnlock()
	return c.files[path]
}

func (c *hashCache) Set(path, hash string, mod time.Time) {
	c.Lock()
	defer c.Unlock()
	c.files[path] = &hashCacheFile{hash, mod}
}
