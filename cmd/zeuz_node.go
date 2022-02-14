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

	// ~/zeuz/payload
	// payloadDir holds the downloads and extract dir
	payloadDir := filepath.Join(zeuzRootDir, "payload")

	// ~/zeuz/zeuz_node_logs
	// logDir centralizes all zeuz node logs
	logDir := filepath.Join(zeuzRootDir, "zeuz_node_logs")

	// ~/zeuz/zeuz_python_node
	nodeDir := filepath.Join(zeuzRootDir, "zeuz_node_python")

	pythonPath, err := python.VerifyAndInstallPython(payloadDir)
	if err != nil {
		return
	}
	zeuz_node.VerifyAndLaunchZeuzNode(pythonPath, nodeDir, payloadDir, logDir)

	log.Println("done. Exiting")
}
