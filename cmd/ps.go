package cmd

import (
	"bypctl/pkg/color"
	"bypctl/pkg/compose"
	"bypctl/pkg/global"
	"bypctl/pkg/i18n"
	"github.com/spf13/cobra"
)

// psCmd represents the version command
var psCmd = &cobra.Command{
	Use:   "ps [SERVICE...]",
	Short: i18n.Translate(`ps_help`),
	Long:  i18n.Translate(`ps_help`),
	Run: func(cmd *cobra.Command, args []string) {
		if err := compose.Ps(global.Conf.System.ComposeFiles, args); err != nil {
			color.PrintRed(i18n.Tf("ps_err", err.Error()))
		}
	},
}

func init() {
	rootCmd.AddCommand(psCmd)
}
