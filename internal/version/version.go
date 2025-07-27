package version

// These variables will be set during compilation via -ldflags
var (
	Version   = "unknown!"
	CommitSHA = "unknown!"
	BuildTime = "unknown!"
)

// GetVersion returns the full version string
func GetVersion() string {
	return Version
}

// GetCommitSHA returns the commit hash
func GetCommitSHA() string {
	return CommitSHA
}

// GetBuildTime returns the build time
func GetBuildTime() string {
	return BuildTime
}

// GetFullVersion returns version with commit hash
func GetFullVersion() string {
	if Version == "unknown" {
		return "dev-" + CommitSHA[:8]
	}
	return Version + "-" + CommitSHA[:8]
}
