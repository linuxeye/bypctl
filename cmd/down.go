package cmd

import (
	"bypctl/pkg/color"
	"bypctl/pkg/compose"
	"bypctl/pkg/global"
	"bypctl/pkg/i18n"
	"github.com/spf13/cobra"
)

// downCmd represents the version command
var downCmd = &cobra.Command{
	Use:   "down [SERVICE...]",
	Short: i18n.Translate("down_help", "Stop and remove containers, networks"),
	Long:  i18n.Translate("down_help", "Stop and remove containers, networks"),
	Run: func(cmd *cobra.Command, args []string) {
		if err := compose.Down(global.Conf.System.ComposeFiles, args); err != nil {
			color.PrintRed(i18n.Tf("down_err", "Stop and remove error: {{ .Err }}", map[string]any{"Err": err.Error()}))
		}
	},
}

func init() {
	rootCmd.AddCommand(downCmd)
}
