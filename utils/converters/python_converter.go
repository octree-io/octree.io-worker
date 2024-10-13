package converters

import (
	"encoding/json"
	"fmt"
	"strings"
)

func JsonToPython(data interface{}) (string, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("error converting object to JSON: %w", err)
	}

	jsonStr := string(jsonData)
	pythonStr := strings.ReplaceAll(jsonStr, "null", "None")
	pythonStr = strings.ReplaceAll(pythonStr, "true", "True")
	pythonStr = strings.ReplaceAll(pythonStr, "false", "False")

	return pythonStr, nil
}
