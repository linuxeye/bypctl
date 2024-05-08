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
	Short: i18n.Translate("stop_help", "Stop services"),
	Long:  i18n.Translate("stop_help", "Stop services"),
	Run: func(cmd *cobra.Command, args []string) {
		if err := compose.Stop(global.Conf.System.ComposeFiles, args); err != nil {
			color.PrintRed(i18n.Tf("stop_err", "Stop Container error: {{ .Err }}", map[string]any{"Err": err.Error()}))
		} else {
			var apps []string
			if len(args) > 0 {
				apps = args
			} else {
				apps = strings.Split(global.Conf.System.ComposeProfiles, ",")
			}
			for _, v := range apps {
				color.PrintGreen(i18n.Tf("stop_succ", "✔ Container {{ .App }} Stopped", map[string]any{"App": v}))
			}
		}
		// fmt.Println("stdout--->", stdout)
	},
}

func init() {
	rootCmd.AddCommand(stopCmd)
}
