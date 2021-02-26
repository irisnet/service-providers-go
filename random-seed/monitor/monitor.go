package monitor

import (
	"context"
	"net/http"
	"time"

	servicesdk "github.com/irisnet/service-sdk-go"
	sdktypes "github.com/irisnet/service-sdk-go/types"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/viper"
	abci "github.com/tendermint/tendermint/abci/types"
	tmsdktypes "github.com/tendermint/tendermint/rpc/core/types"

	"github.com/irisnet/service-providers-go/random-seed/common"
	"github.com/irisnet/service-providers-go/random-seed/types"
)

var (
	baseDenom = "uiris"

	balance = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "balance",
			Help: "",
		},
		nil,
	)
	slashed = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "slashed",
			Help: "",
		},
		nil,
	)
	binding = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "binding",
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

func NewMonitor(viper *viper.Viper) *Monitor {
	rpcURL := viper.GetString("service.node_rpc_addr")
	grpcURL := viper.GetString("service.node_grpc_addr")
	prometheusAddr := viper.GetString("monitor.prometheus_addr")
	interval := viper.GetInt64("monitor.interval")
	providerAddrs := viper.GetStringSlice("monitor.provider_addr")
	threshold := viper.GetInt64("balance.threshold")
	baseDenom = viper.GetString("monitor.base_denom")

	rpcEndpoint := NewEndpointFromURL(rpcURL)
	grpcEndpoint := NewEndpointFromURL(grpcURL)

	cfg := sdktypes.ClientConfig{
		NodeURI:  rpcEndpoint.URL,
		GRPCAddr: grpcEndpoint.URL,
	}
	serviceClient := servicesdk.NewServiceClient(cfg)

	addressMap := make(map[string]bool)
	for _, addr := range providerAddrs {
		addressMap[addr] = true
	}

	startListner(prometheusAddr)

	return &Monitor{
		Client:            serviceClient,
		Interval:          time.Duration(interval) * time.Second,
		Threshold:         threshold,
		ProviderAddresses: addressMap,
	}

}

func startListner(addr string) {
	// Register the summary and the histogram with Prometheus's default registry.
	prometheus.MustRegister(balance)
	prometheus.MustRegister(slashed)
	prometheus.MustRegister(binding)

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
			common.Logger.Error("Prometheus HTTP server ListenAndServe err: ", err)
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
			m.checkBalance(addr)
			m.checkServiceBinding(addr)
		}
		m.lastHeight = endHeight
	}

	m.lastHeight = endHeight
}

func (m *Monitor) parseSlashEvents(blockResult *tmsdktypes.ResultBlockResults) {
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
				requestID, _ := getAttributeValue(event, "request_id")
				slashed.WithLabelValues().Add(1)
				common.Logger.Warnf("slashed for request id %s due to invalid response", requestID)
			}
		}
	}
}

func (m *Monitor) parseSlashEventsFromBlock(endBlockEvents []abci.Event) {
	for _, event := range endBlockEvents {
		if m.IsTargetedSlashEvent(event) {
			requestID, _ := getAttributeValue(event, "request_id")
			slashed.WithLabelValues().Add(1)
			common.Logger.Warnf("slashed for request id %s due to response timeouted", requestID)
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

func (m *Monitor) checkBalance(addr string) {
	baseAccount, err := m.Client.QueryAccount(addr)
	if err != nil {
		common.Logger.Errorf("failed to query balance, err: %s", err)
		return
	}
	balance.WithLabelValues().Set(float64(baseAccount.Coins.AmountOf(baseDenom).Uint64()))
}

func (m *Monitor) checkServiceBinding(addr string) {
	queryServiceBindingResponse, err := m.Client.QueryServiceBinding(types.ServiceName, addr)
	if err != nil {
		binding.WithLabelValues().Set(0)
		common.Logger.Errorf("failed to query balance, err: %s", err)
		return
	}
	if !queryServiceBindingResponse.Available {
		binding.WithLabelValues().Set(0)
	} else {
		binding.WithLabelValues().Set(1)
	}
}

func (m *Monitor) Stop() {
	common.Logger.Info("monitor stopped")
	m.Stopped = true
}

func getAttributeValue(event abci.Event, attributeKey string) (string, error) {
	stringEvents := sdktypes.StringifyEvents([]abci.Event{event})
	return stringEvents.GetValue(event.Type, attributeKey)
}
