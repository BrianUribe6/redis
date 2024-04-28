package store

import (
	"sync"
	"time"
)

var db = make(map[string]*Item)

var mut sync.Mutex = sync.Mutex{}

type Item struct {
	value     string
	expiresAt *time.Time
}

func Set(key string, value string, expiry int64) {
	mut.Lock()
	defer mut.Unlock()
	item := new(Item)
	item.value = value
	if expiry > 0 {
		duration := time.Duration(expiry) * time.Millisecond
		expiresAt := time.Now().Add(duration)
		item.expiresAt = &expiresAt
	}
	db[key] = item
}

func Get(key string) (string, bool) {
	mut.Lock()
	defer mut.Unlock()
	item, exist := db[key]
	if exist {
		if item.expiresAt != nil && item.expiresAt.Before(time.Now()) {
			delete(db, key)
			return "", false
		}
		return item.value, exist
	}
	return "", false
}
