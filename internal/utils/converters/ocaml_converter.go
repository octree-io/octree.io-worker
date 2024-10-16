package converters

import (
	"fmt"
	"strconv"
)

func JsonToOCaml(value interface{}) (string, error) {
	if value == nil {
		return "None", nil
	}

	switch v := value.(type) {
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64), nil
	case bool:
		if v {
			return "true", nil
		}
		return "false", nil
	case string:
		return fmt.Sprintf("\"%s\"", v), nil
	case []interface{}:
		var result string
		for i, elem := range v {
			ocamlVal, err := JsonToOCaml(elem)
			if err != nil {
				return "", err
			}
			if i > 0 {
				result += "; "
			}
			result += ocamlVal
		}
		return "[" + result + "]", nil
	case map[string]interface{}:
		var result string
		for k, v := range v {
			ocamlVal, err := JsonToOCaml(v)
			if err != nil {
				return "", err
			}
			result += fmt.Sprintf("let %s = %s in\n    ", k, ocamlVal)
		}
		return result, nil
	default:
		return "", fmt.Errorf("unsupported type: %T", value)
	}
}
