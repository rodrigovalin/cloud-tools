package main

import (
	"strings"

	"gopkg.in/yaml.v2"
)

// OsMinMaxVersionDefinition bla bla
type OsMinMaxVersionDefinition struct {
	Flavors []OsMinMaxVersionFlavor `yaml:"flavors"`
}

// OsMinMaxVersionFlavor bla bla
type OsMinMaxVersionFlavor struct {
	Name      string                      `yaml:"name"`
	OsVersion []OsMinMaxVersionOsVersions `yaml:"osVersions"`
}

// OsMinMaxVersionOsVersions bla bla
type OsMinMaxVersionOsVersions struct {
	Name string `yaml:"name"`
	Min  string `yaml:"min"`
	Max  string `yaml:"max"`
}

func readOsMinMaxVersions() *OsMinMaxVersionDefinition {
	yamlFile, err := Asset("assets/min_max_versions.yaml")
	// yamlFile, err := ioutil.ReadFile("min_max_versions.yaml")
	if err != nil {
		panic(err)
	}

	definition := &OsMinMaxVersionDefinition{}
	err = yaml.Unmarshal(yamlFile, definition)
	if err != nil {
		panic(err)
	}

	return definition
}

func getMinMaxOsVersionFromURL(flavorName, url string) (string, string) {
	osVersions := readOsMinMaxVersions()

	for _, flavor := range osVersions.Flavors {
		if flavor.Name == flavorName {
			for _, version := range flavor.OsVersion {
				if strings.Contains(url, version.Name) {
					return version.Min, version.Max
				}
			}
		}
	}

	return "", ""
}
