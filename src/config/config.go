package config

import (
	"encoding/json"
	"io/ioutil"

	"github.com/YAWAL/VmFsZXJpaS1ZYXZvcnNreWk/src/database"
)

// Config is a structure for general application configuration
type Config struct {
	Host     string          `json:"host"`
	Database database.Config `json:"database"`
}

// ReadConfig perform reading config data from provided json file
func ReadConfig(file string) (conf Config, err error) {
	confData, err := ioutil.ReadFile(file)
	if err != nil {
		return conf, err
	}
	if err = json.Unmarshal(confData, &conf); err != nil {
		return conf, err
	}
	return conf, nil
}
