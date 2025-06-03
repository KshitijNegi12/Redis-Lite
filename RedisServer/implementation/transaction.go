package implementation

import (
	"Redis/resp"
	"Redis/store"
	"net"
	"strconv"
	"time"
)

func HandleIncr(args []interface{}) []string {
	key := args[0]
	if CheckForKeys(key) {
		if val, ok := store.StoredKeys[key]; ok {
			return resp.ToSimpleString(strconv.Itoa(val.(int)))
		}
	}
	return resp.ToSimpleError("value is not an integer or out of range")
}

func CheckForKeys(key interface{}) bool {
	if val, exists := store.StoredKeys[key]; exists {
		numVal, ok := val.(int)
		if !ok {
			return false
		}

		currTime := time.Now()   
        expTime, hasExpiry := store.ExpiryKeys[key]
        if !hasExpiry || currTime.Before(expTime) {
            store.StoredKeys[key] = numVal + 1
        } else {
            delete(store.ExpiryKeys, key)
			store.StoredKeys[key] = 1
        }

		return true
	}

	store.StoredKeys[key] = 1
	return true
}

func HandleMulti(conn net.Conn) []string {
	if store.AddConnToMultiQueue(conn) {
		return resp.ToSimpleString("OK")
	}
	return resp.ToSimpleError("ERR could not start transaction")
}

func HandleDiscard(conn net.Conn) []string {
	if !store.CheckConnInQueue(conn) {
		return resp.ToSimpleError("ERR DISCARD without MULTI")
	}
	store.DiscardQueueCmds(conn)
	return resp.ToSimpleString("OK")
}
