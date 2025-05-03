// Package version is just for service version.
package version

import "fmt"

// Version is the app version. To be injected.
var Version string

// GitCommit is git commit used for the build.
var GitCommit string

// BuildDate is the date when current binary was built.
var BuildDate string

// GoVersion is the used Go version.
var GoVersion string

// BuildVersionString creates version string.
func BuildVersionString(svs string) string {
	if Version == "" {
		return "<version not set>"
	}

	return fmt.Sprintf("%s %s %s %s. Built with %s", svs, Version, GitCommit, BuildDate, GoVersion)
}
