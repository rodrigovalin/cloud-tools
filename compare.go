package main

import (
	"fmt"
)

func findVersion(m *CloudManifest, version string) *CloudManifestVersion {
	for _, el := range m.Versions {
		if el.Name == version {
			return &el
		}
	}

	return nil
}

func findBuildByUrl(builds []CloudManifestBuild, url string) *CloudManifestBuild {
	for _, b := range builds {
		if b.URL == url {
			return &b
		}
	}

	return nil
}

func compareManifestsForVersion(c1, c2, version string) int {
	c, err := fetchCloudVersionManifest(c1)
	if err != nil {
		fmt.Printf("Could not read %s file\n", c1)
		return 1
	}
	v1 := findVersion(c, version)
	if v1 == nil {
		fmt.Printf("Version %s does not exist in %s file\n", version, c1)
		return 1
	}

	c, err = fetchCloudVersionManifest(c2)
	if err != nil {
		fmt.Printf("Could not read %s file\n", c2)
		return 1
	}
	v2 := findVersion(c, version)
	if v2 == nil {
		fmt.Printf("Version %s does not exist in %s file\n", version, c2)
		return 1
	}

	buildsLess := make([]*CloudManifestBuild, 0)
	buildsMore := make([]*CloudManifestBuild, 0)

	for _, b := range v1.Builds {
		found := findBuildByUrl(v2.Builds, b.URL)
		if found == nil {
			buildsLess = append(buildsLess, &b)
		}
	}
	for _, b := range v2.Builds {
		found := findBuildByUrl(v1.Builds, b.URL)
		if found == nil {
			buildsMore = append(buildsMore, &b)
		}
	}

	if len(buildsLess) > 0 {
		fmt.Printf("Compare failed, %d builds are missing from %s\n", len(buildsLess), c2)
		for _, b := range buildsLess {
			fmt.Println(b.URL)
		}
	}
	if len(buildsMore) > 0 {
		fmt.Printf("Compare failed, %d builds are missing from %s\n", len(buildsMore), c1)
		for _, b := range buildsMore {
			fmt.Println(b.URL)
		}

	}

	if len(buildsLess) == 0 && len(buildsLess) == 0 {
		fmt.Println("All builds were found!")
	}

	return 0
}
