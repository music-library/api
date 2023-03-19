package version

var (
	// The git commit that was compiled. This will be filled in by the compiler.
	GitBranch   string
	GitCommit   string
	GitDescribe string

	Version           = "0.7.35"
	VersionPrerelease = "dev"
	VersionMetadata   = ""
)
