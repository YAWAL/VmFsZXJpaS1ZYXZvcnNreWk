package database

import (
	"fmt"

	"github.com/YAWAL/VmFsZXJpaS1ZYXZvcnNreWk/src/logging"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

const configString = "user=%s dbname=%s sslmode=%s password=%s"

// PGconn performs connection to Postgres
func PGconn(conf Config) (db *gorm.DB, err error) {
	db, err = gorm.Open(conf.Dialect, fmt.Sprintf(configString, conf.User,
		conf.DataBaseName, conf.SSLMode, conf.Password))
	if err != nil {
		logging.Log.Errorf("error during connection to Postgres: %s", err.Error())
		return nil, err
	}
	logging.Log.Info("Connection to Postgres has been established")
	return db, err
}
