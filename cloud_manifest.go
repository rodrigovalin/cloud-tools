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
	Architecture string `json:"architecture"`
	GitVersion   string `json:"gitVersion"`
	Platform     string `json:"platform"`
	URL          string `json:"url"`
	Win2008Plus  bool   `json:"win2008plus,omitempty"`
}

func fetchCloudVersionManifest() (*CloudManifest, error) {
	res, err := http.Get(cloudVersionManifest40)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	s := &CloudManifest{}
	json.Unmarshal(body, &s)

	return s, nil
}
