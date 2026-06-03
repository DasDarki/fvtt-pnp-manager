package database

import (
	"github.com/aetherwright/backend/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Connect(dsn string) (*gorm.DB, error) {
	return gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
}

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.User{},
		&models.RefreshToken{},
		&models.ProviderCredential{},
		&models.ApiKey{},
		&models.Campaign{},
		&models.Folder{},
		&models.Character{},
		&models.Item{},
		&models.Scene{},
		&models.ImageEntry{},
		&models.Asset{},
		&models.Memory{},
		&models.Tag{},
		&models.EntityTag{},
		&models.Link{},
		&models.DalleJob{},
		&models.FoundryConnection{},
	)
}
