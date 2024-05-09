package cmd

import (
	"bypctl/pkg/color"
	"bypctl/pkg/compose"
	"bypctl/pkg/global"
	"bypctl/pkg/i18n"
	"github.com/spf13/cobra"
)

// startCmd represents the version command
var startCmd = &cobra.Command{
	Use:   "start [SERVICE...]",
	Short: i18n.Translate("start_help", "Start services"),
	Long:  i18n.Translate("start_help", "Start services"),
	Run: func(cmd *cobra.Command, args []string) {
		if err := compose.Start(global.Conf.System.ComposeFiles, args); err != nil {
			color.PrintRed(i18n.Tf("start_err", global.Conf.System.ComposeProfiles, err.Error()))
		}
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
