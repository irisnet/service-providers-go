package random_seed

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/tidwall/gjson"

	"github.com/irisnet/service-providers-go/random-seed/common"
	"github.com/irisnet/service-providers-go/random-seed/types"
)

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
	res, err := getRandomSeed()
	fmt.Println(err)
	if err != nil {
		requestResult.State = types.ServiceError
		requestResult.Message = "failed to get random seed"
		return output, requestResult
	}

	seed := string(res)
	seed = strings.Replace(seed, " ", "", -1)
	seed = strings.Replace(seed, "\n", "", -1)
	output.Seed = seed

	requestResult.State = types.Success
	requestResult.Message = "success"

	return output, requestResult
}

func getRandomSeed() ([]byte, error) {
	RandomSeedAPI,err := getAPI()
	 if err != nil {
	 	return nil, err
	 }

	Response, err := http.Get(RandomSeedAPI)
	if err != nil {
		return nil, err
	}
	defer func() {
		err := Response.Body.Close()
		if err != nil {
			common.Logger.Error(err)
			return
		}
	}()

	res, err := ioutil.ReadAll(Response.Body)
	if err !=nil{
		return nil, err
	}
	return res, nil
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
