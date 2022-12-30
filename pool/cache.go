package pool

import (
	"strings"
	"sync/atomic"

	"github.com/TrueCloudLab/frostfs-sdk-go/session"
	lru "github.com/hashicorp/golang-lru/v2"
)

type sessionCache struct {
	cache        *lru.Cache[string, *cacheValue]
	currentEpoch uint64
}

type cacheValue struct {
	token session.Object
}

func newCache() (*sessionCache, error) {
	cache, err := lru.New[string, *cacheValue](100)
	if err != nil {
		return nil, err
	}

	return &sessionCache{cache: cache}, nil
}

// Get returns a copy of the session token from the cache without signature
// and context related fields. Returns nil if token is missing in the cache.
// It is safe to modify and re-sign returned session token.
func (c *sessionCache) Get(key string) (session.Object, bool) {
	value, ok := c.cache.Get(key)
	if !ok {
		return session.Object{}, false
	}

	if c.expired(value) {
		c.cache.Remove(key)
		return session.Object{}, false
	}

	return value.token, true
}

func (c *sessionCache) Put(key string, token session.Object) bool {
	return c.cache.Add(key, &cacheValue{
		token: token,
	})
}

func (c *sessionCache) DeleteByPrefix(prefix string) {
	for _, key := range c.cache.Keys() {
		if strings.HasPrefix(key, prefix) {
			c.cache.Remove(key)
		}
	}
}

func (c *sessionCache) updateEpoch(newEpoch uint64) {
	epoch := atomic.LoadUint64(&c.currentEpoch)
	if newEpoch > epoch {
		atomic.StoreUint64(&c.currentEpoch, newEpoch)
	}
}

func (c *sessionCache) expired(val *cacheValue) bool {
	epoch := atomic.LoadUint64(&c.currentEpoch)
	return val.token.ExpiredAt(epoch)
}
