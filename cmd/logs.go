package cmd

import (
	"bypctl/pkg/compose"
	"bypctl/pkg/global"
	"bypctl/pkg/i18n"
	"github.com/spf13/cobra"
)

var followMode bool

// logsCmd represents the version command
var logsCmd = &cobra.Command{
	Use:   "logs [OPTIONS] [SERVICE...]",
	Short: i18n.Translate("logs_help", "View output from containers"),
	Long:  i18n.Translate("logs_help", "View output from containers"),
	Run: func(cmd *cobra.Command, args []string) {
		compose.Logs(global.Conf.System.ComposeFiles, args, followMode)
	},
}

func init() {
	logsCmd.Flags().BoolVarP(&followMode, "follow", "f", false, i18n.Translate("logs_follow_help", "Follow log output."))
	rootCmd.AddCommand(logsCmd)
}
