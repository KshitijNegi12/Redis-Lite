package store

import "time"

var StoredKeys = make(map[interface{}]interface{})
var ExpiryKeys = make(map[interface{}]time.Time)