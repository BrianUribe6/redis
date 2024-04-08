package store

var db map[string]string = make(map[string]string)

func Set(key string, value string) {
	db[key] = value
}

func Get(key string) (string, bool) {
	value, exist := db[key]
	return value, exist
}
