package cmd

import (
	"bypctl/pkg/color"
	"bypctl/pkg/compose"
	"bypctl/pkg/global"
	"bypctl/pkg/i18n"
	"github.com/spf13/cobra"
)

// statusCmd represents the version command
var statusCmd = &cobra.Command{
	Use:   "status [SERVICE...]",
	Short: i18n.Translate("ps_help", "List containers"),
	Long:  i18n.Translate("ps_help", "List containers"),
	Run: func(cmd *cobra.Command, args []string) {
		if err := compose.Ps(global.Conf.System.ComposeFiles, args); err != nil {
			color.PrintRed(i18n.Tf("ps_err", "List containers error: {{ .Err }}", map[string]any{"Err": err.Error()}))
		}
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
