package resp

import "fmt"

func HandleErrors() []string{
	return []string{ToSimpleError("Invalid Command/Syntax")}
}

func ToSimpleString(str string) string{
	return fmt.Sprintf("+%s\r\n",str)
}

func ToSimpleError(str string) string{
	return fmt.Sprintf("-%s\r\n",str)
}

