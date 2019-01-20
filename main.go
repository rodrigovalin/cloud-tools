package main

/*
 Cloud Tools. A collection of tools for sunny days.

 This tools was conceived to help on routine tasks that the cloud team needs to perform.
 We are trying to accomplish the following goals.

 - find an easy way to distribute tools across the organization
 - make it easy to do routine tasks with them.

*/

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"time"

	"github.com/blang/semver"
	docopt "github.com/docopt/docopt-go"
)

var usage = `Cloud Tools -- Tools for Sunny Days.

Usage:
  cloud-tools version-manifest add-build --new=<mongod-version> --merge-with=<file-merge> [--into=<file-into>] [--om-version=<om-version>]
  cloud-tools version-manifest add-build --new=<mongod-version> --as-pr
  cloud-tools version-manifest list-versions --manifest=<manifest-file>
  cloud-tools version-manifest compare --src=<src-manifest> --dst=<dst-manifest> --compare=<mongod-version>
  cloud-tools -h | --help
  cloud-tools --version
`

const (
	ManifestDir = "server/src/webapp-mms/static/version_manifest"
	ConfDir     = "server/conf"
)

func main() {
	// cwd := getCWD()
	root, err := findRoot()
	if err != nil {
		fmt.Println("Could not find mms repo base directory.")
		fmt.Println("Please run `cloud-tools` from inside the git repo.")
		os.Exit(1)
	}
	args, _ := docopt.ParseDoc(usage)

	addBuild, err := args.Bool("add-build")
	if err == nil && addBuild {
		asPr, err := args.Bool("--as-pr")
		if err == nil && asPr {
			// change all the required files at once, and make them ready for PR
			mongoDVersion, _ := args.String("--new")
			os.Exit(addBuildOperationAsPr(mongoDVersion, root))
		}
		mongoDVersion, _ := args.String("--new")
		fileMerge, _ := args.String("--merge-with")
		fileInto, _ := args.String("--into")
		omVersion, _ := args.String("--om-version")
		os.Exit(addBuildOperation(mongoDVersion, fileMerge, fileInto, omVersion))
	}

	compareVersion, err := args.Bool("compare")
	if err == nil && compareVersion {
		m0, _ := args.String("--src")
		m1, _ := args.String("--dst")
		version, _ := args.String("--compare")

		os.Exit(compareManifestsForVersion(m0, m1, version))
	}

	listVersions, err := args.Bool("list-versions")
	if err == nil && listVersions {
		manifest, _ := args.String("<manifest-file>")
		os.Exit(listVersionsOperation(manifest))
	}
}

func listVersionsOperation(manifest string) int {
	fmt.Printf("Getting versions for manifest: %s\n", manifest)
	versions, err := listVersionsForFile(manifest)
	if err != nil {
		fmt.Println(err)
		return 1
	}

	for _, version := range versions {
		fmt.Println(version)
	}

	return 0
}

func isGreaterThan(a, b string) bool {
	av, _ := semver.Make(a)
	bv, _ := semver.Make(b)

	return av.GT(bv)
}

func addBuildOperationAsPr(mongoVersion, repo string) int {
	serverManifest, err := fetchServerVersionManifest()
	if err != nil {
		fmt.Println(err)
		return 1

	}
	if !serverManifest.HasBuild(mongoVersion) {
		fmt.Printf("MongoDB Version %s does not exists in Server Manifest\n", mongoVersion)
		os.Exit(1)
	}

	var VersionManifestReleases = [...]string{"3.4", "3.6", "4.0"}

	mongo, _ := semver.Make(mongoVersion)
	fmt.Printf("Adding New MongoDB Version %s\n", mongo)
	timestamp := strconv.FormatInt(time.Now().UTC().Truncate(24*time.Hour).Unix(), 10)
	for idx, omVersion := range VersionManifestReleases {
		current, _ := semver.Make(fmt.Sprintf("%s.0", omVersion))
		fname := fmt.Sprintf("%s/%s/%s.json", repo, ManifestDir, omVersion)
		thisUpdated := fmt.Sprintf("%s0%d%d", timestamp, current.Major, current.Minor)

		err = editUpdatedField(fname, thisUpdated)
		if err != nil {
			fmt.Println(err)
			return 1
		}

		if !OpsManagerSupportsMongoDB(current, mongo) {
			continue
		}

		fmt.Printf("Adding version to %s\n", fname)
		err = addVersionToFile(fname, mongoVersion, omVersion, serverManifest)
		if err != nil {
			fmt.Println(err)
			return 1
		}

		if idx == len(VersionManifestReleases)-1 {
			// also modify `mongodb_version_manifest.json`
			err = editUpdatedField(fname, thisUpdated)
			if err != nil {
				fmt.Println(err)
				return 1
			}
			fname = fmt.Sprintf("%s/%s/mongodb_version_manifest.json", repo, ConfDir)
			err = addVersionToFile(fname, mongoVersion, omVersion, serverManifest)
			if err != nil {
				fmt.Println(err)
				return 1
			}
		}
	}

	return 0
}

func addBuildOperation(mongoDVersion, fileMerge, fileInto, omVersion string) int {
	serverManifest, err := fetchServerVersionManifest()
	if err != nil {
		fmt.Println(err)
		return 1
	}

	cloudManifest, err := fetchCloudVersionManifest(fileMerge)
	if err != nil {
		fmt.Println(err)
		return 1
	}

	updatedCloudManifest, err := buildCloudManifestForVersion(mongoDVersion, serverManifest, omVersion)
	if err != nil {
		fmt.Println(err)
		return 1
	}

	cloudManifest.Updated = updatedCloudManifest.Updated

	added := false
	for idx, v := range cloudManifest.Versions {
		if isGreaterThan(v.Name, mongoDVersion) {
			cloudManifest.Versions = append(cloudManifest.Versions, CloudManifestVersion{})
			cloudManifest.Versions = append(cloudManifest.Versions, CloudManifestVersion{})
			copy(cloudManifest.Versions[idx+2:], cloudManifest.Versions[idx:])
			cloudManifest.Versions[idx] = updatedCloudManifest.Versions[0]
			cloudManifest.Versions[idx+1] = updatedCloudManifest.Versions[1]
			added = true
			break
		}
	}
	if !added {
		// If unable to find a place in the array, put it at the end
		cloudManifest.Versions = append(cloudManifest.Versions, updatedCloudManifest.Versions...)
	}

	manifest, err := json.MarshalIndent(cloudManifest, "", "  ")
	if err != nil {
		fmt.Println(err)
		return 1
	}

	if fileInto != "" {
		if isS3File(fileInto) {
			err = writeS3File(fileInto, manifest)
			if err != nil {
				return 1
			}
			return 0
		}
		err = ioutil.WriteFile(fileInto, manifest, 0644)
		if err != nil {
			fmt.Println(err)
			return 1
		}

		return 0
	}
	fmt.Println(string(manifest))
	return 0
}
