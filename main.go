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
)

const serverVersionManifest = "http://downloads.mongodb.org.s3.amazonaws.com/full.json"
const cloudVersionManifest = "https://raw.githubusercontent.com/10gen/mms/master/server/conf/mongodb_version_manifest.json?token=AAHlg1Cl2VYanZz7AaHQCRGKP2GEBCHTks5cI7RbwA%3D%3D"
const cloudVersionManifest36 = "https://raw.githubusercontent.com/10gen/mms/master/server/src/webapp-mms/static/version_manifest/3.6.json?token=AAHlg-Er1Hgv1Cm5_aoI3GMwBiPoeSYOks5cI8F-wA%3D%3D"
const cloudVersionManifest40 = "https://raw.githubusercontent.com/10gen/mms/master/server/src/webapp-mms/static/version_manifest/4.0.json?token=AAHlg2upKuFsW6_iIZYu8txzbDSfxOgTks5cI8HmwA%3D%3D"

func main() {
	fmt.Println("Hola")

	fmt.Println("Fetching server manifest")

	serverManifest, err := fetchServerVersionManifest()
	if err != nil {
		fmt.Println("Error Fetching the server manifest")
	} else {
		serverVersions := make([]string, len(serverManifest.Versions))
		for _, version := range serverManifest.Versions {
			serverVersions = append(serverVersions, version.Version)
		}
		fmt.Printf("There are %d versions in the server manifest\n", len(serverVersions))
		// fmt.Printf("%+v\n", serverManifest.Versions[0])
	}

	fmt.Println("Fetching cloud manifest")
	cloudManifest, err := fetchCloudVersionManifest()
	if err != nil {
		fmt.Println("Error Fetching the cloud manifest")
	} else {
		cloudVersions := make([]string, len(cloudManifest.Versions))
		for _, version := range cloudManifest.Versions {
			cloudVersions = append(cloudVersions, version.Name)
		}
		fmt.Printf("There are %d versions in the server manifest\n", len(cloudVersions))
		// fmt.Printf("%+v\n", cloudManifest.Versions[0])
	}
}
