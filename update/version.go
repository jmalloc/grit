package update

import "github.com/Masterminds/semver"

// IsPreRelease returns true if v is considered a pre-release version.
func IsPreRelease(v *semver.Version) bool {
	return v.Major() == 0 || v.Prerelease() != ""
}
