package main

import "strings"

// FlavorLinux contains supported linux flavors. Also helps detect the "plaform"
var FlavorLinux = [...]string{"suse", "rhel", "ubuntu", "debian", "amazon", "linux"}

func applyLinuxAttributes(serverVersion string, download *ServerManifestDownload, build *CloudManifestBuild) {
	build.Flavor = getLinuxFlavorFromTarget(download.Target)
	minOsVersion, maxOsVersion := getMinMaxOsVersionFromURL(build.Flavor, build.URL)
	build.MinOsVersion = minOsVersion
	build.MaxOsVersion = maxOsVersion

}

func getLinuxFlavorFromTarget(target string) string {
	for _, flavor := range FlavorLinux {
		if strings.Contains(target, flavor) {
			return flavor
		}
	}
	return ""
}

func targetIsLinux(target string) bool {
	for _, el := range FlavorLinux {
		if strings.Contains(target, el) {
			return true
		}
	}
	return false
}
