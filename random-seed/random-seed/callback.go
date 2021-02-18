package random_seed

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/tidwall/gjson"

	"github.com/irisnet/service-providers-go/random-seed/common"
	"github.com/irisnet/service-providers-go/random-seed/types"
)

const RandomSeedAPI = "https://www.random.org/cgi-bin/randbyte?nbytes=32&format=h"

// RequestCallback provider need to supplementary service logic
func RequestCallback(reqID, input string) (
	output *types.ServiceOutput,
	requestResult *types.RequestResult,
) {
	output = &types.ServiceOutput{}
	requestResult = &types.RequestResult{}
	serviceInput, err := parseInput(input)
	if err != nil {
		requestResult.State = types.ClientError
		requestResult.Message = "failed to parse input"
		return output, requestResult
	}
	common.Logger.Info("Get request: \n", serviceInput)
	// Supplementary service logic...
	Response, err := http.Get(RandomSeedAPI)
	defer Response.Body.Close()
	if err != nil {
		requestResult.State = types.ServiceError
		requestResult.Message = "Request for site failed..."
		return nil, requestResult
	}

	res, err := ioutil.ReadAll(Response.Body)
	if err !=nil{
		requestResult.State = types.ServiceError
		requestResult.Message = "Failed to read request body..."
		return nil, requestResult
	}

	seed := string(res)
	seed = strings.Replace(seed, " ", "", -1)
	seed = strings.Replace(seed, "\n", "", -1)
	output.Seed = seed

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
