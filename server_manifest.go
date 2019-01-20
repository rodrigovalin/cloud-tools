package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
)

const (
	serverVersionManifest      = "http://downloads.mongodb.org.s3.amazonaws.com/full.json"
	serverVersionManifestLocal = "full.json"
)

type ServerManifest struct {
	Versions []ServerManifestVersion `json:"versions"`
}

type ServerManifestVersion struct {
	Changes            string                   `json:"changes"`
	Current            bool                     `json:"current"`
	Date               string                   `json:"date"`
	DevelopmentRelease bool                     `json:"development_release"`
	Downloads          []ServerManifestDownload `json:"downloads"`
	Githash            string                   `json:"githash"`
	ProductionRelease  bool                     `json:"production_release"`
	ReleaseCandidate   bool                     `json:"release_candidate"`
	Version            string                   `json:"version"`
	Notes              string                   `json:"notes,omitempty"`
}

type ServerManifestDownload struct {
	Arch     string                `json:"arch,omitempty"`
	Archive  ServerManifestArchive `json:"archive"`
	Edition  string                `json:"edition"`
	Packages []string              `json:"packages,omitempty"`
	Target   string                `json:"target,omitempty"`
	Msi      string                `json:"msi,omitempty"`
}

type ServerManifestArchive struct {
	DebugSymbols string `json:"debug_symbols,omitempty"`
	Sha1         string `json:"sha1"`
	Sha256       string `json:"sha256"`
	URL          string `json:"url"`
}

func (s *ServerManifest) HasBuild(version string) bool {
	for _, v := range s.Versions {
		if v.Version == version {
			return true
		}
	}

	return false
}

func fetchServerVersionManifest() (*ServerManifest, error) {
	var body []byte
	var err error
	_, ok := os.LookupEnv("IM_WORKING")
	// set IM_WORKING to something to read the full.json from disk instead of S3
	// makes testing faster!
	if ok {
		body, err = ioutil.ReadFile(serverVersionManifestLocal)
		if err != nil {
			return nil, err
		}
	} else {
		res, err := http.Get(serverVersionManifest)
		if err != nil {
			return nil, err
		}

		body, err = ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}
	}

	s := &ServerManifest{}
	err = json.Unmarshal(body, &s)
	if err != nil {
		return nil, err
	}

	return s, nil
}
