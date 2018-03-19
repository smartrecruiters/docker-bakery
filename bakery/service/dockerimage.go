package service

import "github.com/Masterminds/semver"

// GetLatestVersion returns latest version of the docker image or "0.0.0" if there was no version defined
func (di *DockerImage) GetLatestVersion() *semver.Version {
	if di.latestVersion != nil {
		return di.latestVersion
	}

	di.latestVersion, _ = semver.NewVersion("0.0.0")
	return di.latestVersion
}

// GetLatestVersionString returns the latest version (as a string) of the docker image or "0.0.0" if there was no version defined
func (di *DockerImage) GetLatestVersionString() string {
	return di.GetLatestVersion().String()
}

// CalculateNextVersion returns next version of the docker image based on the provided scope (major/minor/patch).
// If image had no previous version the the 0.0.0 is used as a base line and appropriately updated with regards
// to the provided scope.
func (di *DockerImage) CalculateNextVersion(scope string) {
	version := di.GetLatestVersion()
	switch scope {
	case "major":
		di.nextVersion = version.IncMajor()
	case "minor":
		di.nextVersion = version.IncMinor()
	case "patch":
		di.nextVersion = version.IncPatch()
	default:
		di.nextVersion = version.IncPatch()
	}
}

// GetNextVersion returns already calculated next version of the docker image.
func (di *DockerImage) GetNextVersion() semver.Version {
	return di.nextVersion
}

// GetNextVersionString returns already calculated next version (as a string) of the docker image.
func (di *DockerImage) GetNextVersionString() string {
	return di.nextVersion.String()
}
