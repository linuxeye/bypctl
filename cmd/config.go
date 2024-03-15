package cmd

import (
	"bypctl/pkg/color"
	"bypctl/pkg/constant"
	"bypctl/pkg/files"
	"bypctl/pkg/global"
	"bypctl/pkg/i18n"
	"bypctl/pkg/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: i18n.Translate(`config_help`),
	Long:  i18n.Translate(`config_help`),
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {

		// 设置语言
		formatLang := []string{"en", "zh-CN"}
		inputLang := util.ReaderTf("config_select_lang", formatLang, global.Conf.System.Lang)
		if len(inputLang) == 0 {
			inputLang = global.Conf.System.Lang
		}
		if !util.SliceItemStrExist(formatLang, inputLang) {
			color.PrintRed(i18n.Tf("reader_input_err", formatLang))
			os.Exit(1)
		}
		if global.Conf.System.Lang != inputLang {
			viper.Set("LANG", inputLang)
			global.Conf.System.Lang = inputLang
		}

		// 设置bypanel路径
		inputBasePath := util.ReaderTf("config_bypanel_path", global.Conf.System.BasePath)
		if len(inputBasePath) == 0 {
			inputBasePath = global.Conf.System.BasePath
		}

		if !strings.HasPrefix(inputBasePath, "/") {
			color.PrintRed(i18n.Translate("reader_path_err"))
			os.Exit(1)
		}

		if global.Conf.System.BasePath != inputBasePath {
			viper.Set("BYPANEL_PATH", inputBasePath)
			global.Conf.System.BasePath = inputBasePath
			// to do mv path
		}

		// 设置volume路径
		inputVolumePath := util.ReaderTf("config_volume_path", global.Conf.System.VolumePath)
		if len(inputVolumePath) == 0 {
			inputVolumePath = global.Conf.System.VolumePath
		}

		if !strings.HasPrefix(inputVolumePath, "/") {
			color.PrintRed(i18n.Translate("reader_path_err"))
			os.Exit(1)
		}

		if global.Conf.System.VolumePath != inputVolumePath {
			viper.Set("VOLUME_PATH", inputVolumePath)
			global.Conf.System.VolumePath = inputVolumePath
			// to do mv path
		}

		// Timezone for Docker container
		inputTimezone := util.ReaderTf("config_timezone", global.Conf.System.Timezone)
		if len(inputTimezone) == 0 {
			inputTimezone = global.Conf.System.Timezone
		}

		_, err := time.LoadLocation(inputTimezone)
		if err != nil {
			color.PrintRed(i18n.Tf("reader_msg_err", err.Error()))
			os.Exit(1)
		}

		if global.Conf.System.Timezone != inputTimezone {
			viper.Set("TIMEZONE", inputTimezone)
			global.Conf.System.Timezone = inputTimezone
		}

		// 设置COMPOSE_PROFILES
		inputApps := util.ReaderTf("config_compose_profiles", global.Conf.System.ComposeProfiles)
		if len(inputApps) == 0 {
			inputApps = global.Conf.System.ComposeProfiles
		}

		appList := strings.Split(inputApps, ",")

		// 判断应用输入是否正确
		f := files.NewFile()
		for _, v := range constant.PHPs {
			appDir := filepath.Join(global.Conf.System.BasePath, "app", v)
			// 判断php应用时，目录结构单独处理
			if util.SliceItemStrExist(constant.PHPs, v) {
				appDir = filepath.Join(global.Conf.System.BasePath, "app", "php", v[3:4]+"."+v[4:])
			}

			if !f.IsDir(appDir) {
				color.PrintRed(i18n.Tf("config_app_err", v))
				os.Exit(1)
			}
		}

		// 判断web只能选择1个
		if util.CheckElements([]string{"nginx", "openresty", "apache"}, appList) {
			color.PrintRed(i18n.Translate("config_web_err"))
			os.Exit(1)
		}

		if global.Conf.System.ComposeProfiles != inputApps {
			viper.Set("COMPOSE_PROFILES", inputApps)
			global.Conf.System.ComposeProfiles = inputApps
		}

		// Nginx
		if util.SliceItemStrExist(appList, "nginx") {
			ngxDir := filepath.Join(global.Conf.System.BasePath, "app", "nginx")
			ngxVers, err := files.GetSubFileNames(ngxDir)
			if err != nil {
				color.PrintRed(i18n.Tf("reader_msg_err", err.Error()))
				os.Exit(1)
			}
			inputNgxVer := util.ReaderTf("config_select_ver", "Nginx", ngxVers, global.Conf.System.NginxVer)
			if len(inputNgxVer) == 0 {
				inputNgxVer = global.Conf.System.NginxVer
			}

			if !f.IsDir(filepath.Join(ngxDir, inputNgxVer)) {
				color.PrintRed(i18n.Tf("config_ver_err", inputNgxVer))
				os.Exit(1)
			}

			if global.Conf.System.NginxVer != inputNgxVer {
				viper.Set("NGINX_SERVER", inputNgxVer)
				global.Conf.System.NginxVer = inputNgxVer
			}
		}

		// MySQL
		if util.SliceItemStrExist(appList, "mysql") {
			// 版本
			mySQLDir := filepath.Join(global.Conf.System.BasePath, "app", "mysql")
			mySQLVers, err := files.GetSubFileNames(mySQLDir)
			if err != nil {
				color.PrintRed(i18n.Tf("reader_msg_err", err.Error()))
				os.Exit(1)
			}

			if runtime.GOARCH == "arm64" {
				mySQLVers = util.RemoveStrings(mySQLVers, constant.MySQLNotArm)
			}

			inputMySQLVer := util.ReaderTf("config_select_ver", "MySQL", mySQLVers, global.Conf.System.MySQLVer)
			if len(inputMySQLVer) == 0 {
				inputMySQLVer = global.Conf.System.MySQLVer
			}

			if !f.IsDir(filepath.Join(mySQLDir, inputMySQLVer)) {
				color.PrintRed(i18n.Tf("config_ver_err", inputMySQLVer))
				os.Exit(1)
			}

			if global.Conf.System.MySQLVer != inputMySQLVer {
				viper.Set("MYSQL_SERVER", inputMySQLVer)
				global.Conf.System.MySQLVer = inputMySQLVer
			}

			// 密码
			inputDBRootPwd := util.ReaderTf("config_db_pwd", "MySQL", "root", global.Conf.System.MySQLRootPwd)
			if len(inputDBRootPwd) == 0 {
				inputDBRootPwd = global.Conf.System.MySQLRootPwd
			}

			if len(inputDBRootPwd) < 5 {
				color.PrintRed(i18n.Translate("config_db_pwd_err"))
				os.Exit(1)
			}

			if global.Conf.System.MySQLRootPwd != inputDBRootPwd {
				viper.Set("MYSQL_ROOT_PASSWORD", inputDBRootPwd)
				global.Conf.System.MySQLRootPwd = inputDBRootPwd
			}
		}

		// PostGreSQL
		if util.SliceItemStrExist(appList, "postgresql") {
			// 版本
			pgSQLDir := filepath.Join(global.Conf.System.BasePath, "app", "postgresql")
			pgsqlVers, err := files.GetSubFileNames(pgSQLDir)
			if err != nil {
				color.PrintRed(i18n.Tf("reader_msg_err", err.Error()))
				os.Exit(1)
			}
			inputPGSQLVer := util.ReaderTf("config_select_ver", "PostgreSQL", pgsqlVers, global.Conf.System.PGSQLVer)
			if len(inputPGSQLVer) == 0 {
				inputPGSQLVer = global.Conf.System.PGSQLVer
			}

			if !f.IsDir(filepath.Join(pgSQLDir, inputPGSQLVer)) {
				color.PrintRed(i18n.Tf("config_ver_err", inputPGSQLVer))
				os.Exit(1)
			}

			if global.Conf.System.PGSQLVer != inputPGSQLVer {
				viper.Set("PGSQL_SERVER", inputPGSQLVer)
				global.Conf.System.PGSQLVer = inputPGSQLVer
			}

			// 密码
			inputDBRootPwd := util.ReaderTf("config_db_pwd", "PGSQL", "postgres", global.Conf.System.PGSQLRootPwd)
			if len(inputDBRootPwd) == 0 {
				inputDBRootPwd = global.Conf.System.PGSQLRootPwd
			}

			if len(inputDBRootPwd) < 5 {
				color.PrintRed(i18n.Translate("config_db_pwd_err"))
				os.Exit(1)
			}

			if global.Conf.System.PGSQLRootPwd != inputDBRootPwd {
				viper.Set("PGSQL_ROOT_USER", inputDBRootPwd)
				global.Conf.System.PGSQLRootPwd = inputDBRootPwd
			}
		}

		// MongoDB
		if util.SliceItemStrExist(appList, "mongo") {
			// 版本
			mongoDir := filepath.Join(global.Conf.System.BasePath, "app", "mongo")
			mongoVers, err := files.GetSubFileNames(mongoDir)
			if err != nil {
				color.PrintRed(i18n.Tf("reader_msg_err", err.Error()))
				os.Exit(1)
			}
			inputMongoLVer := util.ReaderTf("config_select_ver", "MongoDB", mongoVers, global.Conf.System.MongoVer)
			if len(inputMongoLVer) == 0 {
				inputMongoLVer = global.Conf.System.MongoVer
			}

			if !f.IsDir(filepath.Join(mongoDir, inputMongoLVer)) {
				color.PrintRed(i18n.Tf("config_ver_err", inputMongoLVer))
				os.Exit(1)
			}

			if global.Conf.System.MongoVer != inputMongoLVer {
				viper.Set("MONGO_SERVER", inputMongoLVer)
				global.Conf.System.MongoVer = inputMongoLVer
			}
		}

		// Redis
		if util.SliceItemStrExist(appList, "redis") {
			// 版本
			redisDir := filepath.Join(global.Conf.System.BasePath, "app", "redis")
			redisVers, err := files.GetSubFileNames(redisDir)
			if err != nil {
				color.PrintRed(i18n.Tf("reader_msg_err", err.Error()))
				os.Exit(1)
			}
			inputRedisLVer := util.ReaderTf("config_select_ver", "Redis", redisVers, global.Conf.System.RedisVer)
			if len(inputRedisLVer) == 0 {
				inputRedisLVer = global.Conf.System.RedisVer
			}

			if !f.IsDir(filepath.Join(redisDir, inputRedisLVer)) {
				color.PrintRed(i18n.Tf("config_ver_err", inputRedisLVer))
				os.Exit(1)
			}

			if global.Conf.System.RedisVer != inputRedisLVer {
				viper.Set("REDIS_SERVER", inputRedisLVer)
				global.Conf.System.RedisVer = inputRedisLVer
			}
		}

		// Memcached
		if util.SliceItemStrExist(appList, "memcached") {
			// 版本
			memcachedDir := filepath.Join(global.Conf.System.BasePath, "app", "memcached")
			memcachedVers, err := files.GetSubFileNames(memcachedDir)
			if err != nil {
				color.PrintRed(i18n.Tf("reader_msg_err", err.Error()))
				os.Exit(1)
			}
			inputMemcachedLVer := util.ReaderTf("config_select_ver", "Memcached", memcachedVers, global.Conf.System.MemcachedVer)
			if len(inputMemcachedLVer) == 0 {
				inputMemcachedLVer = global.Conf.System.MemcachedVer
			}

			if !f.IsDir(filepath.Join(memcachedDir, inputMemcachedLVer)) {
				color.PrintRed(i18n.Tf("config_ver_err", inputMemcachedLVer))
				os.Exit(1)
			}

			if global.Conf.System.MemcachedVer != inputMemcachedLVer {
				viper.Set("MEMCACHED_SERVER", inputMemcachedLVer)
				global.Conf.System.MemcachedVer = inputMemcachedLVer
			}
		}

		// 更新配置文件
		viper.SetConfigType("env")
		if err := viper.WriteConfig(); err != nil {
			color.PrintRed(i18n.Tf("config_write_err", err.Error()))
		}
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}
