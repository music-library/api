package version

var (
	// The git commit that was compiled. This will be filled in by the compiler.
	GitBranch   string
	GitCommit   string
	GitDescribe string

	Version           = "1.1.7"
	VersionPrerelease = "dev"
	VersionMetadata   = ""
)
