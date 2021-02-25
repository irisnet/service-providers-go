package app

import (
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/irisnet/service-providers-go/random-seed/common"
	"github.com/irisnet/service-providers-go/random-seed/monitor"
	callback "github.com/irisnet/service-providers-go/random-seed/random"
	"github.com/irisnet/service-providers-go/random-seed/service"
	"github.com/irisnet/service-providers-go/random-seed/types"
)

// App represents the provider application
type App struct {
	ServiceClient   service.ServiceClientWrapper
	RequestCallback types.RequestCallback
	Logger          *log.Logger
}

// NewApp constructs a new App instance
func NewApp(serviceClient service.ServiceClientWrapper) App {
	return App{
		ServiceClient:   serviceClient,
		RequestCallback: callback.RequestCallback,
		Logger:          common.Logger,
	}
}

// Start starts the provider process
func (app App) Start() {
	// Subscribe
	err := app.ServiceClient.SubscribeServiceRequest(app.RequestCallback)
	if err != nil {
		app.Logger.Errorf("failed to subscribe service request, err: %s", err.Error())
		return
	}
}

func (app App) StartMonitor(monitor *monitor.Monitor) {
	common.Logger.Infof("monitor started, provider addresses: %v", monitor.ProviderAddresses)

	for {
		monitor.Scan()
		if !monitor.Stopped {
			time.Sleep(monitor.Interval)
			continue
		}
		return
	}
}
