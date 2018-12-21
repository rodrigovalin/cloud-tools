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
	serverVersionManifest = "http://downloads.mongodb.org.s3.amazonaws.com/full.json"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Printf("Usage:\n\t%s add-version <version>\n", os.Args[0])
		os.Exit(1)
	}

	versionArg := os.Args[2]

	serverManifest, err := fetchServerVersionManifest()
	if err != nil {
		fmt.Println("Error Fetching the server manifest")
		os.Exit(1)
	}

	cloudManifest, err := fetchCloudVersionManifest(os.Getenv("GITHUB_TOKEN"))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	updatedCloudManifest, err := buildCloudManifestForVersion(versionArg, serverManifest)
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}

	cloudManifest.Updated = updatedCloudManifest.Updated
	cloudManifest.Versions = append(cloudManifest.Versions, updatedCloudManifest.Versions...)
	manifest, err := json.Marshal(cloudManifest)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(string(manifest))
}

func buildCloudManifestForVersion(newVersion string, server *ServerManifest) (*CloudManifest, error) {
	cloudManifest := &CloudManifest{Updated: time.Now().Unix() * 1000}
	for _, version := range server.Versions {
		if version.Version == newVersion {
			community, enterprise := buildBuildsForCloudManifestVersion(version)
			cloudManifest.Versions = []CloudManifestVersion{{
				Builds: community,
				Name:   newVersion,
			}, {
				Builds: enterprise,
				Name:   newVersion + "-ent",
			}}
		}
	}

	return cloudManifest, nil
}

func buildBuildsForCloudManifestVersion(serverVersion ServerManifestVersion) ([]CloudManifestBuild, []CloudManifestBuild) {
	cloudManifestBuilds := make([]CloudManifestBuild, 0)
	cloudManifestBuildsEnt := make([]CloudManifestBuild, 0)

	for _, download := range serverVersion.Downloads {
		if shouldSkipDownload(&download) {
			continue
		}

		build := CloudManifestBuild{
			Architecture: getCloudArchFromServerArch(download.Arch),
			GitVersion:   serverVersion.Githash,
			Platform:     getPlatformFromTarget(download.Target),
			URL:          getPartialFromFullURL(download.Archive.URL),
		}

		applyLinuxAttributes(serverVersion.Version, &download, &build)
		applyWindowsAttributes(serverVersion.Version, &download, &build)

		if download.Edition == "enterprise" {
			build.GitVersion = serverVersion.Githash + " modules: enterprise"
			build.Modules = []string{"enterprise"}
			cloudManifestBuildsEnt = append(cloudManifestBuildsEnt, build)
		} else {
			cloudManifestBuilds = append(cloudManifestBuilds, build)
		}
	}

	return cloudManifestBuilds, cloudManifestBuildsEnt
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

func getPartialFromFullURL(full string) string {
	splited := strings.Split(full, "/")
	return "/" + strings.Join(splited[len(splited)-2:], "/")
}

func getCloudArchFromServerArch(arch string) string {
	if arch == "x86_64" {
		return "amd64"
	}

	return arch
}

func shouldSkipDownload(download *ServerManifestDownload) bool {
	return download.Edition == "source" || download.Arch == "arm64"
}
