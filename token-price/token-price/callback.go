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
	if len(token) != 2 {
		requestResult.State = types.ClientError
		requestResult.Message = "wrong token input"
		return nil, requestResult
	}

	tokenPriceAPI, err := getAPI()
	if err != nil {
		requestResult.State = types.ServiceError
		requestResult.Message = err.Error()
		return nil, requestResult
	}

	output.Rate, err = getTokenPrice(tokenPriceAPI, token)
	if err !=nil {
		requestResult.State = types.ServiceError
		requestResult.Message = "failes to get token price"
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

func getTokenPrice(tokenPriceAPI string, token []string) (string, error) {
	response, err := http.Get(tokenPriceAPI + token[0] + "," + token[1])
	if err != nil {
		return "", err
	}
	defer func() {
		err := response.Body.Close()
		if err != nil {
			return
		}
	}()

	resbody, err := ioutil.ReadAll(response.Body)
	if err !=nil {
		return "", err
	}

	res := string(resbody)

	token0Price := gjson.Get(res, "data." + strings.ToUpper(token[0]) + ".quote.USD.price").Float()
	token1Price := gjson.Get(res, "data." + strings.ToUpper(token[1]) + ".quote.USD.price").Float()
	if token1Price == 0 {
		return "", errors.New("token1Price = 0")
	}

	rate := strconv.FormatFloat(token0Price/token1Price,'f',precision,64)

	return rate, nil
}
