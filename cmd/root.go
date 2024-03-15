package cmd

import (
	"bypctl/pkg/color"
	"bypctl/pkg/constant"
	"bypctl/pkg/database"
	"bypctl/pkg/files"
	"bypctl/pkg/global"
	"bypctl/pkg/i18n"
	"bypctl/pkg/log"
	"bypctl/pkg/migration"
	"bypctl/pkg/util"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"strings"
)

var (
	cfgFile string
	lang    string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "bypctl",
	Short: "bypctl",
	Long:  i18n.Translate(`root_help`),
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		// fmt.Println("global.Conf.System.Lang,---->", global.Conf.System.Lang)
		// log.Init()
		// database.Init()
		// migration.Init()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.CompletionOptions.DisableDefaultCmd = true

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is /opt/bypanel/.env)")
	rootCmd.PersistentFlags().StringVarP(&lang, "lang", "l", "", "set language")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	f := files.NewFile()
	var useCfgFile string
	if cfgFile != "" {
		useCfgFile = cfgFile
	} else {
		useCfgFile = "/opt/bypanel/.env"
	}

	// 判断师傅存
	if !f.IsExist(useCfgFile) {
		color.PrintRed(i18n.Tf("config_exist_err", useCfgFile))
		os.Exit(1)
	}

	// 判断文件以.env结尾
	if !strings.HasSuffix(useCfgFile, ".env") {
		color.PrintRed(i18n.Tf("config_file_err", useCfgFile))
		os.Exit(1)
	}

	viper.SetConfigFile(useCfgFile)
	viper.AutomaticEnv() // read in environment variables that match
	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
		// Unmarshal 将配置文件转成对象
	}

	global.Conf.System.Mode = logrus.DebugLevel.String()
	global.Conf.Log.Level = logrus.DebugLevel.String()
	global.Conf.System.Lang = viper.GetString("LANG")
	global.Conf.System.BasePath = viper.GetString("BASE_PATH")
	global.Conf.System.LogPath = filepath.Join(global.Conf.System.BasePath, "log")
	global.Conf.System.EnvFile = viper.ConfigFileUsed()
	global.Conf.System.VolumePath = viper.GetString("VOLUME_PATH")
	global.Conf.System.Uid = viper.GetString("NEW_UID")
	global.Conf.System.Gid = viper.GetString("NEW_GID")
	global.Conf.System.Timezone = viper.GetString("TIMEZONE")
	global.Conf.System.ComposeProfiles = viper.GetString("COMPOSE_PROFILES")
	global.Conf.System.NginxVer = viper.GetString("NGINX_SERVER")
	global.Conf.System.MySQLVer = viper.GetString("MYSQL_SERVER")
	global.Conf.System.MySQLRootPwd = viper.GetString("MYSQL_ROOT_PASSWORD")
	global.Conf.System.PGSQLVer = viper.GetString("PGSQL_SERVER")
	global.Conf.System.PGSQLRootUser = viper.GetString("PGSQL_ROOT_USER")
	global.Conf.System.PGSQLRootPwd = viper.GetString("PGSQL_ROOT_PASSWORD")
	global.Conf.System.RedisVer = viper.GetString("REDIS_SERVER")
	global.Conf.System.MongoVer = viper.GetString("MONGO_SERVER")
	global.Conf.System.MemcachedVer = viper.GetString("MEMCACHED_SERVER")

	for _, p := range strings.Split(global.Conf.System.ComposeProfiles, ",") {
		var selectVer string
		switch p {
		case "nginx":
			selectVer = global.Conf.System.NginxVer
		case "mysql":
			selectVer = global.Conf.System.MySQLVer
		case "postgresql":
			selectVer = global.Conf.System.PGSQLVer
		case "redis":
			selectVer = global.Conf.System.RedisVer
		case "mongo":
			selectVer = global.Conf.System.MongoVer
		case "memcached":
			selectVer = global.Conf.System.MemcachedVer
		default:
			selectVer = "latest"
		}
		composeFile := filepath.Join(global.Conf.System.BasePath, "app", p, selectVer, "docker-compose.yml")
		// 判断php应用时，目录结构单独处理
		if util.SliceItemStrExist(constant.PHPs, p) {
			composeFile = filepath.Join(global.Conf.System.BasePath, "app", "php", p[3:4]+"."+p[4:], "docker-compose.yml")
		}
		global.Conf.System.ComposeFiles = append(global.Conf.System.ComposeFiles, composeFile)
	}
	log.Init()
	database.Init()
	migration.Init()
}
