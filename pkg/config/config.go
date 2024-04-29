package config

type Config struct {
	System SystemConfig `json:"system"`
	Log    LogConfig    `json:"log"`
}

type SystemConfig struct {
	BasePath        string // 基础目录，默认/opt/bypanel
	PanelMd5        string
	LogPath         string
	Mode            string
	Lang            string   // 语言
	EnvFile         string   // Compose环境变量文件
	Uid             string   // Uid
	Gid             string   // Uid
	VolumePath      string   // 数据路径
	Timezone        string   // 时区
	ComposeProfiles string   // 启动项目
	ComposeFiles    []string // compose文件
	NginxVer        string   // Nginx版本
	OpenRestyVer    string   // OpenResty版本
	ApacheVer       string   // Apache版本
	MySQLVer        string   // MySQL版本
	MySQLRootPwd    string   // MySQL root密码
	PGSQLVer        string   // pgsql版本
	PGSQLRootUser   string   // PGSQL root用户
	PGSQLRootPwd    string   // PGSQL root密码
	RedisVer        string   // Redis版本
	MemcachedVer    string   // Memcached版本
	MongoVer        string   // mongo版本
}

type LogConfig struct {
	Level     string `json:"level"`
	LogName   string `json:"logName"`
	LogSuffix string `json:"logSuffix"`
	MaxBackup int    `json:"maxBackup"`
}
