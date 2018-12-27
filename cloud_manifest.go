package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
)

type CloudManifest struct {
	Updated  int64                  `json:"updated"`
	Versions []CloudManifestVersion `json:"versions"`
}

type CloudManifestVersion struct {
	Builds []CloudManifestBuild `json:"builds"`
	Name   string               `json:"name"`
}

type CloudManifestBuild struct {
	Architecture string   `json:"architecture"`
	GitVersion   string   `json:"gitVersion"`
	Platform     string   `json:"platform"`
	URL          string   `json:"url"`
	Flavor       string   `json:"flavor,omitempty"`
	MaxOsVersion *string  `json:"maxOsVersion,omitempty"`
	MinOsVersion *string  `json:"minOsVersion,omitempty"`
	Modules      []string `json:"modules,omitempty"`

	Win2008Plus        bool     `json:"win2008plus,omitempty"`
	WinVCRedistDll     string   `json:"winVCRedistDll,omitempty"`
	WinVCRedistOptions []string `json:"winVCRedistOptions,omitempty"`
	WinVCRedistURL     string   `json:"winVCRedistUrl,omitempty"`
	WinVCRedistVersion string   `json:"winVCRedistVersion,omitempty"`
}

const (
	cloudVersionManifest   = "https://raw.githubusercontent.com/10gen/mms/dbc4bd8b8fb0002ddaac2f45a7a6e239af5c8f60/server/conf/mongodb_version_manifest.json"
	cloudVersionManifest36 = "https://raw.githubusercontent.com/10gen/mms/dbc4bd8b8fb0002ddaac2f45a7a6e239af5c8f60/server/src/webapp-mms/static/version_manifest/3.6.json"

	// Pre 3.6.9 added ok?
	pre369Release          = "8edf6d59e1bfc7a0ab193208d6fae61debd06221"
	cloudVersionManifest40 = "https://s3.amazonaws.com/om-kubernetes-conf/4.0.json"
	// cloudVersionManifest40 = "https://raw.githubusercontent.com/10gen/mms/" + pre369Release + "/server/src/webapp-mms/static/version_manifest/4.0.json"

	// The 3.6.9 version was added in this commit dbd14a6db330428b472396860ab80a697a0afdd5, so the next URL
	// points to the version of the 4.0.json file which has been updated with this version
	cloudVersionManifest40Post369 = "https://raw.githubusercontent.com/10gen/mms/dbd14a6db330428b472396860ab80a697a0afdd5/server/src/webapp-mms/static/version_manifest/4.0.json"
)

func fetchCloudVersionManifest(versionManifestFile string) (*CloudManifest, error) {
	if strings.HasPrefix(versionManifestFile, "https://") {
		return fetchCloudVersionManifestFromURL(versionManifestFile)
	}
	return fetchCloudVersionManifestFromFile(versionManifestFile)
}

func fetchCloudVersionManifestFromURL(versionManifestFile string) (*CloudManifest, error) {
	res, err := http.Get(versionManifestFile)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	s := &CloudManifest{}
	err = json.Unmarshal(body, &s)
	if err != nil {
		return nil, err
	}

	return s, nil
}

func fetchCloudVersionManifestFromFile(versionManifestFile string) (*CloudManifest, error) {
	body, err := ioutil.ReadFile(versionManifestFile)
	if err != nil {
		return nil, err
	}

	s := &CloudManifest{}
	err = json.Unmarshal(body, &s)
	if err != nil {
		return nil, err
	}

	return s, nil
}
