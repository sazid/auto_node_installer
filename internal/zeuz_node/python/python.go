package python

import (
	"embed"
	"errors"
	"log"
	"os/exec"
	"path/filepath"

	"github.com/automationsolutionz/zeuz_node/internal/zeuz_node"
	"github.com/automationsolutionz/zeuz_node/internal/zeuz_node/config"
)

const (
	pythonInstallerFilename = "python-3.8.10-amd64.exe"
)

//go:embed embed/*
var embeddedFiles embed.FS

var (
	ErrPythonNotFound = errors.New("failed to find python in PATH or PATHEXT")
)

// isPythonInstalled looks for any previous python installations by observing
// the PATH or PATHEXT. If not found, it'll try to see if an adjacent local
// python installation previously done by this program is available.
func isPythonInstalled(defaultPythonInstallDir string) (pythonPath string, found bool) {
	lookFor := []string{
		"python3",
		"python",
		filepath.Join(defaultPythonInstallDir, "python3"),
		filepath.Join(defaultPythonInstallDir, "python"),
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
func VerifyAndInstallPython(paths config.Paths) (string, error) {
	pythonPath, found := isPythonInstalled(paths.DefaultPythonInstallDir)
	if !found {
		zeuz_node.ExtractFiles(embeddedFiles, paths.ZeuzPayloadDir)

		pythonInstallerPath := filepath.Join(paths.ZeuzPayloadDir, pythonInstallerFilename)
		installPython(pythonInstallerPath, paths.DefaultPythonInstallDir)

		pythonPath, found = isPythonInstalled(paths.DefaultPythonInstallDir)
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
