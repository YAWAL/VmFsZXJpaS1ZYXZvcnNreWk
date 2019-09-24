package config

import (
	"reflect"
	"testing"

	"github.com/YAWAL/VmFsZXJpaS1ZYXZvcnNreWk/src/database"
)

func TestReadConfig(t *testing.T) {

	tests := []struct {
		name     string
		file     string
		wantConf Config
		wantErr  bool
	}{
		{
			"Valid config",
			"./testdata/validConfig.json",
			Config{
				Host: ":9999",
				Database: database.Config{
					DataBaseName: "db_name",
					Dialect:      "dialect",
					Password:     "password",
					SSLMode:      "ssl_mode",
					User:         "user",
				},
			},
			false,
		},
		{
			"Valid config",
			"./testdata/invalidConfig.json",
			Config{
				Host: "",
				Database: database.Config{
					DataBaseName: "",
					Dialect:      "",
					Password:     "",
					SSLMode:      "",
					User:         "",
				},
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotConf, err := ReadConfig(tt.file)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotConf, tt.wantConf) {
				t.Errorf("ReadConfig() = %v, want %v", gotConf, tt.wantConf)
				return
			}
		})
	}
}
