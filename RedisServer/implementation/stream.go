package implementation

import (
	"Redis/resp"
	"Redis/store"
	"fmt"
)

func HandleType(args []interface{}) []string {
	key := args[0]
	var exist bool
	if _, exist = store.StoredKeys[key]; exist {
		return resp.ToSimpleString("string")
	} else {
		return resp.ToSimpleString("none")
	}
}

func HandleXadd(args []interface{}) []string {
	val, ok := args[0].(string)
	if !ok {
		return resp.HandleErrors()
	}
	streamName := val
	val, ok = args[1].(string)
	if !ok {
		return resp.HandleErrors()
	}
	timeline := val

	if _, ok := store.Streams[streamName]; !ok {
		store.Streams[streamName] = make(map[string]map[interface{}]interface{})
	}

	if _, ok := store.Streams[streamName][timeline]; !ok {
		store.Streams[streamName][timeline] = make(map[interface{}]interface{})
	}
	for i := 2; i < len(args); i += 2 {
		key := args[i]
		value := args[i+1]
		store.Streams[streamName][timeline][key] = value
	}
	return resp.ToSimpleString(timeline)
}

func HandleXread(args []interface{}) []string {
	val, ok := args[0].(string)
	if !ok {
		return resp.HandleErrors()
	}
	streamName := val
	val, ok = args[1].(string)
	if !ok {
		return resp.HandleErrors()
	}
	timeline := val

	if _, ok := store.Streams[streamName]; !ok {
		return resp.ToSimpleError("Stream doesn't exist.")
	}

	entry, ok := store.Streams[streamName][timeline]
	if !ok {
		return resp.ToSimpleError("Stream's timeline doesn't exist.")
	}

	var result = fmt.Sprintf("%v\n\t%v\n", streamName, timeline)
	for key, value := range entry {
		result += fmt.Sprintf("\t\t%v\t%v\n", key, value)
	}
	return resp.ToSimpleString(result)
}
