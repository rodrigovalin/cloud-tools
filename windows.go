package main

import (
	"strings"

	yaml "gopkg.in/yaml.v2"
)

type WinVCDefinition struct {
	Versions []WinVCVersion `yaml:"versions"`
}

type WinVCVersion struct {
	Prefix  []string `yaml:"prefix"`
	URL     string   `yaml:"url"`
	DLL     string   `yaml:"dll"`
	Version string   `yaml:"version"`
	Options []string `yaml:"options"`
}

var WinVC = &WinVCDefinition{}

func init() {
	f, err := Asset("assets/winvc_versions.yaml")
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(f, &WinVC)
	if err != nil {
		panic(err)
	}
}

func applyWindowsAttributes(serverVersion string, download *ServerManifestDownload, build *CloudManifestBuild) {
	if !targetIsWindows(download.Target) {
		return
	}

	build.Win2008Plus = getWin2008Plus(download)
	if !build.Win2008Plus {
		// following attributes are only set on 2008plus versions.
		return
	}
	url, dll, version, options := getWinVCAttributes(serverVersion)
	build.WinVCRedistDll = dll
	build.WinVCRedistOptions = options
	build.WinVCRedistURL = url
	build.WinVCRedistVersion = version

}

func isVersionCatchAll(prefix []string) bool {
	if len(prefix) == 0 {
		return true
	}

	if len(prefix) == 1 && prefix[0] == "*" {
		return true
	}

	return false
}

func getWinVCAttributes(version string) (string, string, string, []string) {
	for _, v := range WinVC.Versions {
		// first check that a given version is matched with one of the prefixes
		// skipping the catch-all version
		if isVersionCatchAll(v.Prefix) {
			continue
		}
		for _, prefix := range v.Prefix {
			if strings.HasPrefix(version, prefix) {
				return v.URL, v.DLL, v.Version, v.Options
			}
		}
	}

	// if a version didn't match, return the catch-all, if there's any
	for _, v := range WinVC.Versions {
		if isVersionCatchAll(v.Prefix) {
			return v.URL, v.DLL, v.Version, v.Options
		}
	}

	return "", "", "", []string{}
}

func getWin2008Plus(download *ServerManifestDownload) bool {
	return strings.Contains(download.Archive.URL, "2008plus") || strings.Contains(download.Archive.URL, "windows-64")
}

func targetIsWindows(target string) bool {
	return strings.Contains(target, "windows")
}
