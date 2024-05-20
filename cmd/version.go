package cmd

import (
	"bypctl/pkg/i18n"
	"fmt"
	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: i18n.Translate("version_help", "Show the bypanel version information"),
	Long:  i18n.Translate("version_help", "Show the bypanel version information"),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("v0.0.1")
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
