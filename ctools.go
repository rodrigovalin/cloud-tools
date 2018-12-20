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
	"fmt"
	"os"
)

const (
	serverVersionManifest  = "http://downloads.mongodb.org.s3.amazonaws.com/full.json"
	cloudVersionManifest   = "https://raw.githubusercontent.com/10gen/mms/master/server/conf/mongodb_version_manifest.json"
	cloudVersionManifest36 = "https://raw.githubusercontent.com/10gen/mms/master/server/src/webapp-mms/static/version_manifest/3.6.json"
	cloudVersionManifest40 = "https://raw.githubusercontent.com/10gen/mms/master/server/src/webapp-mms/static/version_manifest/4.0.json"
)

func main() {
	fmt.Println("Hola")

	fmt.Println("Fetching server manifest")

	token := getGitHubToken()
	if token == "" {
		fmt.Println("Please configure GITHUB_TOKEN")
		os.Exit(1)
	}

	serverManifest, err := fetchServerVersionManifest(token)
	if err != nil {
		fmt.Println("Error Fetching the server manifest")
	} else {
		serverVersions := make([]string, len(serverManifest.Versions))
		for _, version := range serverManifest.Versions {
			serverVersions = append(serverVersions, version.Version)
		}
		fmt.Printf("There are %d versions in the server manifest\n", len(serverVersions))
	}

	fmt.Println("Fetching cloud manifest")
	cloudManifest, err := fetchCloudVersionManifest(token)
	if err != nil {
		fmt.Println("Error Fetching the cloud manifest")
	} else {
		cloudVersions := make([]string, len(cloudManifest.Versions))
		for _, version := range cloudManifest.Versions {
			cloudVersions = append(cloudVersions, version.Name)
		}
		fmt.Printf("There are %d versions in the server manifest\n", len(cloudVersions))
	}
}

func getGitHubToken() string {
	token := os.Getenv("GITHUB_TOKEN")
	if token != "" {
		return token
	}

	// TODO: look for the token on a file in home dir (maybe ~.mci/)
	return ""
}
