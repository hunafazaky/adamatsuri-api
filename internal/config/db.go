package config

import (
	"log"

	"github.com/hunafazaky/adamatsuri-api/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DB is kept as a package-level handle for now so existing handlers
// (which call config.DB.___ directly) keep compiling. Phase 2 replaces
// this with a proper repository layer and removes the global entirely.
var DB *gorm.DB

// ConnectDB opens the Postgres connection described by cfg and runs
// auto-migration for all models. It returns the *gorm.DB so callers that
// prefer explicit dependency injection (repositories, tests) can use it
// directly instead of reaching for the global.
func ConnectDB(cfg *Config) *gorm.DB {
	db, err := gorm.Open(postgres.Open(cfg.DSN()), &gorm.Config{TranslateError: true})
	if err != nil {
		log.Fatal("Failed to connect to the database. Err:", err)
	}

	if err := db.AutoMigrate(&model.Event{}, &model.User{}, &model.Booking{}, &model.Tag{}); err != nil {
		log.Fatal("Failed to auto-migrate. Err:", err)
	}

	DB = db
	log.Println("Connected to the database successfully.")
	return db
}
