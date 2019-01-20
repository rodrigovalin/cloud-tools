package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/blang/semver"
	yaml "gopkg.in/yaml.v2"
)

const ModuleEnterprise = "enterprise"

var opsManagerHostSupport OpsManagerHostSupport
var UnsupportedArchs = [...]string{"arm64", "i686", "i386"}

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

type OpsManagerHostSupport struct {
	Versions []OpsManagerVersion `yaml:"opsManagerVersion"`
}

type OpsManagerVersion struct {
	Version            string                         `yaml:"version"`
	SupportedPlatforms []OpsManagerSupportedPlatforms `yaml:"supportedPlatforms"`
}

type OpsManagerSupportedPlatforms struct {
	Platform string             `yaml:"platform"`
	Flavors  []OpsManagerFlavor `yaml:"flavors"`
}

type OpsManagerFlavor struct {
	Name              string   `yaml:"name"`
	SupportedVersions []string `yaml:"supportedVersions"`
	Arch              []string `yaml:"arch"`
}

func init() {
	supportedVersionsFile, err := Asset("assets/ops_manager_host_support.yaml")
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(supportedVersionsFile, &opsManagerHostSupport)
	if err != nil {
		panic(err)
	}
}

func (om *OpsManagerHostSupport) getVersion(omVersion string) *OpsManagerVersion {
	for _, v := range om.Versions {
		if v.Version == omVersion {
			return &v
		}
	}

	return nil
}

func (v *OpsManagerVersion) getPlatform(platform string) *OpsManagerSupportedPlatforms {
	for _, p := range v.SupportedPlatforms {
		if platform == p.Platform {
			return &p
		}
	}

	return nil
}

func (v *OpsManagerSupportedPlatforms) getFlavor(flavor string) *OpsManagerFlavor {
	for _, f := range v.Flavors {
		if f.Name == flavor {
			return &f
		}
	}

	return nil
}

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

func findLatestPatch(manifest *CloudManifest, version semver.Version) semver.Version {
	latest := semver.Version{
		Major: version.Major,
		Minor: version.Minor,
		Patch: 0,
	}
	for _, m := range manifest.Versions {
		this, _ := semver.Make(m.Name)
		if this.Major == latest.Major && this.Minor == latest.Minor && this.Patch > latest.Patch {
			latest = this
		}
	}

	return latest
}

func buildCloudManifestForVersion(newVersion string, server *ServerManifest, omVersion string) (*CloudManifest, error) {
	cloudManifest := &CloudManifest{Updated: time.Now().Unix() * 1000}
	for _, version := range server.Versions {
		if version.Version == newVersion {
			community, enterprise := buildBuildsForCloudManifestVersion(version, omVersion)
			cloudManifest.Versions = []CloudManifestVersion{{
				Builds: community,
				Name:   newVersion,
			}, {
				Builds: enterprise,
				Name:   newVersion + "-ent",
			}}
			break
		}
	}

	return cloudManifest, nil
}

// Returns true if a given version of Ops Manager (v1) should support a version of MongoDB (v2)
// Ops Manager Supports Major MongoDB versions that where released before itself.
// Example:
// Ops Manager 3.6 supports every MongoDB from 3.6 and back
// Ops Manager 4.0 supports every MongoDB from 4.0 and back, so it also supports 3.x
func OpsManagerSupportsMongoDB(v1, v2 semver.Version) bool {
	return v1.Major > v2.Major || (v1.Major == v2.Major && v1.Minor >= v2.Minor)
}

func getBuildsAsString(manifest *CloudManifest) string {
	if len(manifest.Versions) == 0 {
		fmt.Println("Could not get versions from cloud manifest")
		panic(1)
	}
	data, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		panic(err)
	}

	lines := strings.Split(string(data), "\n")

	return strings.Join(lines[3:len(lines)-2], "\n")
}

func buildBuildsForCloudManifestVersion(serverVersion ServerManifestVersion, omVersion string) ([]CloudManifestBuild, []CloudManifestBuild) {
	cloudManifestBuilds := make([]CloudManifestBuild, 0)
	cloudManifestBuildsEnt := make([]CloudManifestBuild, 0)

	for _, download := range serverVersion.Downloads {
		if shouldSkipDownload(&download, &serverVersion, omVersion) {
			continue
		}

		build := CloudManifestBuild{
			Architecture: getCloudArchFromServerArch(download.Arch),
			GitVersion:   serverVersion.Githash,
			Platform:     getPlatformFromTarget(download.Target),
			URL:          getPartialFromFullURL(download.Archive.URL),
		}

		applyLinuxAttributes(serverVersion.Version, &download, &build)
		applyWindowsAttributes(serverVersion.Version, &download, &build)

		if download.Edition == ModuleEnterprise {
			build.Modules = []string{ModuleEnterprise}
			cloudManifestBuildsEnt = append(cloudManifestBuildsEnt, build)
		} else {
			cloudManifestBuilds = append(cloudManifestBuilds, build)
		}
	}

	return cloudManifestBuilds, cloudManifestBuildsEnt
}

func getPlatformFromTarget(target string) string {
	if targetIsLinux(target) {
		return "linux"
	}

	if targetIsMacOS(target) {
		return "osx"
	}

	if targetIsWindows(target) {
		return "windows"
	}
	return ""
}

func getPartialFromFullURL(full string) string {
	splited := strings.Split(full, "/")
	return "/" + strings.Join(splited[len(splited)-2:], "/")
}

func getCloudArchFromServerArch(arch string) string {
	if arch == "x86_64" {
		return "amd64"
	}
	return arch
}

func addVersionToFile(fname, mongoVersion, omVersion string, serverManifest *ServerManifest) error {
	cloudManifest, err := fetchCloudVersionManifest(fname)
	if err != nil {
		return err
	}

	manifest, err := buildCloudManifestForVersion(mongoVersion, serverManifest, omVersion)
	if err != nil {
		return err
	}
	// TODO: check manifest contains builds, if it does not, it means the new version
	// was not found in full.json.

	newBuilds := getBuildsAsString(manifest)

	mongo, _ := semver.Make(mongoVersion)
	latest := findLatestPatch(cloudManifest, mongo)
	ent, _ := semver.NewPRVersion("ent")
	latest.Pre = append(latest.Pre, ent)

	var mySearchFunc = searchFuncForVersion(latest)
	filePos := scanFileFor(fname, mySearchFunc)

	err = insertIntoJsonFile(fname, newBuilds, filePos+2)
	if err != nil {
		return err
	}

	return nil
}

// shouldSkipDownload will skip the download definition from server's full manifest file
// for versions cloud does not support. In this particular case, no "source" or "arm64"
// versions as supported by Cloud.
func shouldSkipDownload(download *ServerManifestDownload, server *ServerManifestVersion, omVersion string) bool {
	if targetIsWindows(download.Target) && strings.Contains(download.Archive.URL, "2008plus") && !strings.Contains(download.Archive.URL, "ssl") {
		return true
	}

	if targetIsMacOS(download.Target) && download.Edition != ModuleEnterprise && strings.Contains(download.Archive.URL, "ssl") && (strings.HasPrefix(server.Version, "3.2") || strings.HasPrefix(server.Version, "3.4")) {
		// Community OSX with SSL not supported on 3.2, 3.4
		return true
	}

	arch := getCloudArchFromServerArch(download.Arch)
	if arch == "s390x" {
		if download.Edition != ModuleEnterprise {
			// skip s390x for non enterprise builds
			return true
		}

		if strings.HasPrefix(server.Version, "3.6") || strings.HasPrefix(server.Version, "3.4") {
			// skip s390x for anything below 4.0
			return true
		}
	}
	if download.Edition != ModuleEnterprise && arch == "s390x" {
		// s390x community builds are not supported.
		return true
	}

	for _, a := range UnsupportedArchs {
		if a == arch {
			return true
		}
	}

	return download.Edition == "source" || !versionIsSupported(download, server, omVersion)
}

// For a given download from the Server Manifest, calculates if it is supported for a given
// version of Ops Manager. For instance, ubuntu1804 is not supported in OM-3.6 so it should
// be skipped.
// This only applies for Linux as all of the builds for Windows & OSX are supported everywhere.
//
func versionIsSupported(version *ServerManifestDownload, server *ServerManifestVersion, omVersion string) bool {
	if version.Arch == "ppc64le" && strings.HasPrefix(server.Version, "3.2.") {
		return false
	}

	if omVersion == "" {
		// If user does not specify a given OM Version I assume they want all the builds.
		return true
	}
	if !targetIsLinux(version.Target) {
		return true
	}

	thisOm := opsManagerHostSupport.getVersion(omVersion)
	if thisOm == nil {
		return false
	}

	thisFlavor := getLinuxFlavorFromTarget(version.Target)
	thisVersion := getLinuxVersionFromTarget(version.Target)
	if thisVersion == "" {
		// targets like "linux_x86_64" are not supported in Ops Manager
		if strings.HasPrefix(version.Target, "amazon") || strings.HasPrefix(version.Target, "amzn") {
			// Amazon target don't have versions on the name, but are still supported
			return true
		}

		if version.Edition == "base" {
			// This is base Linux edition.
			return true
		}
		return false
	}

	platform := thisOm.getPlatform("Linux")
	if platform == nil {
		// if platform is not found, don't add it to cloud manifest
		return false
	}
	flavor := platform.getFlavor(thisFlavor)
	if flavor == nil {
		return false
	}

	return contains(thisVersion, flavor.SupportedVersions) &&
		contains(version.Arch, flavor.Arch)
}

func contains(needle string, haystack []string) bool {
	for _, el := range haystack {
		if needle == el {
			return true
		}
	}

	return false
}
