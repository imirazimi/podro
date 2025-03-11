package pkg

import (
	"interview/adapter"
	"strings"
	"sync"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
)

var (
	ConfigVar    *Config
	ConfigOnce   sync.Once
	ConfigPath   = "config.yml"
	ConfigPrefix = "inteview"
)

type Config struct {
	SQLDB      adapter.SQLDBConfig      `koanf:"postgres"`
	HTTPServer adapter.HTTPServerConfig `koanf:"http_server"`
	Logger     LoggerConfig             `koanf:"logger"`
	Migrator   MigratorConfig           `koanf:"migrator"`
}

func GetConfig() (*Config, error) {

	var err error

	ConfigOnce.Do(func() {

		var k = koanf.New(".")

		lErr := k.Load(file.Provider(ConfigPath), yaml.Parser())

		if lErr != nil {
			err = lErr
			return
		}

		lErr = k.Load(env.Provider(ConfigPrefix, ".", func(s string) string {
			str := strings.Replace(strings.ToLower(
				strings.TrimPrefix(s, ConfigPrefix)), "_", ".", -1)

			return strings.Replace(str, "..", "_", -1)
		}), nil)

		if lErr != nil {
			err = lErr
			return
		}

		uErr := k.Unmarshal("", &ConfigVar)

		if uErr != nil {
			err = uErr
			return
		}

	})

	if err != nil {
		return nil, err
	}

	return ConfigVar, nil

}
