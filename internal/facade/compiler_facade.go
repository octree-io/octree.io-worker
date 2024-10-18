package facade

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os/exec"

	"octree.io-worker/internal/helpers"
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

func CompilerExplorer(language string, code string) (string, error) {
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

func ExecuteJavaScript(language string, code string) (string, string, error) {
	tmpFolderDir, err := helpers.CreateTempNpmPackage(language)
	if err != nil {
		return "", "", fmt.Errorf("failed to create temp npm package: %w", err)
	}

	err = helpers.WriteIndexFile(language, tmpFolderDir, code)
	if err != nil {
		helpers.CleanupTempNpmPackage(tmpFolderDir)
		return "", "", fmt.Errorf("failed to write index file: %w", err)
	}

	err = helpers.RunNpmInstall(tmpFolderDir)
	if err != nil {
		helpers.CleanupTempNpmPackage(tmpFolderDir)
		return "", "", fmt.Errorf("failed to run npm install: %w", err)
	}

	stdout, stderr, err := helpers.BundleNpmPackage(tmpFolderDir)
	if err != nil {
		helpers.CleanupTempNpmPackage(tmpFolderDir)
		return stdout, stderr, fmt.Errorf("failed to bundle npm package: %w", err)
	}

	stdout, stderr, err = helpers.ExecuteWasmtime(tmpFolderDir)
	if err != nil {
		helpers.CleanupTempNpmPackage(tmpFolderDir)
		return stdout, stderr, fmt.Errorf("failed to execute Wasmtime: %w", err)
	}

	err = helpers.CleanupTempNpmPackage(tmpFolderDir)
	if err != nil {
		return "", "", fmt.Errorf("failed to clean up temp package: %w", err)
	}

	return string(stdout), string(stderr), nil
}

func ExecuteTypeScript(language string, code string) (string, string, error) {
	tmpFolderDir, err := helpers.CreateTempNpmPackage(language)
	if err != nil {
		return "", "", fmt.Errorf("failed to create temp npm package: %w", err)
	}

	err = helpers.WriteIndexFile(language, tmpFolderDir, code)
	if err != nil {
		helpers.CleanupTempNpmPackage(tmpFolderDir)
		return "", "", fmt.Errorf("failed to write index file: %w", err)
	}

	err = helpers.RunNpmInstall(tmpFolderDir)
	if err != nil {
		helpers.CleanupTempNpmPackage(tmpFolderDir)
		return "", "", fmt.Errorf("failed to run npm install: %w", err)
	}

	stdout, stderr, err := compileTypeScript(tmpFolderDir)
	if err != nil {
		helpers.CleanupTempNpmPackage(tmpFolderDir)
		return stdout, stderr, fmt.Errorf("failed to compile TypeScript: %w", err)
	}

	stdout, stderr, err = helpers.BundleNpmPackage(tmpFolderDir)
	if err != nil {
		helpers.CleanupTempNpmPackage(tmpFolderDir)
		return stdout, stderr, fmt.Errorf("failed to bundle npm package: %w", err)
	}

	stdout, stderr, err = helpers.ExecuteWasmtime(tmpFolderDir)
	if err != nil {
		helpers.CleanupTempNpmPackage(tmpFolderDir)
		return stdout, stderr, fmt.Errorf("failed to execute Wasmtime: %w", err)
	}

	err = helpers.CleanupTempNpmPackage(tmpFolderDir)
	if err != nil {
		return "", "", fmt.Errorf("failed to clean up temp package: %w", err)
	}

	return string(stdout), string(stderr), nil
}

func compileTypeScript(tmpFolderDir string) (string, string, error) {
	cmd := exec.Command("npx", "tsc", "index.ts")
	cmd.Dir = tmpFolderDir

	stdoutPipe, _ := cmd.StdoutPipe()
	stderrPipe, _ := cmd.StderrPipe()

	if err := cmd.Start(); err != nil {
		log.Fatalf("Failed to start command: %v", err)
		return "", "", err
	}

	stdout, _ := io.ReadAll(stdoutPipe)
	stderr, _ := io.ReadAll(stderrPipe)

	if err := cmd.Wait(); err != nil {
		log.Printf("tsc failed. Error: %v\nstdout: %s\nstderr: %s\n", err, stdout, stderr)
		return string(stdout), string(stderr), err
	}

	return string(stdout), string(stderr), nil
}
