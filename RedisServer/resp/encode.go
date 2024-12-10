package resp

import "fmt"

func ToRESP(obj interface{}) string {
	resp := ""

	switch v := obj.(type) {
		case []interface{}:
			resp += fmt.Sprintf("*%d\r\n", len(v))
			for _, item := range v {
				resp += ToRESP(item)
			}

		case string: 
			resp += fmt.Sprintf("$%d\r\n%s\r\n", len(v), v)

		case int, int64, float64:
			resp += fmt.Sprintf(":%v\r\n", v)

		default:
			resp += ""
	}

	return resp
}

func HandleErrors() []string{
	return ToSimpleError("Invalid Command/Syntax")
}

func ToSimpleString(str string) []string{
	return []string{fmt.Sprintf("+%s\r\n",str)}
}

func ToSimpleError(str string) []string{
	return []string{fmt.Sprintf("-%s\r\n",str)}
}

func ToNullBulkString() []string{
	return ToSimpleError("1")
}