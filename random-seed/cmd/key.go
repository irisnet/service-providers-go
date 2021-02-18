package main

import (
	"fmt"
	"io/ioutil"

	gp "github.com/howeyc/gopass"
	"github.com/spf13/cobra"

	"github.com/irisnet/service-providers-go/random-seed/common"
	"github.com/irisnet/service-providers-go/random-seed/service"
)

var (
	keysCmd = &cobra.Command{
		Use:   "keys",
		Short: "Key management commands",
	}
)

// keysAddCmd implements the keys add command
func keysAddCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add [name] [config-file]",
		Short: "Generate a new key",
		Args:  cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			password := getPasswordWhenAddKey()

			var configPath string

			if len(args) == 1 {
				configPath = common.ConfigPath
			} else {
				configPath = args[1]
			}

			config, err := common.LoadYAMLConfig(configPath)
			if err != nil {
				return err
			}

			serviceClient := service.MakeServiceClientWrapper(service.NewConfig(config), password)

			addr, mnemonic, err := serviceClient.AddKey(args[0], serviceClient.Password)
			if err != nil {
				return err
			}

			fmt.Printf("key generated successfully: \n\nname: %s\naddress: %s\nmnemonic: %s\n\n", args[0], addr, mnemonic)

			return nil
		},
	}

	return cmd
}

// keysShowCmd implements the keys show command
func keysShowCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show [name] [config-file]",
		Short: "Show the key information by name",
		Args:  cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			password := getPassword()

			var configPath string

			if len(args) == 1 {
				configPath = common.ConfigPath
			} else {
				configPath = args[1]
			}

			config, err := common.LoadYAMLConfig(configPath)
			if err != nil {
				return err
			}

			serviceClient := service.MakeServiceClientWrapper(service.NewConfig(config), password)

			addr, err := serviceClient.ShowKey(args[0], serviceClient.Password)
			if err != nil {
				return err
			}

			fmt.Printf("%s\n", addr)

			return nil
		},
	}

	return cmd
}

// keysImportCmd implements the keys import command
func keysImportCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "import [name] [key-file] [config-file]",
		Short: "Import a key from the private key armor file",
		Args:  cobra.RangeArgs(2, 3),
		RunE: func(cmd *cobra.Command, args []string) error {
			password := getPassword()

			var configPath string

			if len(args) == 2 {
				configPath = common.ConfigPath
			} else {
				configPath = args[2]
			}

			config, err := common.LoadYAMLConfig(configPath)
			if err != nil {
				return err
			}

			keyArmor, err := ioutil.ReadFile(args[1])
			if err != nil {
				return err
			}

			serviceClient := service.MakeServiceClientWrapper(service.NewConfig(config), password)

			addr, err := serviceClient.ImportKey(args[0], serviceClient.Password, string(keyArmor))
			if err != nil {
				return err
			}

			fmt.Printf("key imported successfully: %s\n", addr)

			return nil
		},
	}

	return cmd
}

func init() {
	keysCmd.AddCommand(
		keysAddCmd(),
		keysShowCmd(),
		keysImportCmd(),
	)
}

func getPassword() string {
	// Get user's password
	fmt.Print("Please enter your password: ")
	pwd, err := gp.GetPasswd()
	if err != nil {
		panic(err)
	}
	return string(pwd)
}

func getPasswordWhenAddKey() string {
	fmt.Print("Please enter your password: ")
	password0, err := gp.GetPasswd()
	if err != nil {
		panic(err)
	}
	fmt.Print("Confirm your password: ")
	password1, err := gp.GetPasswd()
	if err != nil {
		panic(err)
	}
	if string(password0) != string(password1) {
		panic("The two passwords do not match.")
	}
	return string(password0)
}
