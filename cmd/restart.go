package cmd

import (
	"bypctl/pkg/color"
	"bypctl/pkg/compose"
	"bypctl/pkg/global"
	"bypctl/pkg/i18n"
	"github.com/spf13/cobra"
	"strings"
)

// restartCmd represents the version command
var restartCmd = &cobra.Command{
	Use:   "restart [SERVICE...]",
	Short: i18n.Translate(`restart_help`),
	Long:  i18n.Translate(`restart_help`),
	Run: func(cmd *cobra.Command, args []string) {
		if err := compose.Restart(global.Conf.System.ComposeFiles, args); err != nil {
			color.PrintRed(i18n.Tf("restart_err", err.Error()))
		} else {
			var apps []string
			if len(args) > 0 {
				apps = args
			} else {
				apps = strings.Split(global.Conf.System.ComposeProfiles, ",")
			}
			for _, v := range apps {
				color.PrintGreen(i18n.Tf("restart_succ", v))
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(restartCmd)
}
