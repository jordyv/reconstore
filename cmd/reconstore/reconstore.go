package main

import (
	"fmt"
	"os"

	"github.com/integrii/flaggy"
	"github.com/jordyv/reconstore/internal/cmd"
	"github.com/jordyv/reconstore/internal/entities"
	"github.com/jordyv/reconstore/internal/hooks"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	db *gorm.DB

	debug  bool
	dbFile string
)

func init() {
	flaggy.SetName("reconstore")
	flaggy.SetDescription("Reconstore is a tool to save and query your recon data")
	flaggy.SetVersion("1.0.0")
	flaggy.Bool(&debug, "d", "debug", "Debug output")
	flaggy.String(&dbFile, "f", "database", "Database file")
	flaggy.Parse()
}

func main() {
	var err error
	logMode := logger.Error
	if debug {
		logMode = logger.Info
	}
	if dbFile == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			logrus.WithError(err).Fatal("could not determine user home directory")
		}
		dbFile = fmt.Sprintf("%s/reconstore.db", home)
	}
	db, err = gorm.Open(sqlite.Open(dbFile), &gorm.Config{
		Logger: logger.Default.LogMode(logMode),
	})
	if err != nil {
		logrus.Fatal(err.Error())
	}

	migrate()

	cmd.Init(db)
	hooks.Init(db)

	for _, c := range cmd.Commands {
		if c.ShouldHandle() {
			c.Handle()
			break
		}
	}
}

func migrate() {
	db.AutoMigrate(&entities.Program{})
	db.AutoMigrate(&entities.Subdomain{})
	db.AutoMigrate(&entities.Tag{})
}
