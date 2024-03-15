package migrations

import (
	"bypctl/internal/models"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

var AddInitTable = &gormigrate.Migration{
	ID: "20240315-add-init-table",
	Migrate: func(tx *gorm.DB) error {
		if err := tx.AutoMigrate(
			&models.WebsiteDnsAccount{},
			&models.WebsiteAcmeAccount{},
		); err != nil {
			return err
		}
		return nil
	},
}
