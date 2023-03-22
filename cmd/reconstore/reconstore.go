package main

import (
	"github.com/integrii/flaggy"
	"github.com/jordyv/reconstore/internal/cmd"
	"github.com/jordyv/reconstore/internal/config"
	"github.com/jordyv/reconstore/internal/entities"
	"github.com/jordyv/reconstore/internal/hooks"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	db *gorm.DB

	debug                bool
	dbType, dbConnString string
)

func init() {
	flaggy.SetName("reconstore")
	flaggy.SetDescription("Reconstore is a tool to save and query your recon data")
	flaggy.SetVersion("1.0.0")
	flaggy.Bool(&debug, "d", "debug", "Debug output")
	flaggy.Parse()

	config.Initialize()
	dbType = config.GetString(config.DBType)
	dbConnString = config.GetString(config.DBConnectionString)
}

func getDBDialector() gorm.Dialector {
	switch dbType {
	case "sqlite":
		return sqlite.Open(dbConnString)
	case "postgres":
		return postgres.Open(dbConnString)
	default:
		logrus.Fatalf("unsupported database '%s'", dbType)
		return nil
	}
}

func main() {
	var err error
	logMode := logger.Silent
	if debug {
		logMode = logger.Info
	}
	db, err = gorm.Open(getDBDialector(), &gorm.Config{
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
