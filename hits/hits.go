package hits

import (
	"sync"

	"github.com/go-redis/redis"
)

// Hits store information about API hits
type Hits struct {
	sync.RWMutex
	store map[string]int
	redis *redis.Client
}

// NewHits return new hits
func NewHits() *Hits {
	return &Hits{
		store: make(map[string]int),
	}
}

// NewHitsRedis return Hits that
// uses Redis as storage backend
func NewHitsRedis(url string) *Hits {
	client := redis.NewClient(&redis.Options{
		Addr:     url,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	return &Hits{
		redis: client,
	}
}

// Incr increments by key
func (h *Hits) Incr(key string) {
	if h.redis != nil {
		h.redis.Incr(key)
		return
	}

	h.Lock()
	defer h.Unlock()
	h.store[key]++
}

// Stats show stats
func (h *Hits) Stats() (map[string]int, error) {
	if h.redis != nil {
		return h.statsRedis()
	}

	h.RLock()
	defer h.RUnlock()
	s := make(map[string]int)

	for k, v := range h.store {
		s[k] = v
	}

	return h.store, nil
}

func (h *Hits) statsRedis() (map[string]int, error) {
	keys, err := h.redis.Keys("*").Result()
	if err != nil {
		return nil, err
	}

	s := make(map[string]int)
	for _, v := range keys {
		res, err := h.redis.Get(v).Int()
		if err != nil {
			return nil, err
		}

		s[v] = res
	}

	return s, nil
}

// CloseRedis closes connection to redis
func (h *Hits) CloseRedis() error {
	return h.redis.Close()
}
