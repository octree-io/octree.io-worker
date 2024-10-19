package utils

import (
	"encoding/json"
	"fmt"
	"strings"
)

func ConvertToString(value interface{}) string {
	return fmt.Sprintf("%v", value)
}

func ConvertToJSONString(value interface{}) (string, error) {
	if str, ok := value.(string); ok {
		return str, nil
	}

	if b, ok := value.(bool); ok {
		return fmt.Sprintf("%v", b), nil
	}

	if i, ok := value.(int); ok {
		return fmt.Sprintf("%d", i), nil
	}

	jsonBytes, err := json.Marshal(value)
	if err != nil {
		return "", fmt.Errorf("failed to convert to JSON: %v", err)
	}
	return string(jsonBytes), nil
}

func SplitStringIntoParts(input string, delimiter string) []string {
	parts := strings.Split(input, delimiter)

	var result []string
	for _, part := range parts {
		if strings.TrimSpace(part) != "" {
			result = append(result, part)
		}
	}

	return result
}
