package config

import "fmt"

var (
	Version   = "1.3.0"
	GitBranch = "unknown"
	GitCommit = "ffffff"
)

func FormattedVersion() string {
	return fmt.Sprintf("%s-%s:%s", Version, GitBranch, GitCommit)
}
