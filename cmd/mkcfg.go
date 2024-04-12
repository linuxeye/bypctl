package cmd

import (
	"bypctl/pkg/color"
	"bypctl/pkg/constant"
	"bypctl/pkg/files"
	"bypctl/pkg/global"
	"bypctl/pkg/i18n"
	"bypctl/pkg/ssl"
	"bypctl/pkg/util"
	"fmt"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/spf13/cobra"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// mkCfgCmd represents the config command
var mkCfgCmd = &cobra.Command{
	Use:   "mkcfg [WEB...]",
	Short: i18n.Translate(`mkcfg_help`),
	Long:  i18n.Translate(`mkcfg_help`),
	Run: func(cmd *cobra.Command, args []string) {
		var webs []string
		if len(args) > 0 {
			if !util.IsSubSlice(constant.Webs, args) {
				color.PrintRed(i18n.Tf("mkcfg_web_input_err", constant.Webs))
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
			color.PrintRed(i18n.Translate("mkcfg_web_err"))
			os.Exit(1)
		}

		// 设置语言
		inputSSLs := []uint{1, 2, 3}
		inputSSL := gconv.Uint(util.ReaderTf("mkcfg_ssl", 2))
		if !util.SliceItemUintExist(inputSSLs, inputSSL) {
			color.PrintRed(i18n.Tf("mkcfg_input_err", inputSSLs))
			os.Exit(1)
		}

		inputDomains := strings.Split(util.ReaderTf("mkcfg_domain"), ",")
		for i := range inputDomains {
			inputDomains[i] = strings.TrimSpace(inputDomains[i])
		}

		if !util.ValidateDomains(inputDomains) {
			color.PrintRed(i18n.Translate("mkcfg_domain_err"))
			os.Exit(1)
		}

		color.PrintGreen(i18n.Tf("mkcfg_domain_list", inputDomains))

		defaultWebroot := filepath.Join(global.Conf.System.VolumePath, "webroot", inputDomains[0])
		inputWebroot := util.ReaderTf("mkcfg_webroot", defaultWebroot)
		if len(inputWebroot) == 0 {
			inputWebroot = defaultWebroot
		}

		if !strings.HasPrefix(inputWebroot, filepath.Join(global.Conf.System.VolumePath, "webroot")) {
			color.PrintRed(i18n.Translate("reader_path_err"))
			os.Exit(1)
		}

		color.PrintGreen(i18n.Tf("mkcfg_domain_webroot", inputWebroot))
		color.PrintYellow(i18n.Translate("mkcfg_webroot_permission"))

		f := files.NewFile()
		if err := f.CreateDir(inputWebroot, fs.ModeDir); err != nil {
			color.PrintRed(err.Error())
			// os.Exit(1)
		}
		if err := f.ChownR(inputWebroot, global.Conf.System.Uid, global.Conf.System.Gid, true); err != nil {
			color.PrintRed(err.Error())
			// os.Exit(1)
		}

		// 域名重定向
		var inputRedirectFlag string
		if len(inputDomains) >= 2 {
			if util.SliceItemUintExist([]uint{2, 3}, inputSSL) {
				inputRedirectFlag = strings.ToLower(util.ReaderTf("mkcfg_redirect_flag", inputDomains[1:], inputDomains[0]))
				if !util.SliceItemStrExist(constant.FlagYN, inputRedirectFlag) {
					color.PrintRed(i18n.Tf("reader_input_err", constant.FlagYN))
				}
			}
		}

		// https跳转
		var inputToHttpsFlag string
		if util.SliceItemUintExist([]uint{2, 3}, inputSSL) {
			inputToHttpsFlag = strings.ToLower(util.ReaderTf("mkcfg_to_https_flag"))
			if !util.SliceItemStrExist(constant.FlagYN, inputToHttpsFlag) {
				color.PrintRed(i18n.Tf("reader_input_err", constant.FlagYN))
			}
		}

		if util.SliceItemUintExist([]uint{2, 3}, inputSSL) {

			// var (
			// 	country          string
			// 	organization     string
			// 	organizationUint string
			// 	province         string
			// 	city             string
			// )
			if util.SliceItemUintExist([]uint{2, 3}, inputSSL) {
				country := util.ReaderTf("mkcfg_self_ssl_country")
				if len(country) == 0 {
					country = "CN"
				}
				province := util.ReaderTf("mkcfg_self_ssl_province")
				if len(province) == 0 {
					province = "Shanghai"
				}

				city := util.ReaderTf("mkcfg_self_ssl_city")
				if len(city) == 0 {
					city = "Shanghai"
				}

				organization := util.ReaderTf("mkcfg_self_ssl_organization")
				if len(organization) == 0 {
					organization = "Example Inc."
				}

				organizationUint := util.ReaderTf("mkcfg_self_ssl_organizationuint")
				if len(organizationUint) == 0 {
					organizationUint = "IT Dept."
				}

				if err := ssl.GenerateSelfPem(ssl.SelfSSL{
					Domains:          inputDomains,
					CommonName:       inputDomains[0],
					Country:          country,
					Organization:     organization,
					OrganizationUint: organizationUint,
					Name:             inputDomains[0],
					KeyType:          "P256",
					Province:         province,
					City:             city,
					CertificatePath:  filepath.Join(global.Conf.System.BasePath, "cfg", webs[0], "cert", inputDomains[0]+".crt"),
					PrivateKeyPath:   filepath.Join(global.Conf.System.BasePath, "cfg", webs[0], "cert", inputDomains[0]+".key"),
				}); err != nil {
					color.PrintRed(err.Error())
				}
			}
		}

		if util.SliceItemStrExist([]string{"nginx", "openresty"}, webs[0]) {
			antiHotlinkFlag := strings.ToLower(util.ReaderTf("mkcfg_hotlink"))
			if !util.SliceItemStrExist(constant.FlagYN, antiHotlinkFlag) {
				color.PrintRed(i18n.Tf("reader_input_err", constant.FlagYN))
				os.Exit(1)
			}

			inputRewriteFlag := strings.ToLower(util.ReaderTf("mkcfg_rewrite_flag"))
			if !util.SliceItemStrExist(constant.FlagYN, inputRewriteFlag) {
				color.PrintRed(i18n.Tf("reader_input_err", constant.FlagYN))
			}
		}

		domains := util.GetUniqueDomains(inputDomains)
		fmt.Println("domains---->", domains)

	},
}

func init() {
	rootCmd.AddCommand(mkCfgCmd)
}
