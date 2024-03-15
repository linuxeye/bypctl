package cmd

import (
	"bypctl/pkg/color"
	"bypctl/pkg/compose"
	"bypctl/pkg/global"
	"bypctl/pkg/i18n"
	"github.com/spf13/cobra"
	"strings"
)

// downCmd represents the version command
var downCmd = &cobra.Command{
	Use:   "down [SERVICE...]",
	Short: i18n.Translate(`down_help`),
	Long:  i18n.Translate(`down_help`),
	Run: func(cmd *cobra.Command, args []string) {
		if err := compose.Down(global.Conf.System.ComposeFiles, args); err != nil {
			color.PrintRed(i18n.Tf("down_err", err.Error()))
		} else {
			var apps []string
			if len(args) > 0 {
				apps = args
			} else {
				apps = strings.Split(global.Conf.System.ComposeProfiles, ",")
			}
			for _, v := range apps {
				color.PrintGreen(i18n.Tf("down_succ", v))
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(downCmd)
}
