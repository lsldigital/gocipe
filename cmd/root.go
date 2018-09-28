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
	generateCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Verbose.")
	generateCmd.Flags().BoolVarP(&generateBootstrap, "Bootstrap", "", true, "")
	generateCmd.Flags().BoolVarP(&generateSchema, "Schema", "", true, "")
	generateCmd.Flags().BoolVarP(&generateCrud, "Crud", "", true, "")
	generateCmd.Flags().BoolVarP(&generateAdmin, "Admin", "", true, "")
	generateCmd.Flags().BoolVarP(&generateAuth, "Auth", "", true, "")
	generateCmd.Flags().BoolVarP(&generateUtils, "Utils", "", true, "")
	generateCmd.Flags().BoolVarP(&generateVuetify, "Vuetify", "", true, "")

	rootCmd.AddCommand(generateCmd, versionCmd)
}

// Execute starts the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
