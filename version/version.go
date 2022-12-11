package version

import (
	"bytes"
	"fmt"
)

// VersionInfo
type VersionInfo struct {
	Revision          string
	Version           string
	VersionPrerelease string
	VersionMetadata   string
}

func GetVersion() *VersionInfo {
	ver := Version
	rel := VersionPrerelease
	md := VersionMetadata
	if GitDescribe != "" {
		ver = GitDescribe
	}
	if GitDescribe == "" && rel == "" && VersionPrerelease != "" {
		rel = "dev"
	}

	return &VersionInfo{
		Revision:          GitCommit,
		Version:           ver,
		VersionPrerelease: rel,
		VersionMetadata:   md,
	}
}

func (c *VersionInfo) VersionNumber() string {
	if Version == "unknown" && VersionPrerelease == "unknown" {
		return "Version unknown"
	}

	version := fmt.Sprintf("%s", c.Version)

	if c.VersionPrerelease != "" {
		version = fmt.Sprintf("%s-%s", version, c.VersionPrerelease)
	}

	if c.VersionMetadata != "" {
		version = fmt.Sprintf("%s+%s", version, c.VersionMetadata)
	}

	return version
}

func (c *VersionInfo) FullVersionNumber(rev bool) string {
	var versionString bytes.Buffer

	if Version == "unknown" && VersionPrerelease == "unknown" {
		return "music-api [Version unknown]"
	}

	fmt.Fprintf(&versionString, "music-api [Version %s", c.Version)
	if c.VersionPrerelease != "" {
		fmt.Fprintf(&versionString, "-%s", c.VersionPrerelease)
	}

	if c.VersionMetadata != "" {
		fmt.Fprintf(&versionString, "+%s", c.VersionMetadata)
	}

	if rev && c.Revision != "" {
		fmt.Fprintf(&versionString, " %s", c.Revision)
	}

	fmt.Fprintf(&versionString, "]")

	return versionString.String()
}

func PrintTitle() {
	// Get version info
	versionStruct := GetVersion()

	// Check if dev release
	// Hide git-sha for main releases
	isDev := false
	if versionStruct.VersionPrerelease != "" {
		isDev = true
	}

	// Get full version string
	versionString := versionStruct.FullVersionNumber(isDev)

	fmt.Println(versionString)
	fmt.Println("(c) MerrittCorp. All rights reserved.")
	fmt.Println()
}
