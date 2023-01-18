package database

import (
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Driver is database driver to use in the app.
//
// Supported database driver: MySQL, Postgres(PostgreSQL)
type Driver string

const (
	MySQL    Driver = "mysql"
	Postgres Driver = "postgres"
)

func New(driver Driver, cfg *Config) *gorm.DB {
	var db *gorm.DB
	var err error
	switch driver {
	case MySQL:
		db, err = mysqlConn(cfg)
		break
	case Postgres:
		db, err = postgresConn(cfg)
		break
	}
	if err != nil {
		log.Fatalf("[database error] failed to connect database. %v", err)
	}
	return db
}

func Migrate(db *gorm.DB, models ...any) error {
	err := db.AutoMigrate(models...)
	if err != nil {
		log.Printf("[database error] error migrate. %v", err)
		return err
	}
	return nil
}

func mysqlConn(cfg *Config) (*gorm.DB, error) {
	dns := cfg.DatabaseUser + ":" + cfg.DatabasePassword + "@tcp(" + cfg.DatabaseHost + ":" + cfg.DatabasePort + ")/" + cfg.DatabaseName + "?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dns), &gorm.Config{})
	if err != nil {
		log.Fatalf("[database error] failed to connect database. %v", err)
		return nil, err
	}
	return db, nil
}

func postgresConn(cfg *Config) (*gorm.DB, error) {
	dsn := "host=" + cfg.DatabaseHost + " user=" + cfg.DatabaseUser + " password=" + cfg.DatabasePassword + " dbname=" + cfg.DatabaseName + " port=" + cfg.DatabasePort + " sslmode=disable TimeZone=Asia/Jakarta"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("[database error] failed to connect to postgres database. %v", err)
		return nil, err
	}
	return db, nil
}
