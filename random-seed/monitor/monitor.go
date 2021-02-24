package monitor

import (
	"context"
	"net/http"
	"time"

	"github.com/spf13/viper"

	servicesdk "github.com/irisnet/service-sdk-go"
	"github.com/irisnet/service-sdk-go/types"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	abci "github.com/tendermint/tendermint/abci/types"
	tmtypes "github.com/tendermint/tendermint/rpc/core/types"

	"github.com/irisnet/service-providers-go/random-seed/common"
)

var (
	balance = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "balance",
			Help: "",
		},
		nil,
	)
	slashed = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "slashed",
			Help: "",
		},
		nil,
	)
)

const (
	ServiceSlashingEventType = "service_slash"
)

type Monitor struct {
	Client            servicesdk.ServiceClient
	RPCEndpoint       Endpoint
	GRPCEndpoint      Endpoint
	Interval          time.Duration
	Threshold         int64
	ProviderAddresses map[string]bool
	lastHeight        int64
	Stopped           bool
}

func NewMonitor(
	rpcEndpoint Endpoint,
	grpcEndpoint Endpoint,
	prometheusAddr string,
	interval time.Duration,
	threshold int64,
	providerAddresses []string,
) *Monitor {
	cfg := types.ClientConfig{
		NodeURI:  rpcEndpoint.URL,
		GRPCAddr: grpcEndpoint.URL,
	}
	serviceClient := servicesdk.NewServiceClient(cfg)

	addressMap := make(map[string]bool)
	for _, addr := range providerAddresses {
		addressMap[addr] = true
	}

	startListner(prometheusAddr)

	return &Monitor{
		Client:            serviceClient,
		RPCEndpoint:       rpcEndpoint,
		GRPCEndpoint:      grpcEndpoint,
		Interval:          interval,
		Threshold:         threshold,
		ProviderAddresses: addressMap,
	}

}

func NewConfig(viper *viper.Viper) (Endpoint, Endpoint, string, time.Duration, int64, []string) {
	rpcURL := viper.GetString("monitors.rpc_endpoint")
	gRPCURL := viper.GetString("monitors.grpc_endpoint")
	prometheusAddr := viper.GetString("monitors.prometheus_addr")
	interval := viper.GetInt64("monitors.interval")
	providerAddrs := viper.GetStringSlice("monitors.provider_addr")
	threshold := viper.GetInt64("monitors.threshold")

	rpcEndpoint := NewEndpointFromURL(rpcURL)
	grpcEndpoint := NewEndpointFromURL(gRPCURL)

	return rpcEndpoint, grpcEndpoint, prometheusAddr, time.Duration(interval) * time.Second, threshold, providerAddrs
}

func startListner(addr string) {
	// Register the summary and the histogram with Prometheus's default registry.
	prometheus.MustRegister(balance)
	prometheus.MustRegister(slashed)
	srv := &http.Server{
		Addr: addr,
		Handler: promhttp.InstrumentMetricHandler(
			prometheus.DefaultRegisterer, promhttp.HandlerFor(
				prometheus.DefaultGatherer,
				promhttp.HandlerOpts{MaxRequestsInFlight: 10},
			),
		),
	}

	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			// Error starting or closing listener:
			panic(err.Error())
		}
	}()
}

func (m *Monitor) Scan() {
	currentHeight, err := m.getLatestHeight()
	if err != nil {
		common.Logger.Warnf("failed to retrieve the latest block height: %s", err)
		return
	}

	common.Logger.Infof("block height: %d", currentHeight)

	if m.lastHeight == 0 {
		m.lastHeight = currentHeight - 1
	}

	m.scanByRange(m.lastHeight+1, currentHeight)
}

func (m Monitor) getLatestHeight() (int64, error) {
	res, err := m.Client.Status(context.Background())
	if err != nil {
		return -1, err
	}

	return res.SyncInfo.LatestBlockHeight, nil
}

func (m *Monitor) scanByRange(startHeight int64, endHeight int64) {
	for h := startHeight; h <= endHeight; h++ {
		_, err := m.Client.BlockResults(context.Background(), &h)
		if err != nil {
			common.Logger.Warnf("failed to retrieve the block result, height: %d, err: %s", h, err)
			continue
		}
	}

	for addr := range m.ProviderAddresses {
		for h := startHeight; h <= endHeight; h++ {
			blockResult, err := m.Client.BlockResults(context.Background(), &h)
			if err != nil {
				common.Logger.Warnf("failed to retrieve the block result, height: %d, err: %s", h, err)
				continue
			}
			m.parseSlashEvents(blockResult)
		}
		m.lastHeight = endHeight

		baseAccount, err := m.Client.QueryAccount(addr)
		if err != nil {
			common.Logger.Errorf("failed to query balance, err: %s", err)
		}
		amount := baseAccount.Coins.AmountOf("uiris")
		balance.WithLabelValues().Set(float64(amount.Quo(types.NewIntWithDecimal(10, 6)).Int64()))
	}

	m.lastHeight = endHeight
}

func (m *Monitor) parseSlashEvents(blockResult *tmtypes.ResultBlockResults) {
	if len(blockResult.TxsResults) > 0 {
		m.parseSlashEventsFromTxs(blockResult.TxsResults)
	}

	if len(blockResult.EndBlockEvents) > 0 {
		m.parseSlashEventsFromBlock(blockResult.EndBlockEvents)
	}
}

func (m *Monitor) parseSlashEventsFromTxs(txsResults []*abci.ResponseDeliverTx) {
	for _, txResult := range txsResults {
		for _, event := range txResult.Events {
			if m.IsTargetedSlashEvent(event) {
				slashed.WithLabelValues().Set(1)
			}
		}
	}
}

func (m *Monitor) parseSlashEventsFromBlock(endBlockEvents []abci.Event) {
	for _, event := range endBlockEvents {
		if m.IsTargetedSlashEvent(event) {
			slashed.WithLabelValues().Set(1)
		}
	}
}

func (m *Monitor) IsTargetedSlashEvent(event abci.Event) bool {
	if event.Type != ServiceSlashingEventType {
		return false
	}

	providerAddr, err := getAttributeValue(event, "provider")
	if err != nil {
		return false
	}

	if _, ok := m.ProviderAddresses[providerAddr]; !ok {
		return false
	}

	return true
}

func (m *Monitor) Stop() {
	common.Logger.Info("monitor stopped")
	m.Stopped = true
}

func getAttributeValue(event abci.Event, attributeKey string) (string, error) {
	stringEvents := types.StringifyEvents([]abci.Event{event})
	return stringEvents.GetValue(event.Type, attributeKey)
}
