package store

import "time"

var db map[string]*Item = make(map[string]*Item)

type Item struct {
	value     string
	expiresAt *time.Time
}

func Set(key string, value string, expiry int64) {
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
	item, exist := db[key]
	if exist && item.expiresAt != nil && item.expiresAt.Compare(time.Now()) < 0 {
		delete(db, key)
		return "", false

	}
	return item.value, exist
}
