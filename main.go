package main

import (
	"embed"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
)

const (
	pythonInstallerFilename = "python-3.8.10-amd64.exe"
	refreshEnvFilename      = "RefreshEnv.cmd"
	payloadDir              = "payload"
)

//go:embed embed/*
var embeddedFiles embed.FS

// extractFiles extracts all the embedded files into the "payload" directory.
func extractFiles() {
	dir := payloadDir

	err := os.MkdirAll(dir, os.ModeDir)
	if err != nil {
		log.Fatalf("failed to create `%s` directory to extract required files", dir)
	}

	// extract all the embedded files
	efDirEntries, err := embeddedFiles.ReadDir("embed")
	if err != nil {
		log.Fatal("failed to read embedded files")
	}
	for _, efDirEntry := range efDirEntries {
		// creating an anonymous func() so we can use `defer` to properly close
		// any resources.
		func() {
			log.Printf("extracting `%s` into `%s` directory", efDirEntry.Name(), dir)
			embeddedFile, err := embeddedFiles.Open(path.Join("embed", efDirEntry.Name()))
			if err != nil {
				log.Fatalf("failed to open embedded file: %s", efDirEntry.Name())
			}
			defer embeddedFile.Close()

			osFile, err := os.OpenFile(
				filepath.Join(dir, efDirEntry.Name()),
				os.O_RDWR|os.O_CREATE|os.O_TRUNC,
				0644,
			)
			if err != nil {
				log.Fatalf("failed to open file: %s", efDirEntry.Name())
			}
			defer osFile.Close()

			_, err = io.Copy(osFile, embeddedFile)
			if err != nil {
				log.Fatalf("failed to extract file `%s` to `%s`", efDirEntry.Name(), osFile.Name())
			}
		}()
	}
}

// isPythonInstalled looks for any previous python installations by observing
// the PATH or PATHEXT. If not found, it'll try to see if an adjacent local
// python installation previously done by this program is available.
func isPythonInstalled() (pythonPath string, found bool) {
	lookFor := []string{
		"python3",
		"python",
		filepath.Join(getLocalPythonInstallDir(), "python3"),
		filepath.Join(getLocalPythonInstallDir(), "python"),
	}

	for _, loc := range lookFor {
		pythonPath, err := exec.LookPath(loc)
		if err == nil {
			return pythonPath, true
		}
	}

	return "", false
}

func getLocalPythonInstallDir() string {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("failed to get current working directory: %v", err)
	}
	return filepath.Join(cwd, "python")
}

// installPython installs the extracted Python installation file.
func installPython() {
	log.Println("installing python")

	pythonInstaller := filepath.Join(payloadDir, pythonInstallerFilename)

	installerArgs := []string{
		"/passive",
		"InstallAllUsers=0",
		"Include_launcher=0",
		"Include_test=0",
		"Include_launcher=1",
		"Include_pip=1",
		"Include_tools=1",
		"Include_exe=1",
		"Include_dev=1",
		"PrependPath=1",
		"AssociateFiles=1",
		fmt.Sprintf("DefaultJustForMeTargetDir=%s", getLocalPythonInstallDir()),
	}
	cmd := exec.Command(pythonInstaller, installerArgs...)

	err := cmd.Run()
	if err != nil {
		log.Fatalf("failed to install python: %v", err)
	}
}

// downloadZeuzNode downloads and makes ZeuZ Node available in the current
// directory. Right now, this will only check whether there's a `zeuz_node`
// directory present in the current location.
//
// TODO: This should ideally check for a config file - which version of node to
// download or use.
func downloadZeuzNode() (zeuzNodePath string) {
	log.Println("")
	log.Println("downloading ZeuZ Node")
}

// launchZeuzNode launches `node_cli.py` with the log directory set to `qlogs` in
// the current directory.
func launchZeuzNode(pythonPath, zeuzNodePath string) {

}

func main() {
	log.Println("starting ZeuZ Node")
	extractFiles()
	pythonPath, found := isPythonInstalled()
	if !found {
		installPython()

		pythonPath, found = isPythonInstalled()
		if !found {
			log.Fatal("failed to find python in PATH or PATHEXT after installing python")
		}
	}
	if !found {
		log.Fatal("failed to find python in PATH or PATHEXT")
	}
	log.Printf("found python at: %v", pythonPath)
	zeuzNodePath := downloadZeuzNode()
	launchZeuzNode(pythonPath, zeuzNodePath)
}
