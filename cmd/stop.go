package cmd

import (
	"bypctl/pkg/color"
	"bypctl/pkg/compose"
	"bypctl/pkg/global"
	"bypctl/pkg/i18n"
	"github.com/spf13/cobra"
	"strings"
)

// stopCmd represents the version command
var stopCmd = &cobra.Command{
	Use:   "stop [SERVICE...]",
	Short: i18n.Translate(`stop_help`),
	Long:  i18n.Translate(`stop_help`),
	Run: func(cmd *cobra.Command, args []string) {
		if err := compose.Stop(global.Conf.System.ComposeFiles, args); err != nil {
			color.PrintRed(i18n.Tf("stop_err", err.Error()))
		} else {
			var apps []string
			if len(args) > 0 {
				apps = args
			} else {
				apps = strings.Split(global.Conf.System.ComposeProfiles, ",")
			}
			for _, v := range apps {
				color.PrintGreen(i18n.Tf("stop_succ", v))
			}
		}
		// fmt.Println("stdout--->", stdout)
	},
}

func init() {
	rootCmd.AddCommand(stopCmd)
}
