package resp

import (
	"Redis/myConfig"
	"fmt"
	"strings"
	"strconv"
)

func ParseMessage(data []interface{}) myConfig.Content{
	parsedData := ParseRESP(data)
	if get, chk := parsedData[0].([]interface{}); chk {
		parsedData = get
	} else{
		parsedData = parsedData[:len(parsedData)-1]
	}

	var result myConfig.Content
	result.Cmd,_ = parsedData[0].(string) 
	if strings.ToUpper(result.Cmd) == "CONFIG" {
		part2,_ := parsedData[1].(string) 
		result.Cmd = fmt.Sprintf("%v %v",strings.ToUpper(result.Cmd), strings.ToUpper(part2))
		result.Args = parsedData[2:]
	}else{
		result.Args = parsedData[1:]
	}

	return result
}

func ParseRESP(data []interface{}) []interface{}{
	for len(data) > 0{
		element,_ := data[0].(string)
		data = data[1:]
		switch element[0] {
			case '+':
				var arr []interface{}
				arr = append(arr, element[1:])
				arr = append(arr, data)
				return arr
			
			case '*':
				arrlen,_ := strconv.Atoi(element[1:])
				var arr []interface{}
				var values []interface{}
				for j := 0; j < arrlen; j++ {
					parsedContent := ParseRESP(data)
					values = append(values, parsedContent[0])
					data,_ = parsedContent[1].([]interface{})
				}
				arr = append(arr, values)
				arr = append(arr, data)
				return arr;

			case '$':
				var arr []interface{}
				str,_ := data[0].(string)
				data = data[1:]
				arr = append(arr, str)
				arr = append(arr, data)
				return arr
			
			case ':':
				var arr []interface{}
				n,_ := strconv.Atoi(element[1:])
				arr = append(arr, n)
				arr = append(arr, data)
				return arr
		}
	}
	return make([]interface{}, 0)
}
