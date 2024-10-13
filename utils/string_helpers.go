package utils

import "fmt"

func ConvertToString(value interface{}) string {
	return fmt.Sprintf("%v", value)
}
