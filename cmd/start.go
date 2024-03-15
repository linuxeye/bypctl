package cmd

import (
	"bypctl/pkg/color"
	"bypctl/pkg/compose"
	"bypctl/pkg/global"
	"bypctl/pkg/i18n"
	"github.com/spf13/cobra"
	"strings"
)

// startCmd represents the version command
var startCmd = &cobra.Command{
	Use:   "start [SERVICE...]",
	Short: i18n.Translate(`start_help`),
	Long:  i18n.Translate(`start_help`),
	Run: func(cmd *cobra.Command, args []string) {
		if err := compose.Start(global.Conf.System.ComposeFiles, args); err != nil {
			color.PrintRed(i18n.Tf("start_err", global.Conf.System.ComposeProfiles, err.Error()))
		} else {
			var apps []string
			if len(args) > 0 {
				apps = args
			} else {
				apps = strings.Split(global.Conf.System.ComposeProfiles, ",")
			}
			for _, v := range apps {
				color.PrintGreen(i18n.Tf("start_succ", v))
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
