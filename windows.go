package main

import "strings"

/**
Dealing with Windows Attributes
*/

const (
	WinVCRedistDll34 = "vcruntime140.dll"
	WinVCRedistUrl34 = "http://download.microsoft.com/download/6/D/F/6DF3FF94-F7F9-4F0B-838C-A328D1A7D0EE/vc_redist.x64.exe"

	WinVCRedistDll = "msvcr120.dll"
	WinVCRedistUrl = "http://download.microsoft.com/download/2/E/6/2E61CFA4-993B-4DD4-91DA-3737CD5CD6E3/vcredist_x64.exe"

	WinVCRedistVersion   = "10.0.40219.325"
	WinVCRedistVersion3  = "12.0.21005.1"
	WinVCRedistVersion34 = "14.0.24212.0"
)

func applyWindowsAttributes(serverVersion string, download *ServerManifestDownload, build *CloudManifestBuild) {
	if !targetIsWindows(download.Target) {
		return
	}

	build.Win2008Plus = getWin2008Plus()
	build.WinVCRedistDll = getWinVCRedistDll(serverVersion)
	build.WinVCRedistOptions = getWinVCRedistOptions(serverVersion)
	build.WinVCRedistURL = getWinVCRedistURL(serverVersion)
	build.WinVCRedistVersion = getWinVCRedistVersion(serverVersion)

}

func getWinVCRedistDll(version string) string {
	dll, _ := getWinRCRedistDll(version)
	return dll
}

func getWinVCRedistURL(version string) string {
	_, url := getWinRCRedistDll(version)
	return url
}

func getWinVCRedistOptions(version string) []string {
	return []string{"/quiet", "/norestart"}
}

func getWinVCRedistVersion(version string) string {
	return WinVCRedistVersion
}

func getWin2008Plus() bool {
	return true
}

func getMsi() string {
	return "COMPLETE_ME"
}

func getWinRCRedistDll(version string) (string, string) {
	if strings.HasPrefix(version, "3.4") ||
		strings.HasPrefix(version, "3.6") ||
		strings.HasPrefix(version, "4.") {
		return WinVCRedistDll34, WinVCRedistUrl34
	}

	return WinVCRedistDll, WinVCRedistUrl
}

func targetIsWindows(target string) bool {
	return strings.Contains(target, "windows")
}
