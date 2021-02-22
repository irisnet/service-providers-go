package token_price

import (
	"github.com/irisnet/service-providers-go/token-price/common"
)

const (
	precision = 6
	prefix = "API"
	CMC = "cmc"
)

func getAPI() (string, error) {
	config , err := common.LoadYAMLConfig(common.ConfigPath)
	if err != nil {
		return "", err
	}
	cmc := config.GetString(common.GetConfigKey(prefix, CMC))
	return cmc, nil
}
