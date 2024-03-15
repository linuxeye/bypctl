package log

import (
	"bypctl/pkg/global"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
	"path/filepath"
)

func Init() {
	// Create a new instance of the logger. You can have any number of instances.
	logger := logrus.New()

	// 标准输出和文件都输出
	logger.SetOutput(io.MultiWriter(os.Stdout, &lumberjack.Logger{
		Filename:   filepath.Join(global.Conf.System.LogPath, global.Conf.Log.LogName+global.Conf.Log.LogSuffix),
		MaxBackups: global.Conf.Log.MaxBackup,
	}))

	// 判断系统运行模式
	if global.Conf.System.Mode == logrus.DebugLevel.String() {
		logger.Level = logrus.DebugLevel
	}

	logLevel, err := logrus.ParseLevel(global.Conf.Log.Level)
	if err != nil {
		panic(err)
	}
	logger.SetLevel(logLevel)

	// If you wish to add the calling method as a field
	// logger.SetReportCaller(true)
	logger.SetFormatter(&logrus.JSONFormatter{
		PrettyPrint:     false,                 // 格式化
		TimestampFormat: "2006-01-02 15:04:05", // 时间格式
	})

	global.Log = logger
}
