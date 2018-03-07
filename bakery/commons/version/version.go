package version

import (
	"fmt"
)

var (
	AppName   string
	Version   string
	BuildTime string
	GitCommit string
)

func String() string {
	return fmt.Sprintf("%s\nUTC Build Time: %s\nGit Commit Hash: %s", Version, BuildTime, GitCommit)
}
