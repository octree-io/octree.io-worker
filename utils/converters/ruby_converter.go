package converters

import (
	"fmt"
)

func JsonToRuby(data interface{}) (string, error) {
	switch val := data.(type) {
	case map[string]string:
		result := "{"
		for k, v := range val {
			result += fmt.Sprintf("\"%s\" => \"%s\", ", k, v)
		}
		if len(val) > 0 {
			result = result[:len(result)-2]
		}
		result += "}"
		return result, nil

	case []map[string]interface{}:
		result := "["
		for _, testCase := range val {
			caseStr := "{"
			for k, v := range testCase {
				switch v := v.(type) {
				case string:
					caseStr += fmt.Sprintf("\"%s\" => \"%s\", ", k, v)
				case int:
					caseStr += fmt.Sprintf("\"%s\" => %d, ", k, v)
				case []interface{}:
					listStr := ""
					for _, item := range v {
						switch item := item.(type) {
						case int:
							listStr += fmt.Sprintf("%d, ", item)
						case string:
							listStr += fmt.Sprintf("\"%s\", ", item)
						}
					}
					if len(v) > 0 && len(listStr) > 2 {
						listStr = listStr[:len(listStr)-2]
					}
					caseStr += fmt.Sprintf("\"%s\" => [%s], ", k, listStr)
				default:
					caseStr += fmt.Sprintf("\"%s\" => %v, ", k, v)
				}
			}
			if len(testCase) > 0 {
				caseStr = caseStr[:len(caseStr)-2]
			}
			caseStr += "}, "
			result += caseStr
		}
		if len(val) > 0 {
			result = result[:len(result)-2]
		}
		result += "]"
		return result, nil

	default:
		return "", fmt.Errorf("unsupported input type")
	}
}
