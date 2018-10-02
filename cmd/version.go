package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Versioning info
var (
	appVersion = "n/a"
	appCommit  = "n/a"
	appBuilt   = "n/a"
)

var versionCmd = &cobra.Command{
	Use: "version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Version : %v \nCommit : %v\nBuilt: %v\n", appVersion, appCommit, appBuilt)
	},
}

// SetVersionInfo initializes version, commit and built information
func SetVersionInfo(version, commit, built string) {
	appVersion = version
	appCommit = commit
	appBuilt = built
}
