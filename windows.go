package main

import "strings"

/**
Dealing with Windows Attributes

Requirements for VCRedist* are hardcoded here. Good for now. I'm not sure if it needs more attention considering how
few changes we do to this.
*/

const (
	WinVCRedistUrl34 = "http://download.microsoft.com/download/6/D/F/6DF3FF94-F7F9-4F0B-838C-A328D1A7D0EE/vc_redist.x64.exe"

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
	if strings.HasPrefix(version, "2.") {
		return "msvcr100.dll"
	}

	if strings.HasPrefix(version, "3.0") || strings.HasPrefix(version, "3.2") {
		return "msvcr120.dll"
	}

	return "vcruntime140.dll"
}

func getWinVCRedistURL(version string) string {
	if strings.HasPrefix(version, "2.") {
		return "http://download.microsoft.com/download/1/6/5/165255E7-1014-4D0A-B094-B6A430A6BFFC/vcredist_x64.exe"
	}

	if strings.HasPrefix(version, "3.0") || strings.HasPrefix(version, "3.2") {
		return "http://download.microsoft.com/download/2/E/6/2E61CFA4-993B-4DD4-91DA-3737CD5CD6E3/vcredist_x64.exe"
	}

	// next version is used from 3.4 onwards
	return "http://download.microsoft.com/download/6/D/F/6DF3FF94-F7F9-4F0B-838C-A328D1A7D0EE/vc_redist.x64.exe"
}

func getWinVCRedistVersion(version string) string {
	if strings.HasPrefix(version, "2.") {
		return "10.0.40219.325"
	}

	if strings.HasPrefix(version, "3.0") || strings.HasPrefix(version, "3.2") {
		return "12.0.21005.1"
	}

	// next version is used from 3.4 onwards
	return "14.0.24212.0"
}

func getWinVCRedistOptions(version string) []string {
	if strings.HasPrefix(version, "2.") {
		return []string{"/q", "/norestart"}
	}

	return []string{"/quiet", "/norestart"}
}

func getWin2008Plus() bool {
	return true
}

// Not sure what this is
func getMsi() string {
	return "COMPLETE_ME"
}

func targetIsWindows(target string) bool {
	return strings.Contains(target, "windows")
}
