package python

import (
	"embed"
	"errors"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/automationsolutionz/zeuz_node/internal/zeuz_node"
)

const (
	pythonInstallerFilename = "python-3.8.10-amd64.exe"
)

//go:embed embed/*
var embeddedFiles embed.FS

var (
	ErrPythonNotFound = errors.New("failed to find python in PATH or PATHEXT")
)

// fromCwd returns the `name` joined with the current working directory
func fromCwd(name string) string {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("failed to get current working directory: %v", err)
	}
	return filepath.Join(cwd, name)
}

// isPythonInstalled looks for any previous python installations by observing
// the PATH or PATHEXT. If not found, it'll try to see if an adjacent local
// python installation previously done by this program is available.
func isPythonInstalled() (pythonPath string, found bool) {
	lookFor := []string{
		"python3",
		"python",
		filepath.Join(fromCwd("python"), "python3"),
		filepath.Join(fromCwd("python"), "python"),
	}

	for _, loc := range lookFor {
		pythonPath, err := exec.LookPath(loc)
		if err == nil {
			return pythonPath, true
		}
	}

	return "", false
}

// VerifyAndInstallPython verifies whether a python installation is already
// available, if not it'll auto install the Python (only for Windows). It
// returns the path to the `python` executable.
func VerifyAndInstallPython(payloadDir string) (string, error) {
	pythonPath, found := isPythonInstalled()
	if !found {
		zeuz_node.ExtractFiles(embeddedFiles, payloadDir)
		installPython(payloadDir, pythonInstallerFilename)

		pythonPath, found = isPythonInstalled()
		if !found {
			log.Println("failed to find python in PATH or PATHEXT after installing python")
			return "", ErrPythonNotFound
		}
	}
	if !found {
		log.Println("failed to find python in PATH or PATHEXT")
		return "", ErrPythonNotFound
	}
	log.Printf("found python at: %v", pythonPath)

	return pythonPath, nil
}