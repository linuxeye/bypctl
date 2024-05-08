package cmd

import (
	"bypctl/pkg/color"
	"bypctl/pkg/compose"
	"bypctl/pkg/global"
	"bypctl/pkg/i18n"
	"github.com/spf13/cobra"
	"strings"
)

// pullCmd represents the version command
var pullCmd = &cobra.Command{
	Use:   "pull [SERVICE...]",
	Short: i18n.Translate("pull_help", "Pull service images"),
	Long:  i18n.Translate("pull_help", "Pull service images"),
	Run: func(cmd *cobra.Command, args []string) {
		if err := compose.Pull(global.Conf.System.ComposeFiles, args); err != nil {
			color.PrintRed(i18n.Tf("pull_docker_err", "Pull {{ .Apps }} image, error: {{ .Err }}", map[string]any{"Apps": global.Conf.System.ComposeProfiles, "Err": err.Error()}))
		} else {
			var apps []string
			if len(args) > 0 {
				apps = args
			} else {
				apps = strings.Split(global.Conf.System.ComposeProfiles, ",")
			}
			for _, v := range apps {
				color.PrintGreen(i18n.Tf("pull_docker_succ", "✔ {{ .App }} Pulled", map[string]any{"App": v}))
			}
		}
		// fmt.Println(stdout)
	},
}

func init() {
	rootCmd.AddCommand(pullCmd)
}
