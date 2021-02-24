package main

import (
	"github.com/spf13/cobra"

	"github.com/irisnet/service-providers-go/random-seed/app"
	"github.com/irisnet/service-providers-go/random-seed/common"
	"github.com/irisnet/service-providers-go/random-seed/monitor"
	"github.com/irisnet/service-providers-go/random-seed/service"
)

func startCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "start",
		Short:   "Start provider daemon",
		Example: `random-sp start [config-file]`,
		Args:    cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			password := "1234567890"

			var configPath string

			if len(args) == 0 {
				configPath = common.ConfigPath
			} else {
				configPath = args[0]
			}

			config, err := common.LoadYAMLConfig(configPath)
			if err != nil {
				return err
			}

			serviceClient := service.MakeServiceClientWrapper(service.NewConfig(config), password)

			appInstance := app.NewApp(serviceClient)
			appInstance.Start()

			m := monitor.NewMonitor(config)
			appInstance.StartMonitor(m)

			select {}
		},
	}
	return cmd
}
