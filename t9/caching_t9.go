package t9

import (
	"fmt"

	"github.com/hashicorp/golang-lru"
	"github.com/pkg/errors"
)

// NewCachingT9 instantiates a new T9 structure that caches lookup results using an LRU cache
// of the designated maximum size.
func NewCachingT9(cacheSize int) (T9, error) {

	cache, err := lru.New(cacheSize)
	if err != nil {
		return nil, errors.Wrap(err, `could not instantiate LRU cache`)
	}

	return &cachingT9{
		T9:    New(),
		cache: cache,
	}, nil
}

type cachingT9 struct {
	T9
	cache *lru.Cache
}

func (t9 *cachingT9) GetWords(digits string, exact bool) ([]string, error) {

	cacheKey := fmt.Sprintf(`digits=%s | exact=%t`, digits, exact)

	if wordsInterface, ok := t9.cache.Get(cacheKey); ok {
		return wordsInterface.([]string), nil
	}

	words, err := t9.T9.GetWords(digits, exact)
	if err != nil {
		return nil, err
	}

	t9.cache.Add(cacheKey, words)
	return words, nil
}
