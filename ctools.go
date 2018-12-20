package main

/*
This is the easiest program ever written in Go. It has the follow functionality:

1. Read from the server builds json
2. Read from the cloud builds json
3. Compare the two
4. Spit out the set differences (each element a build)

Then

1. Authenticate to Jira
2. Read the release ticket (https://jira.mongodb.org/browse/CLOUDP-35176)
3. Generate a new cloud json file combining:
  - the old server team json man -> http://downloads.mongodb.org.s3.amazonaws.com/full.json
  - the old cloud json man -> https://github.com/10gen/mms/blob/master/server/conf/mongodb_version_manifest.json
  - the information in the ticket about the new build to be added
  -
4. Write this new file into disk

Then

1. Make sure the new file is correct? Simple validations
  - The file should be json valid
  - All of the builds should be HEAD-able
2. Make sure the new builds are downloadable and that the SHA checks
3.
*/

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

const (
	serverVersionManifest  = "http://downloads.mongodb.org.s3.amazonaws.com/full.json"
	cloudVersionManifest   = "https://raw.githubusercontent.com/10gen/mms/master/server/conf/mongodb_version_manifest.json"
	cloudVersionManifest36 = "https://raw.githubusercontent.com/10gen/mms/master/server/src/webapp-mms/static/version_manifest/3.6.json"
	cloudVersionManifest40 = "https://raw.githubusercontent.com/10gen/mms/master/server/src/webapp-mms/static/version_manifest/4.0.json"

	WinVCRedistDll34 = "vcruntime140.dll"
	WinVCRedistUrl34 = "http://download.microsoft.com/download/6/D/F/6DF3FF94-F7F9-4F0B-838C-A328D1A7D0EE/vc_redist.x64.exe"

	// TODO add this
	WinVCRedistDll = "msvcr120.dll"
	WinVCRedistUrl = "http://download.microsoft.com/download/2/E/6/2E61CFA4-993B-4DD4-91DA-3737CD5CD6E3/vcredist_x64.exe"

	WinVCRedistVersion   = "10.0.40219.325"
	WinVCRedistVersion3  = "12.0.21005.1"
	WinVCRedistVersion34 = "14.0.24212.0"
)

var FlavorLinux = [...]string{"suse", "rhel", "ubuntu", "debian", "amazon"}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: ctools add-version <version>")
		os.Exit(1)
	}

	versionArg := os.Args[2]

	token := getGitHubToken()
	if token == "" {
		fmt.Println("Please configure GITHUB_TOKEN")
		os.Exit(1)
	}

	serverManifest, err := fetchServerVersionManifest(token)
	if err != nil {
		fmt.Println("Error Fetching the server manifest")
		os.Exit(1)
	}

	cloudManifest, err := buildCloudManifestForVersion(versionArg, serverManifest)
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
	manifest, err := json.Marshal(cloudManifest)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(string(manifest))
}

func buildCloudManifestForVersion(newVersion string, server *ServerManifest) (*CloudManifest, error) {
	cloudManifest := &CloudManifest{Updated: time.Now().Unix()}
	for _, version := range server.Versions {
		if version.Version == newVersion {
			builds := buildBuildsForCloudManifestVersion(version)
			cloudManifest.Versions = []CloudManifestVersion{{
				Builds: builds,
				Name:   newVersion,
			}}
		}
	}

	return cloudManifest, nil
}

func newCloudManifestVersion(version, gitsha string, serverVersion ServerManifestVersion) (*CloudManifestVersion, error) {
	builds := buildBuildsForCloudManifestVersion(serverVersion)
	return &CloudManifestVersion{Name: version, Builds: builds}, nil
}

func buildBuildsForCloudManifestVersion(serverVersion ServerManifestVersion) []CloudManifestBuild {
	cloudManifestBuilds := make([]CloudManifestBuild, len(serverVersion.Downloads))
	for _, download := range serverVersion.Downloads {
		build := CloudManifestBuild{
			Architecture: download.Arch,
			GitVersion:   serverVersion.Githash,
			Platform:     getPlatformFromTarget(download.Target),
			URL:          getPartialFromFullURL(download.Archive.URL),
		}

		if targetIsLinux(download.Target) {
			build.Flavor = getFlavorFromTarget(download.Target)
		}

		if targetIsWindows(download.Target) {
			// TODO: lots of rules for windows specially different versions of the redistdll
			dll, url := getWinRCRedistDll(serverVersion.Version)
			if strings.Contains(download.Target, "2008plus") {
				build.Win2008Plus = true
				build.WinVCRedistDll = dll
				build.WinVCRedistOptions = []string{"/quiet", "/norestart"}
				build.WinVCRedistURL = url
			}
		}
		if targetIsMacOS(download.Target) {
			// nothing actually
		}

		cloudManifestBuilds = append(cloudManifestBuilds, build)
	}

	return cloudManifestBuilds
}

func getGitHubToken() string {
	token := os.Getenv("GITHUB_TOKEN")
	if token != "" {
		return token
	}

	// TODO: look for the token on a file in home dir (maybe ~.mci/)
	return ""
}

func getPlatformFromTarget(target string) string {
	if targetIsLinux(target) {
		return "linux"
	}

	if targetIsMacOS(target) {
		return "macos"
	}

	if targetIsWindows(target) {
		return "windows"
	}
	return ""
}

func targetIsLinux(target string) bool {
	for _, el := range FlavorLinux {
		if strings.Contains(target, el) {
			return true
		}
	}
	return false
}

func targetIsMacOS(target string) bool {
	return strings.Contains(target, "macos") || strings.Contains(target, "osx")
}

func targetIsWindows(target string) bool {
	return strings.Contains(target, "windows")
}

func getPartialFromFullURL(full string) string {
	splited := strings.Split(full, "/")
	return strings.Join(splited[len(splited)-2:], "/")
}

func getFlavorFromTarget(target string) string {
	for _, flavor := range FlavorLinux {
		if strings.Contains(target, flavor) {
			return flavor
		}
	}
	return ""
}

func getWinRCRedistDll(version string) (string, string) {
	if strings.HasPrefix(version, "3.4") ||
		strings.HasPrefix(version, "3.6") ||
		strings.HasPrefix(version, "4.") {
		return WinVCRedistDll34, WinVCRedistUrl34
	}

	return WinVCRedistDll, WinVCRedistUrl
}
