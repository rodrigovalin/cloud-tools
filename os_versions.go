package main

import (
	"strings"

	yaml "gopkg.in/yaml.v2"
)

type OsMinMaxVersionDefinition struct {
	Flavors []OsMinMaxVersionFlavor `yaml:"flavors"`
}

type OsMinMaxVersionFlavor struct {
	Name      string                      `yaml:"name"`
	OsVersion []OsMinMaxVersionOsVersions `yaml:"osVersions"`
}

type OsMinMaxVersionOsVersions struct {
	Name string `yaml:"name"`
	Min  string `yaml:"min"`
	Max  string `yaml:"max"`
}

var OsMinMaxVersions = &OsMinMaxVersionDefinition{}

func init() {
	yamlFile, err := Asset("assets/min_max_versions.yaml")
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(yamlFile, &OsMinMaxVersions)
	if err != nil {
		panic(err)
	}
}

func getMinMaxOsVersionFromURL(flavorName, url string) (string, string) {
	for _, flavor := range OsMinMaxVersions.Flavors {
		if strings.HasPrefix(flavorName, flavor.Name) {
			for _, version := range flavor.OsVersion {
				if strings.Contains(url, version.Name) {
					return version.Min, version.Max
				}
			}
		}
	}

	return "", ""
}
