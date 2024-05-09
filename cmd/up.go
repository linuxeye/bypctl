package cmd

import (
	"bypctl/pkg/color"
	"bypctl/pkg/compose"
	"bypctl/pkg/global"
	"bypctl/pkg/i18n"
	"github.com/spf13/cobra"
)

var detachMode bool

// upCmd represents the version command
var upCmd = &cobra.Command{
	Use:   "up [OPTIONS] [SERVICE...]",
	Short: i18n.Translate("up_help", "Create and start containers"),
	Long:  i18n.Translate("up_help", "Create and start containers"),
	Run: func(cmd *cobra.Command, args []string) {
		if err := compose.Up(global.Conf.System.ComposeFiles, args, detachMode); err != nil {
			color.PrintRed(i18n.Tf("start_err", global.Conf.System.ComposeProfiles, err.Error()))
		}
	},
}

func init() {
	upCmd.Flags().BoolVarP(&detachMode, "detach", "d", false, i18n.Translate("up_detach_help", "Detached mode: Run containers in the background"))
	rootCmd.AddCommand(upCmd)
}
