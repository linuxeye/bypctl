package cmd

import (
	"bypctl/pkg/color"
	"bypctl/pkg/compose"
	"bypctl/pkg/constant"
	"bypctl/pkg/global"
	"bypctl/pkg/i18n"
	"bypctl/pkg/util"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

// reloadCmd represents the version command
var reloadCmd = &cobra.Command{
	Use:   "reload [SERVICE...]",
	Short: i18n.Translate("reload_help", "Reload Web service"),
	Long:  i18n.Translate("reload_help", "Reload Web service"),
	Run: func(cmd *cobra.Command, args []string) {
		var webs []string
		if len(args) > 0 {
			if !util.IsSubSlice(constant.Webs, args) {
				color.PrintRed(i18n.Tf("web_input_err", "Make web config, \nvalue range is: {{ .ValueRange }}. ", map[string]any{"ValueRange": constant.Webs}))
				os.Exit(1)
			}
			webs = args
		} else {
			for _, v := range strings.Split(global.Conf.System.ComposeProfiles, ",") {
				switch v {
				case "nginx", "openresty", "apache", "caddy", "pingora":
					webs = []string{v}
				}
			}
		}
		if len(webs) == 0 {
			color.PrintRed(i18n.Translate("web_err", "Parameter error or no web startup."))
			os.Exit(1)
		}

		if util.SliceItemStrExist([]string{"nginx", "openresty"}, webs[0]) {
			if err := compose.Exec(global.Conf.System.ComposeFiles, webs[0], "nginx -s reload"); err != nil {
				color.PrintRed(i18n.Tf("reload_err", "Reloaded error: {{ .Err }}", map[string]any{"Err": err.Error()}))
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(reloadCmd)
}
