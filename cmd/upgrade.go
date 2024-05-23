package cmd

import (
	"bypctl/pkg/color"
	"bypctl/pkg/files"
	"bypctl/pkg/i18n"
	"context"
	"github.com/gogf/gf/v2/crypto/gmd5"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/spf13/cobra"
	"os"
	"runtime"
	"strings"
)

// upgradeCmd represents the version command
var upgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: i18n.Translate("upgrade_help", "Upgrade bypanel"),
	Long:  i18n.Translate("upgrade_help", "Upgrade bypanel"),
	Run: func(cmd *cobra.Command, args []string) {
		f := files.NewFile()
		bypBin := "/usr/local/bin/bypctl"
		if !f.IsExist(bypBin) {
			color.PrintRed(i18n.Translate("upgrade_err", "Upgrade failed, bypanel does not exist"))
			os.Exit(1)
		}
		bypMd5Now := gmd5.MustEncryptFile(bypBin)
		panelUrl := "http://mirrors.linuxeye.com"
		resp := g.Client().GetContent(context.TODO(), panelUrl+"/md5sum.txt", nil)
		bypRemoteBinFile := "bypctl-linux-" + runtime.GOARCH
		panelRemoteUrl := panelUrl + "/bypanel/" + bypRemoteBinFile
		var bypMd5Remote string
		for _, line := range strings.Split(resp, "\n") {
			if strings.Contains(line, bypRemoteBinFile) {
				result := strings.Fields(line)
				bypMd5Remote = result[0]
			}
		}
		if bypMd5Now != bypMd5Remote {
			bypBinTmp := "/tmp/bypctl"
			color.PrintYellow(i18n.Translate("upgrade_ing", "Upgrade ing..."))
			if err := f.DownloadFile(panelRemoteUrl, bypBinTmp); err != nil {
				color.PrintRed(err.Error())
				os.Exit(1)
			}
			if err := f.Chmod(bypBinTmp, 0755); err != nil {
				color.PrintRed(err.Error())
				os.Exit(1)
			}
			if err := f.Rename(bypBinTmp, bypBin); err != nil {
				color.PrintRed(err.Error())
				os.Exit(1)
			}
			color.PrintGreen(i18n.Translate("upgrade_succ", "Congratulations! bypanel upgrade successful "))
		} else {
			color.PrintYellow(i18n.Translate("upgrade_warn", "Your bypanel already has the latest version or does not need to be upgraded! "))
		}
	},
}

func init() {
	rootCmd.AddCommand(upgradeCmd)
}
