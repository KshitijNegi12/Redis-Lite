package implementation

import (
	"Redis/myConfig"
	"Redis/resp"
	"Redis/store"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"
)

func HandlePing() []string {
	return resp.ToSimpleString("PONG")
}

func HandleEcho(args []interface{}) []string {
	var str string
	for _, item := range args {
		str += fmt.Sprintf("%v ", item)
	}
	return resp.ToSimpleString(str)
}

func HandleSet(args []interface{}, config *myConfig.Config) []string {
	key := args[0]
	value := args[1]
	delete(store.ExpiryKeys, key)
	if len(args) > 2 {
		err := handleKeyExpiry(args, key)
		if err != nil {
			return resp.HandleErrors()
		}
	}
	store.StoredKeys[key] = value
	if config.Role == "master" {
		data := []interface{}{"SET", key, value}
		sendPropogationToReplicas(resp.ToRESP(data), config)
		return resp.ToSimpleString("OK")
	}
	return resp.ToSimpleString("")
}

func handleKeyExpiry(args []interface{}, key interface{}) error {
	typeExpiry, ok := args[2].(string)
	if !ok {
		return errors.New("invalid type")
	}
	typeExpiry = strings.ToUpper(typeExpiry)

	timeExpiry, err := parseInt(args[3])
	if err != nil {
		return err
	}

	var expiryTime time.Time
	if typeExpiry == "PX" {
		duration := time.Duration(timeExpiry) * time.Millisecond
		expiryTime = time.Now().Add(duration)
	} else if typeExpiry == "EX" {
		duration := time.Duration(timeExpiry) * time.Second
		expiryTime = time.Now().Add(duration)
	} else {
		return errors.New("invalid type")
	}

	store.ExpiryKeys[key] = expiryTime
	return nil
}

func parseInt(str interface{}) (int, error) {
	result, ok := str.(int)
	if ok {
		return result, nil
	}
	return -1, errors.New("invalid type")
}

func HandleGet(args []interface{}) []string {
	key := args[0]
	var val interface{}
	var exist bool
	if val, exist = store.StoredKeys[key]; exist {
		currTime := time.Now()
		expTime, hasExpiry := store.ExpiryKeys[key]
		value := convertToString(val)
		if !hasExpiry || currTime.Before(expTime) {
			return resp.ToSimpleString(value)
		} else {
			return resp.ToNullBulkString()
		}
	} else {
		return resp.ToNullBulkString()
	}
}

func convertToString(val interface{}) string {
	switch v := val.(type) {
	case string:
		return v
	case int:
		return fmt.Sprintf("%d", v)
	default:
		log.Println("Unexpected value type, returning empty string")
		return ""
	}
}

func HandleDel(args []interface{}) []string {
	key := args[0]
	var exist bool
	if _, exist = store.StoredKeys[key]; exist {
		delete(store.StoredKeys, key)
		if _, exist = store.ExpiryKeys[key]; exist {
			delete(store.ExpiryKeys, key)
		}
		return resp.ToSimpleString("1")
	} else {
		return resp.ToNullBulkString()
	}
}

func sendPropogationToReplicas(data string, config *myConfig.Config) {
	for conn := range config.ConnectedSlaves {
		if conn != nil {
			_, err := conn.Write([]byte(data))
			if err != nil {
				continue
			}
		}
	}
}
