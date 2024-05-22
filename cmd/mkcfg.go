package cmd

import (
	"bypctl/pkg/color"
	"bypctl/pkg/compose"
	"bypctl/pkg/constant"
	"bypctl/pkg/files"
	"bypctl/pkg/global"
	"bypctl/pkg/i18n"
	"bypctl/pkg/ssl"
	"bypctl/pkg/util"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/spf13/cobra"
	"io/fs"
	"net/mail"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// WebSiteTpl .
type WebSiteTpl struct {
	ServerName     string
	ServerAlias    []string
	WebRoot        string
	Email          string
	PHPVer         string
	SSL            bool
	AccessLog      bool
	Hotlink        bool
	HotlinkDomains string
	RedirectDomain bool
	HTTPToHTTPS    bool
	Rewrite        string
}

// mkCfgCmd represents the config command
var mkCfgCmd = &cobra.Command{
	Use:   "mkcfg [WEB...]",
	Short: i18n.Translate("mkcfg_help", "Make web config."),
	Long:  i18n.Translate("mkcfg_help", "Make web config."),
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

		webSite := new(WebSiteTpl)
		// 设置ssl
		inputSSLs := []uint{1, 2, 3}
		inputSSL := gconv.Uint(util.ReaderTf("mkcfg_ssl", "What Are You Doing?\n\t1. Use HTTP Only\n\t2. Use your own SSL Certificate and Key\n\t3. Use Let's Encrypt to Create SSL Certificate and Key\nPlease input the correct option (default: {{ .Option }}): ", map[string]any{"Option": 2}))
		if !util.SliceItemUintExist(inputSSLs, inputSSL) {
			color.PrintRed(i18n.Tf("mkcfg_input_err", "Input error, the value should be set within the range: {{ .ValueRange }}", map[string]any{"ValueRange": inputSSL}))
			os.Exit(1)
		}

		switch inputSSL {
		case 2, 3:
			webSite.SSL = true
		}

		// 选择php版本
		_, _, PHPs := util.CompareStrList(strings.Split(global.Conf.System.ComposeProfiles, ","), constant.PHPs)
		if len(PHPs) == 1 {
			webSite.PHPVer = PHPs[0]
		} else if len(PHPs) > 1 {
			defaultPHP := PHPs[0]
			inputPHP := util.ReaderTf("mkcfg_select_php", "\nPlease enter the PHP version, the value range is: {{ .ValueRange }}, (default: {{ .DefaultValue }}): ", map[string]any{"ValueRange": PHPs, "DefaultValue": defaultPHP})
			if len(inputPHP) == 0 {
				inputPHP = defaultPHP
			}
			if util.SliceItemStrExist(PHPs, inputPHP) {
				webSite.PHPVer = inputPHP
			}
		}

		inputDomains := strings.Split(util.ReaderTf("mkcfg_domain", "\nPlease enter the domains (separate multiple domains with commas, e.g.: bypanel.com,www.bypanel.com): ", nil), ",")
		for i := range inputDomains {
			inputDomains[i] = strings.TrimSpace(inputDomains[i])
		}

		if !util.ValidateDomains(inputDomains) {
			color.PrintRed(i18n.Translate("mkcfg_domain_err", "Invalid domain format."))
			os.Exit(1)
		}

		color.PrintGreen(i18n.Tf("mkcfg_domain_list", "Domains: {{ .Domains }}", map[string]any{"Domains": inputDomains}))
		webSite.ServerName = inputDomains[0]
		if len(inputDomains) > 1 {
			webSite.ServerAlias = inputDomains[1:]
		}

		defaultWebroot := filepath.Join(global.Conf.System.VolumePath, "webroot", inputDomains[0])
		inputWebroot := util.ReaderTf("mkcfg_webroot", "\nPlease enter the webroot directory of the domains(default: {{ .DefaultValue }}): ", map[string]any{"DefaultValue": defaultWebroot})
		if len(inputWebroot) == 0 {
			inputWebroot = defaultWebroot
		}

		if !strings.HasPrefix(inputWebroot, filepath.Join(global.Conf.System.VolumePath, "webroot")) {
			color.PrintRed(i18n.Translate("reader_path_err", "Input error, wrong path format"))
			os.Exit(1)
		}
		webSite.WebRoot = filepath.Join("/var/www", filepath.Base(inputWebroot))
		color.PrintGreen(i18n.Tf("mkcfg_domain_webroot", "\nDomain webroot directory: {{ .WebRoot }}", map[string]any{"WebRoot": inputWebroot}))
		color.PrintYellow(i18n.Translate("mkcfg_webroot_permission", "\nCreate webroot directory......\nset permissions of webroot directory......"))

		f := files.NewFile()
		if err := f.CreateDir(inputWebroot, fs.ModeDir); err != nil {
			color.PrintRed(err.Error())
			// os.Exit(1)
		}
		if err := f.ChownR(inputWebroot, global.Conf.System.Uid, global.Conf.System.Gid, true); err != nil {
			color.PrintRed(err.Error())
			// os.Exit(1)
		}

		// Apache ServerAdmin
		if webs[0] == "apache" {
			email := util.ReaderTf("mkcfg_email", "\nPlease enter your Email: ", nil)
			_, err := mail.ParseAddress(email)
			if err != nil {
				color.PrintRed(err.Error())
				os.Exit(1)
			}
			webSite.Email = email
		}

		if util.SliceItemUintExist([]uint{2, 3}, inputSSL) {
			if util.SliceItemUintExist([]uint{2, 3}, inputSSL) {
				country := util.ReaderTf("mkcfg_self_ssl_country", "\nYou are about to be asked to enter information that will be incorporated into your certificate request.\nWhat you are about to enter is what is called a Distinguished Name or a DN.\nThere are quite a few fields but you can leave some blank\nFor some fields there will be a default value,\nIf you enter '.', the field will be left blank.\n\nCountry Name (2 letter code) [CN]: ", nil)
				if len(country) == 0 {
					country = "CN"
				}
				province := util.ReaderTf("mkcfg_self_ssl_province", "\nState or Province Name (full name) [Shanghai]: ", nil)
				if len(province) == 0 {
					province = "Shanghai"
				}

				city := util.ReaderTf("mkcfg_self_ssl_city", "\nLocality Name (eg, city) [Shanghai]: ", nil)
				if len(city) == 0 {
					city = "Shanghai"
				}

				organization := util.ReaderTf("mkcfg_self_ssl_organization", "\nOrganization Name (eg, company) [Example Inc.]: ", nil)
				if len(organization) == 0 {
					organization = "Example Inc."
				}

				organizationUint := util.ReaderTf("mkcfg_self_ssl_organizationuint", "\nOrganizational Unit Name (eg, section) [IT Dept.]: ", nil)
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
					KeyType:          "2048",
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
			// 域名重定向
			var inputRedirectFlag string
			if len(inputDomains) >= 2 {
				if util.SliceItemUintExist([]uint{2, 3}, inputSSL) {
					inputRedirectFlag = strings.ToLower(util.ReaderTf("mkcfg_redirect_flag", "\nDo you want to redirect from {{ .MoreDomains }} to {{ .Domain }}？[y/n]: ", map[string]any{"MoreDomains": inputDomains[1:], "Domain": inputDomains[0]}))
					if !util.SliceItemStrExist(constant.FlagYN, inputRedirectFlag) {
						color.PrintRed(i18n.Tf("reader_input_err", "Input error, the value should be set within the range: {{ .ValueRange }}", map[string]any{"ValueRange": constant.FlagYN}))
					}
				}
			}
			if inputRedirectFlag == "y" {
				webSite.RedirectDomain = true
			}

			// https跳转
			var inputToHttpsFlag string
			if util.SliceItemUintExist([]uint{2, 3}, inputSSL) {
				inputToHttpsFlag = strings.ToLower(util.ReaderTf("mkcfg_to_https_flag", "\nDo you want to redirect all HTTP requests to HTTPS? [y/n]: ", nil))
				if !util.SliceItemStrExist(constant.FlagYN, inputToHttpsFlag) {
					color.PrintRed(i18n.Tf("reader_input_err", "Input error, the value should be set within the range: {{ .ValueRange }}", map[string]any{"ValueRange": constant.FlagYN}))
				}
			}

			if inputToHttpsFlag == "y" {
				webSite.HTTPToHTTPS = true
			}

			// 防盗链
			antiHotlinkFlag := strings.ToLower(util.ReaderTf("mkcfg_hotlink", "\nDo you want to add hotlink protection? [y/n]: ", nil))
			if !util.SliceItemStrExist(constant.FlagYN, antiHotlinkFlag) {
				color.PrintRed(i18n.Tf("reader_input_err", "Input error, the value should be set within the range: {{ .ValueRange }}", map[string]any{"ValueRange": constant.FlagYN}))
				os.Exit(1)
			}

			if antiHotlinkFlag == "y" {
				webSite.Hotlink = true
				webSite.HotlinkDomains = strings.Join(util.GetUniqueDomains(inputDomains), " ")
			}

			inputRewriteName := strings.ToLower(util.ReaderTf("mkcfg_rewrite", "\nPlease input the rewrite of programme :\nwordpress,opencart,magento2,drupal,joomla,codeigniter,laravel\nthinkphp,discuz,typecho,ecshop,nextcloud,zblog,whmcs rewrite was exist.\n(Default rewrite: none): ", nil))
			if len(inputRewriteName) == 0 {
				inputRewriteName = "none"
			}
			rewriteFile := filepath.Join(global.Conf.System.BasePath, "cfg", webs[0], "rewrite", inputRewriteName+".conf")
			if !f.IsExist(rewriteFile) {
				color.PrintYellow(i18n.Tf("mkcfg_rewrite_file", "\nThe rewrite file does not exist, it will create {{ .RewriteFile }}, please add the rewrite rules to this file.", map[string]any{"RewriteFile": rewriteFile}))
				if err := f.CreateFile(rewriteFile); err != nil {
					color.PrintRed(err.Error())
					os.Exit(1)
				}
			}
			webSite.Rewrite = inputRewriteName

			nginxLogFlag := strings.ToLower(util.ReaderTf("mkcfg_nginx_log", "\nDo you want to enable the {{ .App }} access log? [y/n]: ", map[string]any{"App": webs[0]}))
			if !util.SliceItemStrExist(constant.FlagYN, nginxLogFlag) {
				color.PrintRed(i18n.Tf("reader_input_err", "Input error, the value should be set within the range: {{ .ValueRange }}", map[string]any{"ValueRange": constant.FlagYN}))
				os.Exit(1)
			}

			if nginxLogFlag == "y" {
				webSite.AccessLog = true
			}

			nginxConf := `server {
  listen 80;
  listen [::]:80;
  {{- if .SSL }}
  listen 443 ssl;
  listen [::]:443 ssl;
  ssl_certificate cert/{{ .ServerName }}.crt;
  ssl_certificate_key cert/{{ .ServerName }}.key;
  ssl_protocols TLSv1.2 TLSv1.3;
  ssl_ecdh_curve X25519:prime256v1:secp384r1:secp521r1;
  ssl_ciphers ECDHE-ECDSA-AES256-GCM-SHA384:ECDHE-RSA-AES256-GCM-SHA384:ECDHE-ECDSA-CHACHA20-POLY1305:ECDHE-RSA-CHACHA20-POLY1305:ECDHE-ECDSA-AES128-GCM-SHA256:ECDHE-RSA-AES128-GCM-SHA256:DHE-RSA-AES256-GCM-SHA384:DHE-RSA-AES128-GCM-SHA256;
  ssl_conf_command Ciphersuites TLS_AES_256_GCM_SHA384:TLS_CHACHA20_POLY1305_SHA256:TLS_AES_128_GCM_SHA256;
  ssl_conf_command Options PrioritizeChaCha;
  ssl_prefer_server_ciphers on;
  ssl_session_timeout 10m;
  ssl_session_cache shared:SSL:10m;
  ssl_buffer_size 2k;
  add_header Strict-Transport-Security max-age=15768000;
  ssl_stapling on;
  ssl_stapling_verify on;
  {{- end }}
  server_name {{ .ServerName }} {{- range .ServerAlias }} {{.}} {{- end}};
  {{- if .AccessLog }}
  access_log /var/log/nginx/{{ .ServerName }}_nginx.log combined;
  {{- else }}
  access_log off;
  {{- end }}
  index index.html index.htm index.php;
  root {{ .WebRoot }};
  {{ if .HTTPToHTTPS }}
  if ($ssl_protocol = "") { return 301 https://$host$request_uri; }
  {{ end }}
  {{- if .RedirectDomain }}
  if ($host != {{ .ServerName }}) {  return 301 $scheme://{{ .ServerName }}$request_uri;  }
  {{ end }}
  #error_page 404 /404.html;
  #error_page 502 /502.html;
  {{ if .Hotlink }}
  location ~ .*\.(wma|wmv|asf|mp3|mmf|zip|rar|jpg|gif|png|swf|flv|mp4)$ {
    valid_referers none blocked {{ .HotlinkDomains }};
    if ($invalid_referer) {
        return 403;
    }
  }
  {{ end }}
  include rewrite/{{ .Rewrite }}.conf;
  {{ if .PHPVer }}
  include enable-{{ .PHPVer }}.conf;
  {{ end }}
  location ~ .*\.(gif|jpg|jpeg|png|bmp|swf|flv|mp4|ico)$ {
    expires 30d;
    access_log off;
  }
  location ~ .*\.(js|css)?$ {
    expires 7d;
    access_log off;
  }
  location ~ /(\.user\.ini|\.ht|\.git|\.svn|\.project|LICENSE|README\.md) {
    deny all;
  }
  location /.well-known {
    allow all;
  }
}
`
			nginxTpl, err := template.New("ngx").Parse(nginxConf)
			if err != nil {
				color.PrintRed(err.Error())
			}

			fn, err := os.Create(filepath.Join(global.Conf.System.BasePath, "cfg", webs[0], "conf.d", webSite.ServerName+".conf"))
			if err != nil {
				color.PrintRed(err.Error())
			}
			defer fn.Close()

			// fmt.Println("ngx----->", gconv.String(ngx))
			if err = nginxTpl.Execute(fn, webSite); err != nil {
				color.PrintRed(err.Error())
			}

			if err = compose.Exec(global.Conf.System.ComposeFiles, webs[0], "nginx -t"); err != nil {
				color.PrintRed(err.Error())
			} else {
				compose.Exec(global.Conf.System.ComposeFiles, webs[0], "nginx -s reload")
			}
		}
		if webs[0] == "apache" {
			apacheConf := `<VirtualHost *:80>
  ServerAdmin {{ .Email }}
  DocumentRoot {{ .WebRoot }}
  ServerName {{ .ServerName }}
  ServerAlias {{- range .ServerAlias }} {{.}} {{- end}}
  {{- if .AccessLog }}
  ErrorLog /var/log/httpd/{{ .ServerName }}_error_apache.log
  CustomLog /var/log/httpd/{{ .ServerName }}_apache.log common
  {{- else }}
  ErrorLog /dev/null
  CustomLog /dev/null common
  {{- end }}
  <Files ~ (\.user.ini|\.htaccess|\.git|\.svn|\.project|LICENSE|README.md)$>
    Order allow,deny
    Deny from all
  </Files>
  <FilesMatch \.php$>
    SetHandler proxy:fcgi://{{ .PHPVer }}:9000
  </FilesMatch>
<Directory {{ .WebRoot }}>
  SetOutputFilter DEFLATE
  Options FollowSymLinks ExecCGI
  Require all granted
  AllowOverride All
  Order allow,deny
  Allow from all
  DirectoryIndex index.html index.php
</Directory>
</VirtualHost>
{{- if .SSL }}
<VirtualHost *:443>
  ServerAdmin {{ .Email }}
  DocumentRoot {{ .WebRoot }}
  ServerName {{ .ServerName }}
  ServerAlias {{- range .ServerAlias }} {{.}} {{- end}}
  SSLEngine on
  SSLCertificateFile /etc/httpd/cert/{{ .ServerName }}.crt
  SSLCertificateKeyFile /etc/httpd/cert/{{ .ServerName }}.key
  {{- if .AccessLog }}
  ErrorLog /var/log/httpd/{{ .ServerName }}_error_apache.log
  CustomLog /var/log/httpd/{{ .ServerName }}_apache.log common
  {{- else }}
  ErrorLog /dev/null
  CustomLog /dev/null common
  {{- end }}
  <Files ~ (\.user.ini|\.htaccess|\.git|\.svn|\.project|LICENSE|README.md)$>
    Order allow,deny
    Deny from all
  </Files>
  <FilesMatch \.php$>
    SetHandler proxy:fcgi://{{ .PHPVer }}:9000
  </FilesMatch>
<Directory {{ .WebRoot }}>
  SetOutputFilter DEFLATE
  Options FollowSymLinks ExecCGI
  Require all granted
  AllowOverride All
  Order allow,deny
  Allow from all
  DirectoryIndex index.html index.php
</Directory>
</VirtualHost>
{{- end }}
`
			apacheTpl, err := template.New("apache").Parse(apacheConf)
			if err != nil {
				color.PrintRed(err.Error())
			}

			fn, err := os.Create(filepath.Join(global.Conf.System.BasePath, "cfg", webs[0], "conf.d", webSite.ServerName+".conf"))
			if err != nil {
				color.PrintRed(err.Error())
			}
			defer fn.Close()

			if err = apacheTpl.Execute(fn, webSite); err != nil {
				color.PrintRed(err.Error())
			}

			if err = compose.Exec(global.Conf.System.ComposeFiles, webs[0], "apachectl -t"); err != nil {
				color.PrintRed(err.Error())
			} else {
				compose.Exec(global.Conf.System.ComposeFiles, webs[0], "apachectl -k graceful")
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(mkCfgCmd)
}
