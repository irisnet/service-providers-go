package random_seed

import (
	"github.com/irisnet/service-providers-go/random-seed/common"
)

const (
	prefix = "API"
	RANDOMSEEDAPI = "random-seed-api"
)

func getAPI() (string, error) {
	config , err := common.LoadYAMLConfig(common.ConfigPath)
	if err != nil {
		return "", err
	}
	cmc := config.GetString(common.GetConfigKey(prefix, RANDOMSEEDAPI))
	return cmc, nil
}
