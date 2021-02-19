package token_price

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/tidwall/gjson"

	"github.com/irisnet/service-providers-go/token-price/common"
	"github.com/irisnet/service-providers-go/token-price/types"
)

var TokenPriceAPI = "https://pro-api.coinmarketcap.com/v1/cryptocurrency/quotes/latest?CMC_PRO_API_KEY=43bf626f-772c-4d81-9bd4-db23a26f01d4&symbol="

// RequestCallback provider need to supplement service logic
func RequestCallback(reqID, input string) (
	output *types.ServiceOutput,
	requestResult *types.RequestResult,
) {
	output =&types.ServiceOutput{}
	requestResult = &types.RequestResult{}
	serviceInput, err := parseInput(input)
	if err != nil {
		requestResult.State = types.ClientError
		requestResult.Message = "failed to parse input"
		return nil, requestResult
	}
	common.Logger.Info("Get request: \n", serviceInput)
	// Supplementary service logic...
	pair := serviceInput.Pair
	token := strings.Split(pair, "-")

	response, err := http.Get(TokenPriceAPI + token[0] + "," + token[1])
	defer response.Body.Close()
	if err != nil {
		requestResult.State = types.ServiceError
		requestResult.Message = err.Error()
		return nil, requestResult
	}

	resbody, err := ioutil.ReadAll(response.Body)
	if err !=nil {
		requestResult.State = types.ServiceError
		requestResult.Message = "Failed to read request body..."
		return nil, requestResult
	}

	output.Rate, err = getPrice(string(resbody), token)
	if err != nil {
		requestResult.State = types.ServiceError
		requestResult.Message = "Failed to get token price..."
		return nil, requestResult
	}

	requestResult.State = types.Success
	requestResult.Message = "success"

	return output, requestResult
}

func parseInput(input string) (serviceInput *types.ServiceInput, err error) {
	serviceInput = &types.ServiceInput{}
	input = gjson.Get(input, "body").String()
	if len(input) == 0 {
		return serviceInput, nil
	}
	err = json.Unmarshal([]byte(input), serviceInput)

	return serviceInput, err
}

func getPrice(resbody string, tokens []string) (string, error) {
	token0Price := gjson.Get(resbody, "data." + strings.ToUpper(tokens[0]) + ".quote.USD.price").Float()
	token1Price := gjson.Get(resbody, "data." + strings.ToUpper(tokens[1]) + ".quote.USD.price").Float()
	if token1Price == 0 {
		return "", errors.New("get token price err")
	}

	rate := strconv.FormatFloat(token0Price/token1Price, 'E', -1, 64)
	return rate, nil
}
