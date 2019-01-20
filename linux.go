package main

import "strings"

// FlavorLinux contains supported linux flavors. Also helps detect the "plaform"
var FlavorLinux = [...]string{"suse", "rhel", "ubuntu", "debian", "amazon2", "amazon", "amzn64", "linux"}

func applyLinuxAttributes(serverVersion string, download *ServerManifestDownload, build *CloudManifestBuild) {
	flavor := getLinuxFlavorFromTarget(download.Target)
	if download.Edition != "base" {
		build.Flavor = flavor
	}

	minOsVersion, maxOsVersion := getMinMaxOsVersionFromURL(flavor, build.URL)
	var emptyOsVersion = ""
	if minOsVersion != "" {
		if minOsVersion == "<empty>" {
			build.MinOsVersion = &emptyOsVersion
		} else {
			build.MinOsVersion = &minOsVersion
		}
	}

	if maxOsVersion != "" {
		if maxOsVersion == "<empty>" {
			build.MaxOsVersion = &emptyOsVersion
		} else {
			build.MaxOsVersion = &maxOsVersion
		}
	}
}

func getLinuxFlavorFromTarget(target string) string {
	for _, flavor := range FlavorLinux {
		if strings.Contains(target, flavor) {
			if flavor == "amzn64" {
				return "amazon"
			}
			return flavor
		}
	}
	return ""
}

func getLinuxVersionFromTarget(target string) string {
	if strings.HasPrefix(target, "linux") {
		// Linux appears as "linux_x86_64" which is not a version
		return ""
	}
	flavorLen := len(getLinuxFlavorFromTarget(target))
	return target[flavorLen:]
}

func targetIsLinux(target string) bool {
	for _, el := range FlavorLinux {
		if strings.Contains(target, el) {
			return true
		}
	}
	return false
}
