package cmd

import (
	"bypctl/pkg/color"
	"bypctl/pkg/compose"
	"bypctl/pkg/global"
	"bypctl/pkg/i18n"
	"github.com/spf13/cobra"
)

// stopCmd represents the version command
var stopCmd = &cobra.Command{
	Use:   "stop [SERVICE...]",
	Short: i18n.Translate("stop_help", "Stop services"),
	Long:  i18n.Translate("stop_help", "Stop services"),
	Run: func(cmd *cobra.Command, args []string) {
		if err := compose.Stop(global.Conf.System.ComposeFiles, args); err != nil {
			color.PrintRed(i18n.Tf("stop_err", "Stop Container error: {{ .Err }}", map[string]any{"Err": err.Error()}))
		}
		// fmt.Println("stdout--->", stdout)
	},
}

func init() {
	rootCmd.AddCommand(stopCmd)
}
