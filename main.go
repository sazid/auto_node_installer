// Author: Mohammed Sazid Al Rashid
// GitHub: sazid

package main

import (
	"archive/zip"
	"embed"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

const (
	pythonInstallerFilename = "python-3.8.10-amd64.exe"
	payloadDir              = "payload"
	activeZeuzNodeDir       = "zeuz_node-active"
)

var (
	logDir = fromCwd("logs")
)

//go:embed embed/*
var embeddedFiles embed.FS

// extractTo extracts the file specified by dirEntry into destDir
func extractTo(destDir string, dirEntry fs.DirEntry) {
	log.Printf("extracting `%s` into `%s` directory", dirEntry.Name(), destDir)
	embeddedFile, err := embeddedFiles.Open(path.Join("embed", dirEntry.Name()))
	if err != nil {
		log.Fatalf("failed to open embedded file: %s", dirEntry.Name())
	}
	defer embeddedFile.Close()

	osFile, err := os.Create(filepath.Join(destDir, dirEntry.Name()))
	if err != nil {
		log.Fatalf("failed to open file: %s", dirEntry.Name())
	}
	defer osFile.Close()

	_, err = io.Copy(osFile, embeddedFile)
	if err != nil {
		log.Fatalf("failed to extract file `%s` to `%s`", dirEntry.Name(), osFile.Name())
	}
}

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
		extractTo(dir, efDirEntry)
	}
}

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
		"Include_tcltk=1",
		"PrependPath=1",
		"AssociateFiles=1",
		fmt.Sprintf("DefaultJustForMeTargetDir=%s", fromCwd("python")),
	}
	cmd := exec.Command(pythonInstaller, installerArgs...)

	err := cmd.Run()
	if err != nil {
		log.Fatalf("failed to install python: %v", err)
	}
}

// getZeuzNode downloads and makes ZeuZ Node available in the current directory
// if not already available. Right now, this will only check whether there's a
// `zeuz_node` directory present in the current location.
//
// TODO: This should ideally check for a config file - which version of node to
// download or use.
func getZeuzNode() (zeuzNodeDir string) {
	zeuzNodeDir = fromCwd(activeZeuzNodeDir)
	_, err := os.Stat(zeuzNodeDir)
	if err == nil {
		log.Printf("found zeuz node at: %v", zeuzNodeDir)
		return
	}

	log.Println("downloading ZeuZ Node")

	out, err := os.Create("zeuz_node_download.zip")
	if err != nil {
		log.Fatalf("failed to create output zip file for downloading zeuz node: %v", err)
	}
	defer out.Close()
	resp, err := http.Get("https://github.com/AutomationSolutionz/Zeuz_Python_Node/archive/refs/heads/beta.zip")
	if err != nil {
		log.Fatalf("failed to connect to internet to download zeuz node: %v", err)
	}
	defer resp.Body.Close()
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		log.Fatalf("failed to download zeuz node: %v", err)
	}

	_, err = Unzip(out.Name(), fromCwd("."))
	os.RemoveAll(activeZeuzNodeDir)
	os.Rename("Zeuz_Python_Node-beta", activeZeuzNodeDir)
	return
}

// launchZeuzNode launches `node_cli.py` with the log directory set to `qlogs` in
// the current directory.
func launchZeuzNode(pythonPath, zeuzNodePath, logDir string) {
	nodeCliPath := filepath.Join(zeuzNodePath, "node_cli.py")

	err := os.MkdirAll(logDir, os.ModeDir)
	if err != nil {
		log.Fatalf("failed to create log directory: %v", err)
	}

	var nodeCliArgs = []string{
		nodeCliPath,
		"-d",
		logDir,
	}

	cmd := exec.Command(pythonPath, nodeCliArgs...)
	cmd.Dir = zeuzNodePath

	// Connect standard in/out devices.
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin

	err = cmd.Run()
	if err != nil {
		log.Fatalf("failed to launch zeuz node: %v", err)
	}
}

// Unzip will decompress a zip archive, moving all files and folders
// within the zip file (parameter 1) to an output directory (parameter 2).
func Unzip(src, dest string) ([]string, error) {

	var filenames []string

	r, err := zip.OpenReader(src)
	if err != nil {
		return filenames, err
	}
	defer r.Close()

	for _, f := range r.File {

		// Store filename/path for returning and using later on
		fpath := filepath.Join(dest, f.Name)

		// Check for ZipSlip. More Info: http://bit.ly/2MsjAWE
		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return filenames, fmt.Errorf("%s: illegal file path", fpath)
		}

		filenames = append(filenames, fpath)

		if f.FileInfo().IsDir() {
			// Make Folder
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		// Make File
		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return filenames, err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return filenames, err
		}

		rc, err := f.Open()
		if err != nil {
			return filenames, err
		}

		_, err = io.Copy(outFile, rc)

		// Close the file without defer to close before next iteration of loop
		outFile.Close()
		rc.Close()

		if err != nil {
			return filenames, err
		}
	}
	return filenames, nil
}

func main() {
	log.Println("starting ZeuZ Node")
	pythonPath, found := isPythonInstalled()
	if !found {
		extractFiles()
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
	zeuzNodeDir := getZeuzNode()
	launchZeuzNode(pythonPath, zeuzNodeDir, logDir)

	log.Println("done. Exiting")
}
