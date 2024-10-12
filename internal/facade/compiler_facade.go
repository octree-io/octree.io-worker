package facade

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const API_URL = "https://godbolt.org/api"

var COMPILERS = map[string]string{
	"python":     "python312",
	"java":       "java2102",
	"cpp":        "g142",
	"csharp":     "dotnet80csharpcoreclr",
	"typescript": "tsc_0_0_35_gc",
	"ruby":       "ruby334",
	"go":         "gl1221",
	"rust":       "r1810",
	"ocaml":      "ocaml5200",
}

func Compile(language string, code string) (string, error) {
	compiler, exists := COMPILERS[language]

	if !exists {
		return "", fmt.Errorf("unsupported language: %s", language)
	}

	lang := language
	if language == "cpp" {
		lang = "c++"
	}

	payload := map[string]interface{}{
		"source":   code,
		"compiler": compiler,
		"options": map[string]interface{}{
			"userArguments":     "",
			"executeParameters": map[string]interface{}{"args": "", "stdin": "", "runtimeTools": []interface{}{}},
			"compilerOptions":   map[string]interface{}{"executorRequest": true, "skipAsm": true, "overrides": []interface{}{}},
			"filters":           map[string]interface{}{"execute": true},
			"tools":             []interface{}{},
			"libraries":         []interface{}{},
		},
		"lang":                lang,
		"files":               []interface{}{},
		"allowStoreCodeDebug": true,
	}

	url := fmt.Sprintf("%s/compiler/%s/compile", API_URL, compiler)
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return "", fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set("User-Agent", "octree.io")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("error response from API: %s", string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response: %w", err)
	}

	return string(body), nil
}
