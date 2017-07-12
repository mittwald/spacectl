package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/mittwald/spacectl/buildinfo"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display current spacectl version",
	Long: `Display current spacectl version`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("spacectl %s (commit ID %s)\n", buildinfo.Version, buildinfo.Hash)
		fmt.Printf("built at %s\n", buildinfo.BuildDate)
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
