package zeuz_node

import (
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

// getZeuzNode downloads and makes ZeuZ Node available in the current directory
// if not already available. Right now, this will only check whether there's a
// `zeuz_node` directory present in the current location.
//
// TODO: This should ideally check for a config file - which version of node to
// download or use.
func getZeuzNode(zeuzNodeDir, payloadDir string) {
	_, err := os.Stat(zeuzNodeDir)
	if err == nil {
		log.Printf("found zeuz node at: %v", zeuzNodeDir)
		return
	}

	log.Println("downloading ZeuZ Node")

	os.MkdirAll(payloadDir, os.ModePerm)

	downloadPath := filepath.Join(payloadDir, "zeuz_node_download.zip")
	os.Remove(downloadPath)
	extractPath := filepath.Join(payloadDir, "Zeuz_Python_Node-beta")
	os.RemoveAll(extractPath)

	out, err := os.Create(downloadPath)
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

	_, err = Unzip(payloadDir, out.Name())
	log.Printf("extract path: %v\nzeuz node dir: %v", extractPath, zeuzNodeDir)
	if err = os.Rename(extractPath, zeuzNodeDir); err != nil {
		log.Fatalf("failed to move zeuz node from `%v` to `%v` with error: %v", extractPath, zeuzNodeDir, err)
	}
}

// launchZeuzNode launches `node_cli.py` with the log directory set to `qlogs` in
// the current directory.
func launchZeuzNode(pythonPath, zeuzNodePath, logDir string) {
	nodeCliPath := filepath.Join(zeuzNodePath, "node_cli.py")

	err := os.MkdirAll(logDir, os.ModePerm)
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

// VerifyAndLaunchZeuzNode updates to latest zeuz node if not already available
// on the local machine and then launches it.
func VerifyAndLaunchZeuzNode(pythonPath, zeuzNodeDir, payloadDir, logDir string) {
	getZeuzNode(zeuzNodeDir, payloadDir)
	launchZeuzNode(pythonPath, zeuzNodeDir, logDir)
}
