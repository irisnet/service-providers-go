package common

import (
	"fmt"

	"github.com/spf13/viper"

	"github.com/irisnet/service-providers-go/random-seed/utils"
)

var (
	ConfigPath = utils.MustGetHomeDir() + "/.random-seed-sp/config.yaml"
)

// LoadYAMLConfig loads the YAML config file
func LoadYAMLConfig(configPath string) (*viper.Viper, error) {
	v := viper.New()

	v.SetConfigFile(configPath)
	v.SetConfigType("yaml")

	err := v.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to read the config file: %s", err)
	}

	return v, nil
}

// GetConfigKey returns the key with the given prefix
func GetConfigKey(prefix string, key string) string {
	return fmt.Sprintf("%s.%s", prefix, key)
}
