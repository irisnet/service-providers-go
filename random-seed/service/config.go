package service

import (
	"github.com/irisnet/service-sdk-go/types"
	"github.com/spf13/viper"

	"github.com/irisnet/service-providers-go/random-seed/common"
	"github.com/irisnet/service-providers-go/random-seed/utils"
)

// default config variables
var (
	defaultChainID       = "iris-hub"
	defaultNodeRPCAddr   = "http://127.0.0.1:26657"
	defaultNodeGRPCAddr  = "127.0.0.1:9090"
	defaultKeyPath       = utils.MustGetHomeDir() + "/.iriscli"
	defaultGas           = uint64(200000)
	defaultFee           = "4point"
	defaultBroadcastMode = types.Commit
	defaultKeyAlgorithm  = "sm2"
)

const (
	Prefix       = "service"
	ChainID      = "chain_id"
	NodeRPCAddr  = "node_rpc_addr"
	NodeGRPCAddr = "node_grpc_addr"
	KeyPath      = "key_path"
	KeyName      = "key_name"
	Fee          = "fee"
	KeyAlgorithm = "key_algorithm"
)

// Config is a config struct for service
type Config struct {
	ChainID      string `yaml:"chain_id"`
	NodeRPCAddr  string `yaml:"node_rpc_addr"`
	NodeGRPCAddr string `yaml:"node_grpc_addr"`
	KeyPath      string `yaml:"key_path"`
	KeyName      string `yaml:"key_name"`
	Fee          string `yaml:"fee"`
	KeyAlgorithm string `yaml:"key_algorithm`
}

// NewConfig constructs a new Config from viper
func NewConfig(v *viper.Viper) Config {
	return Config{
		ChainID:      v.GetString(common.GetConfigKey(Prefix, ChainID)),
		NodeRPCAddr:  v.GetString(common.GetConfigKey(Prefix, NodeRPCAddr)),
		NodeGRPCAddr: v.GetString(common.GetConfigKey(Prefix, NodeGRPCAddr)),
		KeyPath:      v.GetString(common.GetConfigKey(Prefix, KeyPath)),
		KeyName:      v.GetString(common.GetConfigKey(Prefix, KeyName)),
		Fee:          v.GetString(common.GetConfigKey(Prefix, Fee)),
		KeyAlgorithm: v.GetString(common.GetConfigKey(Prefix, KeyAlgorithm)),
	}
}
