// Author: Mohammed Sazid Al Rashid
// GitHub: sazid

package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/automationsolutionz/zeuz_node/internal/zeuz_node"
	"github.com/automationsolutionz/zeuz_node/internal/zeuz_node/python"
)

func main() {
	log.Println("starting ZeuZ Node")

	// setup all the required paths
	// ~/
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("failed to get current user's home directory: %v", err)
	}

	// ~/zeuz
	zeuzRootDir := filepath.Join(homeDir, "zeuz")

	// ~/zeuz/python
	// default path where we automatically install python to.
	defaultPythonInstallDir := filepath.Join(zeuzRootDir, "python")

	// ~/zeuz/payload
	// payloadDir holds the temporary downloads and extract dir
	payloadDir := filepath.Join(zeuzRootDir, "payload")

	// cleanup payloadDir after we're done as it contains transient data
	defer os.RemoveAll(payloadDir)

	// ~/zeuz/zeuz_node_logs
	// logDir centralizes all zeuz node logs
	logDir := filepath.Join(zeuzRootDir, "zeuz_node_logs")

	// ~/zeuz/zeuz_node_python
	nodeDir := filepath.Join(zeuzRootDir, "zeuz_node_python")

	pythonPath, err := python.VerifyAndInstallPython(payloadDir, defaultPythonInstallDir)
	if err != nil {
		defer os.Exit(1)
		return
	}
	zeuz_node.VerifyAndLaunchZeuzNode(pythonPath, nodeDir, payloadDir, logDir)

	log.Println("done. Exiting")
}
