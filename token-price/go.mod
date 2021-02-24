module github.com/irisnet/service-providers-go/token-price

go 1.15

replace (
	github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.2-alpha.regen.4
	github.com/tendermint/tendermint => github.com/bianjieai/tendermint v0.34.1-irita-210113
)

require (
	github.com/howeyc/gopass v0.0.0-20190910152052-7cb4b85ec19c
	github.com/irisnet/service-sdk-go e6c49b22a89bcde6b551314c6fc19b
	github.com/sirupsen/logrus v1.6.0
	github.com/tendermint/tendermint master
	github.com/spf13/cobra v1.1.1
	github.com/spf13/viper v1.7.1
	github.com/tidwall/gjson v1.6.1
)
