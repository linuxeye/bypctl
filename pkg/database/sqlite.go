package database

import (
	"bypctl/pkg/files"
	"bypctl/pkg/global"
	"fmt"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"os"
	"path/filepath"
	"time"
)

func Init() {
	var newLogger logger.Interface
	switch global.Conf.Log.Level {
	case "info":
		newLogger = logger.Default.LogMode(logger.Info)
	case "warn":
		newLogger = logger.Default.LogMode(logger.Warn)
	case "error":
		newLogger = logger.Default.LogMode(logger.Error)
	default:
		newLogger = logger.Default.LogMode(logger.Silent)
	}

	newFile := files.NewFile()
	filePath := filepath.Join(global.Conf.System.BasePath, "db")
	fileName := "bypanel.db"
	if !newFile.Stat(filePath) {
		if err := newFile.CreateDir(filePath, os.ModePerm); err != nil {
			panic(fmt.Errorf("init db dir falied, err: %v", err))
		}
	}

	db, err := gorm.Open(sqlite.Open(filepath.Join(filePath, fileName)), &gorm.Config{
		Logger:                                   newLogger,
		DisableForeignKeyConstraintWhenMigrating: true,
		NamingStrategy:                           schema.NamingStrategy{SingularTable: true},
	})
	if err != nil {
		panic(err)
	}
	_ = db.Exec("PRAGMA journal_mode = WAL;")
	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}
	sqlDB.SetConnMaxIdleTime(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	global.DB = db
	// global.Log.Info("init db successfully")
}
