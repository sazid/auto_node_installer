// Author: Mohammed Sazid Al Rashid
// GitHub: sazid

package main

import (
	"bytes"
	"flag"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/automationsolutionz/zeuz_node/internal/zeuz_node"
	"github.com/automationsolutionz/zeuz_node/internal/zeuz_node/config"
	"github.com/automationsolutionz/zeuz_node/internal/zeuz_node/python"
)

var (
	customLocation *string = flag.String("location", "", "specify a custom location where zeuz node will store all its data")
)

// checkWriteable checks to see if the selected location is writable - if not,
// we'll use the current working directory.
func checkWriteable(homeDir string) string {
	testDir := filepath.Join(homeDir, "zeuzwritetestdir")
	err := os.Mkdir(testDir, os.ModePerm)
	if err != nil {
		homeDir, err = os.Getwd()
		if err != nil {
			log.Fatalf("failed to get current working directory: %v", err)
		}
	}
	os.RemoveAll(testDir)

	return homeDir
}

func handleSignals() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Printf("%s received, exiting...", <-sigs)
		os.Exit(0)
	}()
}

func main() {
	handleSignals()

	log.Println("starting ZeuZ Node")

	flag.Parse()

	// setup all the required paths
	homeDir := ""
	if customLocation == nil || len(*customLocation) == 0 {
		// ~/
		var err error
		homeDir, err = os.UserHomeDir()
		if err != nil {
			log.Fatalf("failed to get current user's home directory: %v", err)
		}
	} else {
		err := os.MkdirAll(*customLocation, os.ModePerm)
		if err != nil {
			log.Fatalf("failed to create zeuz dir at the specified custom location: %v", err)
		}
		homeDir = *customLocation
	}

	homeDir = checkWriteable(homeDir)

	// ~/zeuz
	zeuzRootDir := filepath.Join(homeDir, "zeuz")

	// ~/zeuz/python
	// default path where we automatically install python to.
	defaultPythonInstallDir := filepath.Join(zeuzRootDir, "python")

	// ~/zeuz/payload
	// payloadDir holds the temporary downloads and extract dir
	payloadDir := filepath.Join(zeuzRootDir, "payload")

	// cleanup payloadDir after we're done as it contains transient data
	// TODO: This is not working right now because of os.Exit(1)
	defer os.RemoveAll(payloadDir)

	// ~/zeuz/zeuz_node_logs
	// logDir centralizes all zeuz node logs
	logDir := filepath.Join(zeuzRootDir, "zeuz_node_logs")

	// ~/zeuz/zeuz_node_python
	nodeDir := filepath.Join(zeuzRootDir, "zeuz_node_python")

	configPath := filepath.Join(zeuzRootDir, "config.json")

	paths := config.Paths{
		HomeDir:                 homeDir,
		WorkingDir:              zeuzRootDir,
		ZeuzNodeDir:             nodeDir,
		ZeuzLogDir:              logDir,
		ZeuzPayloadDir:          payloadDir,
		DefaultPythonInstallDir: defaultPythonInstallDir,
		ConfigPath:              configPath,
	}

	os.MkdirAll(paths.WorkingDir, os.ModePerm)

	var err error
	paths.PythonPath, err = python.VerifyAndInstallPython(paths)
	if err != nil {
		defer os.Exit(1)
		return
	}

	var conf config.Config
	f, err := os.Open(paths.ConfigPath)
	if err != nil {
		log.Println("no previous config file found, using the default config.")
		conf, err = config.NewConfig(bytes.NewBufferString(config.DefaultConfig))
	} else {
		defer f.Close()

		conf, err = config.NewConfig(f)
		if err != nil {
			log.Fatalf("failed to read config file")
		}
	}

	zeuz_node.VerifyAndLaunchZeuzNode(paths, conf)

	log.Println("done. Exiting")
}
