package database

import (
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func New(cfg *Config) (*gorm.DB, error) {
	dns := cfg.DatabaseUser + ":" + cfg.DatabasePassword + "@tcp(" + cfg.DatabaseHost + ":" + cfg.DatabasePort + ")/" + cfg.DatabaseName + "?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dns), &gorm.Config{})
	if err != nil {
		log.Fatalf("[database] failed to connect database. %v", err)
		return nil, err
	}
	return db, nil
}

func Migrate(db *gorm.DB, models ...any) error {
	err := db.AutoMigrate(models...)
	if err != nil {
		log.Printf("[database] error migrate. %v", err)
		return err
	}
	return nil
}
