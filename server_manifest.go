package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type ServerManifest struct {
	Versions []ServerManifestVersion `json:"versions"`
}

type ServerManifestVersion struct {
	Changes            string                 `json:"changes"`
	Current            bool                   `json:"current"`
	Date               string                 `json:"date"`
	DevelopmentRelease bool                   `json:"development_release"`
	Downloads          ServerManifestDownload `json:"downloads"`
	Githash            string                 `json:"githash"`
	ProductionRelease  bool                   `json:"production_release"`
	ReleaseCandidate   bool                   `json:"release_candidate"`
	Version            string                 `json:"version"`
	Notes              string                 `json:"notes,omitempty"`
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

func fetchServerVersionManifest(token string) (*ServerManifest, error) {
	url := serverVersionManifest
	if token != "" {
		url = serverVersionManifest + "?token=" + token
	}
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	s := &ServerManifest{}
	err = json.Unmarshal(body, &s)
	if err != nil {
		return nil, err
	}
	return s, nil
}
