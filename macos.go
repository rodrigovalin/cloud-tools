package main

import "strings"

func targetIsMacOS(target string) bool {
	return strings.Contains(target, "macos") || strings.Contains(target, "osx")
}
