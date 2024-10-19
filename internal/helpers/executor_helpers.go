package helpers

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

func CreateTempNpmPackage(language string) (string, error) {
	uuidFolder := uuid.New().String()
	tmpFolderDir := fmt.Sprintf("/tmp/%s", uuidFolder)

	// Create the temp folder
	err := os.Mkdir(tmpFolderDir, os.ModePerm)
	if err != nil {
		log.Printf("Failed to create folder %s: %s", uuidFolder, err)
	}

	// Copy over prebuilt NPM package to temp folder
	// Make sure to have these packages in /root/untrusted-code-exec
	dummyPkgPath := "/root/untrusted-code-exec/dummy-pkg"
	if language == "typescript" {
		dummyPkgPath = "/root/untrusted-code-exec/dummy-ts-pkg"
	}
	err = copyDirectory(dummyPkgPath, tmpFolderDir)
	if err != nil {
		log.Printf("Failed to copy files from NPM package to %s: %s", tmpFolderDir, err)
		return "", fmt.Errorf("failed to copy NPM package")
	}

	return tmpFolderDir, nil
}

func WriteIndexFile(language string, tmpFolderDir string, code string) error {
	extension := ".js"
	if language == "typescript" {
		extension = ".ts"
	}

	filePath := fmt.Sprintf("%s/index%s", tmpFolderDir, extension)

	err := os.WriteFile(filePath, []byte(code), 0644)
	if err != nil {
		log.Printf("Failed to write to index file: %s", err)
		return fmt.Errorf("failed to write index file")
	}

	return nil
}

func RunNpmInstall(tmpFolderDir string) error {
	npmInstallCmd := exec.Command("npm", "install")
	npmInstallCmd.Dir = tmpFolderDir

	stdout, stderr, err := RunCommandWithOutput(npmInstallCmd)
	if err != nil {
		fmt.Printf("npm install failed with error: %s\n", err)
		fmt.Printf("npm install stdout: %s\n", stdout)
		fmt.Printf("npm install stderr: %s\n", stderr)
		log.Printf("npm install failed: %s", err)
		return fmt.Errorf("error while bundling with npm install")
	}

	return nil
}

func BundleNpmPackage(tmpFolderDir string) (string, string, error) {
	esbuildCmd := exec.Command("esbuild", "index.js", "--bundle", "--outfile=dist/bundle.js")
	esbuildCmd.Dir = tmpFolderDir

	stdout, stderr, err := RunCommandWithOutput(esbuildCmd)
	if err != nil {
		fmt.Printf("esbuild failed with error: %s\n", err)
		fmt.Printf("esbuild stdout: %s\n", stdout)
		fmt.Printf("esbuild stderr: %s\n", stderr)
		log.Printf("esbuild failed: %s", err)
		return stdout, stderr, fmt.Errorf("error while bundling with esbuild")
	}

	return stdout, stderr, nil
}

func ExecuteWasmtime(tmpFolderDir string) (string, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Make sure to have js.wasm in /root/untrusted-code-exec/js.wasm and wasmtime installed
	wasmtimeCmd := exec.CommandContext(ctx, "wasmtime", "run", "--dir=.", "--", "/root/untrusted-code-exec/js.wasm", "dist/bundle.js")
	wasmtimeCmd.Dir = tmpFolderDir

	stdout, stderr, err := RunCommandWithOutput(wasmtimeCmd)
	if ctx.Err() == context.DeadlineExceeded {
		return "", "Time limit exceeded", fmt.Errorf("wasm timed out after 10s")
	}

	if err != nil {
		fmt.Printf("wasmtime failed with error: %s\n", err)
		fmt.Printf("wasmtime stdout: %s\n", stdout)
		fmt.Printf("wasmtime stderr: %s\n", stderr)
		log.Printf("wasmtime failed: %s", err)
		return "", "", fmt.Errorf("failed to execute code: %s", err)
	}

	return string(stdout), string(stderr), nil
}

func CleanupTempNpmPackage(tmpFolderDir string) error {
	err := os.RemoveAll(tmpFolderDir)
	if err != nil {
		log.Printf("Failed to delete temp folder: %s", err)
		return fmt.Errorf("failed to delete temp folder")
	}

	return nil
}

func RunCommandWithOutput(cmd *exec.Cmd) (string, string, error) {
	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf

	err := cmd.Run()
	return stdoutBuf.String(), stderrBuf.String(), err
}

func copyDirectory(srcDir string, destDir string) error {
	err := filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(srcDir, path)
		if err != nil {
			return err
		}

		destPath := filepath.Join(destDir, relPath)

		if info.IsDir() {
			return os.MkdirAll(destPath, os.ModePerm)
		}

		return copyFile(path, destPath)
	})
	return err
}

func copyFile(src, dest string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return err
	}

	return destFile.Sync()
}
