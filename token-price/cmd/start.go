package main

import (
	"github.com/spf13/cobra"

	"github.com/irisnet/service-providers-go/token-price/app"
	"github.com/irisnet/service-providers-go/token-price/common"
	"github.com/irisnet/service-providers-go/token-price/monitor"
	"github.com/irisnet/service-providers-go/token-price/service"
)

func startCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "start",
		Short:   "Start provider daemon",
		Example: `token-price-sp start [config-file]`,
		Args:    cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			password := getPassword()

			var configPath string

			if len(args) == 0 {
				configPath = common.ConfigPath
			} else {
				common.ConfigPath = args[0]
			}

			config, err := common.LoadYAMLConfig(configPath)
			if err != nil {
				return err
			}

			serviceClient := service.MakeServiceClientWrapper(service.NewConfig(config), password)

			appInstance := app.NewApp(serviceClient)
			appInstance.Start()

			m := monitor.NewMonitor(monitor.NewConfig(config))
			appInstance.StartMonitor(m)

			return nil
		},
	}

	return cmd
}
