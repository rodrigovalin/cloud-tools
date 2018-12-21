package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
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
	MaxOsVersion string   `json:"maxOsVersion,omitempty"`
	MinOsVersion string   `json:"minOsVersion,omitempty"`
	Modules      []string `json:"modules,omitempty"`

	Win2008Plus        bool     `json:"win2008plus,omitempty"`
	WinVCRedistDll     string   `json:"winVCRedistDll,omitempty"`
	WinVCRedistOptions []string `json:"winVCRedistOptions,omitempty"`
	WinVCRedistURL     string   `json:"winVCRedistURL,omitempty"`
	WinVCRedistVersion string   `json:"winVCRedistVersion,omitempty"`
}

const (
	cloudVersionManifest   = "https://raw.githubusercontent.com/10gen/mms/dbc4bd8b8fb0002ddaac2f45a7a6e239af5c8f60/server/conf/mongodb_version_manifest.json"
	cloudVersionManifest36 = "https://raw.githubusercontent.com/10gen/mms/dbc4bd8b8fb0002ddaac2f45a7a6e239af5c8f60/server/src/webapp-mms/static/version_manifest/3.6.json"
	cloudVersionManifest40 = "https://raw.githubusercontent.com/10gen/mms/dbc4bd8b8fb0002ddaac2f45a7a6e239af5c8f60/server/src/webapp-mms/static/version_manifest/4.0.json"
)

func fetchCloudVersionManifest(token string) (*CloudManifest, error) {
	url := cloudVersionManifest40
	if token != "" {
		url = cloudVersionManifest40 + "?token=" + token
	}
	res, err := http.Get(url)
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
