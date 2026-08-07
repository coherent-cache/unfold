package db

import (
	"runtime"

	"github.com/rs/zerolog/log"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var Conn *gorm.DB

func Init(dbPath string) {
	var err error
	Conn, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		log.Error().Err(err).Msg("Failed to connect to the database")
		runtime.Goexit()
	}

	if err := Conn.AutoMigrate(&Transactions{}); err != nil {
		log.Error().Err(err).Msg("Failed to migrate the database schema")
		runtime.Goexit()
	}

	log.Debug().Msg("Initiated a new database connection")
}
