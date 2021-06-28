package timescaledb

import (
	"fmt"
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/seb7887/janus/internal/config"
)

const (
	dbType = "postgres"
)

var (
	DB  *gorm.DB
	c   = config.GetConfig()
	dsn = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", c.TSHost, c.TSPort, c.TSUser, c.TSPassword, c.TSDatabase)
)

func InitTimescaleDB() {
	database, err := gorm.Open(dbType, dsn)
	if err != nil {
		log.Fatalf("error initializating TimescaleDB %s", err.Error())
	}

	// Set connection pool
	database.DB().SetMaxIdleConns(20)
	database.DB().SetMaxOpenConns(200)
	DB = database
}

func AutoMigrate() error {
	err := DB.AutoMigrate(&Telemetry{}, &Log{}).Error
	if err != nil {
		log.Fatalf("error executing migrations %s", err.Error())
	}
	return nil
}
