package migration

import (
	"bypctl/pkg/global"
	"bypctl/pkg/migration/migrations"
	"github.com/go-gormigrate/gormigrate/v2"
)

func Init() {
	m := gormigrate.New(global.DB, gormigrate.DefaultOptions, []*gormigrate.Migration{
		migrations.AddInitTable,
	})
	if err := m.Migrate(); err != nil {
		global.Log.Error(err)
	}
	// global.Log.Info("Migration run successfully")
}
