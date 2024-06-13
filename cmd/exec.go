package cmd

import (
	"bypctl/pkg/color"
	"bypctl/pkg/compose"
	"bypctl/pkg/global"
	"bypctl/pkg/i18n"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

// execCmd represents the exec command
var execCmd = &cobra.Command{
	Use:   "exec SERVICE COMMAND [ARGS...]",
	Short: i18n.Translate("exec_help", "Execute a command in a running container"),
	Long:  i18n.Translate("exec_help", "Execute a command in a running container"),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			color.PrintRed(i18n.Translate("exec_args_err", "At least 2 parameters, the first parameter is the service name, and the second parameter is the command."))
			os.Exit(1)
		}
		if err := compose.Exec(global.Conf.System.ComposeFiles, args[0], strings.Join(args[1:], " ")); err != nil {
			color.PrintRed(i18n.Tf("exec_docker_err", "Exec {{ .Apps }}, error: {{ .Err }}", map[string]any{"Apps": global.Conf.System.ComposeProfiles, "Err": err.Error()}))
		}
		// fmt.Println(stdout)
	},
}

func init() {
	rootCmd.AddCommand(execCmd)
}
