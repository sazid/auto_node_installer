// Author: Mohammed Sazid Al Rashid
// GitHub: sazid

package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/automationsolutionz/zeuz_node/internal/zeuz_node"
	"github.com/automationsolutionz/zeuz_node/internal/zeuz_node/python"
)

var (
	customLocation *string = flag.String("location", "", "specify a custom location where zeuz node will store all its data")
)

func main() {
	log.Println("starting ZeuZ Node")

	flag.Parse()

	// setup all the required paths
	homeDir := ""
	log.Println(*customLocation)
	if customLocation == nil || len(*customLocation) == 0 {
		// ~/
		var err error
		homeDir, err = os.UserHomeDir()
		if err != nil {
			log.Fatalf("failed to get current user's home directory: %v", err)
		}
	} else {
		err := os.MkdirAll(*customLocation, os.ModeDir)
		if err != nil {
			log.Fatalf("failed to create or locate custom location: %v", err)
		}
		homeDir = *customLocation
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
