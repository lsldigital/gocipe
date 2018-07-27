package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gocipe",
	Short: "Gocipe is a web / grpc app scaffold generator",
	Long:  `A web / grpc scaffold Generator`,
}

func init() {
	generateCmd.Flags().BoolVarP(&noSkip, "overwrite", "o", false, "Overwrite.")

	rootCmd.AddCommand(generateCmd, versionCmd)
}

// Execute starts the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
