package checker

import "time"

// retrieveResultFromCache returns the latest check result from cache, if any.
// If the result is expired or none is available, this function will return nil.
func (c *Checker) retrieveResultFromCache() *Result {
	if c.cachedResult != nil && c.cachedResult.expiration.After(time.Now()) {
		return c.cachedResult.result
	}

	return nil
}

// cacheResult sets a check result to the cache and expires it after the
// Checker.cacheTTL is exceeded.
func (c *Checker) cacheResult(result *Result) {
	c.cachedResult = &CachedResult{
		result:     result,
		expiration: time.Now().Add(c.cacheTTL),
	}
}
