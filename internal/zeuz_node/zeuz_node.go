package zeuz_node

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/automationsolutionz/zeuz_node/internal/zeuz_node/config"
	"golang.org/x/mod/semver"
)

const (
	betaGithubUrl          = "https://github.com/AutomationSolutionz/Zeuz_Python_Node/archive/refs/heads/beta.zip"
	githubReleasesEndpoint = "https://api.github.com/repos/AutomationSolutionz/Zeuz_Python_Node/releases"
)

// getZeuzNode downloads and makes ZeuZ Node available in the current directory
// if not already available. Right now, this will only check whether there's a
// `zeuz_node` directory present in the current location.
//
// TODO: This should ideally check for a config file - which version of node to
// download or use.
func getZeuzNode(zeuzNodeDir, payloadDir, url string, updateAvailable bool) {
	_, err := os.Stat(zeuzNodeDir)
	if !updateAvailable && err == nil {
		log.Printf("found zeuz node at: %v", zeuzNodeDir)
		return
	}

	os.RemoveAll(zeuzNodeDir)

	log.Println("downloading ZeuZ Node")

	os.MkdirAll(payloadDir, os.ModePerm)

	downloadPath := filepath.Join(payloadDir, "zeuz_node_download.zip")
	os.Remove(downloadPath)

	out, err := os.Create(downloadPath)
	if err != nil {
		log.Fatalf("failed to create output zip file for downloading zeuz node: %v", err)
	}
	defer out.Close()
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("failed to connect to internet to download zeuz node: %v", err)
	}
	defer resp.Body.Close()
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		log.Fatalf("failed to download zeuz node: %v", err)
	}

	_, err = Unzip(payloadDir, out.Name())

	// name of the dir inside the zip is of the form: username-repo_name-commit_hash
	var extractPath string
	dirEntries, err := os.ReadDir(payloadDir)
	for _, ent := range dirEntries {
		if strings.HasPrefix(ent.Name(), "AutomationSolutionz") {
			extractPath = filepath.Join(payloadDir, ent.Name())
			break
		}
	}
	defer os.RemoveAll(extractPath)

	log.Printf("extract path: %v\nzeuz node dir: %v", extractPath, zeuzNodeDir)
	if err = os.Rename(extractPath, zeuzNodeDir); err != nil {
		log.Fatalf("failed to move zeuz node from `%v` to `%v` with error: %v", extractPath, zeuzNodeDir, err)
	}
	os.Remove(downloadPath)
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

type githubRelease struct {
	Name       string `json:"name"`
	ZipballUrl string `json:"zipball_url"`
}

// fetchLatestVersionInfo returns whether we have the latest zeuz node installed by
// inspecting a config file.
func fetchLatestVersionInfo(paths config.Paths, conf config.Config) (string, bool) {
	log.Println("checking for new updates")

	resp, err := http.Get(githubReleasesEndpoint)
	if err != nil {
		log.Printf("failed to get the latest releases from github: %v", err)
		return "", false
	}
	defer resp.Body.Close()

	var releases []githubRelease
	err = json.NewDecoder(resp.Body).Decode(&releases)
	if err != nil {
		log.Printf("failed to get the latest releases from github: %v", err)
		return "", false
	}

	// Add 'v' prefix if not present, e.g 'v1.2.3'
	for i := 0; i < len(releases); i++ {
		if len(releases[i].Name) > 0 && releases[i].Name[0] != 'v' {
			releases[i].Name = fmt.Sprintf("v%s", releases[i].Name)
		} else {
			continue
		}
	}

	var gr *githubRelease

	// If we've never downloaded zeuz node before, we'll fetch the latest
	// version and download it.
	if conf.CurrentVersion == config.FirstRunVersion {
		for i := 0; i < len(releases); i++ {
			r := releases[i]

			// If release version is less or equal to the one locally present,
			// we skip, otherwise we store it as the new version.
			if semver.Compare(r.Name, conf.CurrentVersion) <= 0 {
				continue
			}
			conf.CurrentVersion = r.Name
			gr = &r
		}

		// write to config file.
		conf.WriteToFile(paths.ConfigPath)

		return gr.ZipballUrl, true
	}

	largestVersion := conf.CurrentVersion
	for i := 0; i < len(releases); i++ {
		r := releases[i]
		if semver.Compare(largestVersion, r.Name) < 0 {
			largestVersion = r.Name
			gr = &r
		}
	}

	if gr != nil {
		log.Printf(
			"new update for ZeuZ Node is available. Current version: %v, Latest version: %v",
			conf.CurrentVersion,
			gr.Name,
		)

		// write to config file.
		conf.CurrentVersion = gr.Name
		conf.WriteToFile(paths.ConfigPath)

		return gr.ZipballUrl, true
	}

	log.Println("we're already running the latest version")
	return "", false
}

// VerifyAndLaunchZeuzNode updates to latest zeuz node if not already available
// on the local machine and then launches it.
func VerifyAndLaunchZeuzNode(paths config.Paths, conf config.Config) {
	zipUrl := betaGithubUrl
	updateUrl, newUpdateAvailable := fetchLatestVersionInfo(paths, conf)
	if newUpdateAvailable {
		zipUrl = updateUrl
	}

	getZeuzNode(
		paths.ZeuzNodeDir,
		paths.ZeuzPayloadDir,
		zipUrl,
		newUpdateAvailable,
	)
	launchZeuzNode(paths.PythonPath, paths.ZeuzNodeDir, paths.ZeuzLogDir)
}
