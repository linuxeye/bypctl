package cmd

import (
	"bypctl/pkg/color"
	"bypctl/pkg/compose"
	"bypctl/pkg/global"
	"bypctl/pkg/i18n"
	"github.com/spf13/cobra"
	"strings"
)

// execCmd represents the exec command
var execCmd = &cobra.Command{
	Use:   "exec SERVICE COMMAND [ARGS...]",
	Short: i18n.Translate("exec_help", "Execute a command in a running container"),
	Long:  i18n.Translate("exec_help", "Execute a command in a running container"),
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		service := args[0]
		cmdStr := strings.Join(args[1:], " ")
		if err := compose.Exec(global.Conf.System.ComposeFiles, service, cmdStr); err != nil {
			color.PrintRed(i18n.Tf("exec_docker_err", "Exec {{ .Apps }}, error: {{ .Err }}", map[string]any{"Apps": global.Conf.System.ComposeProfiles, "Err": err.Error()}))
		}
		// fmt.Println(stdout)
	},
}

func init() {
	rootCmd.AddCommand(execCmd)
	execCmd.Flags().SetInterspersed(false)
}
