/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/user"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

// setupCmd represents the setup command
var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Setup the wakatime configration in ~/.wakatime.cfg",
	Run:   prompt,
}

var (
	username   string
	apiKey     string
	configFile string
)

func promptUserName() string {
	validate := func(input string) error {
		if len(input) < 3 {
			return errors.New("Username must have more than 3 characters")
		}
		return nil
	}

	var username string
	u, err := user.Current()
	if err == nil {
		username = u.Username
	}

	prompt := promptui.Prompt{
		Label:    "Username",
		Validate: validate,
		Default:  username,
	}

	result, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return ""
	}

	return result
}

func promptAPIKey() string {
	validate := func(input string) error {
		if len(input) < 6 {
			return errors.New("API Key must have more than 6 characters")
		}
		return nil
	}

	prompt := promptui.Prompt{
		Label:    "API Key",
		Validate: validate,
		Mask:     '*',
	}

	result, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return ""
	}

	return result
}

func prompt(cmd *cobra.Command, args []string) {
	username = promptUserName()
	apiKey = promptAPIKey()

	if configFile == "" {
		configFile = "~/.wakatime-view.cfg"
	}

	// Check ini file exists
	if _, err := os.Stat(configFile); err == nil {
		fmt.Printf("%s already exists. Overwrite? (y/n) ", configFile)
		prompt := promptui.Prompt{
			Label: "Overwrite",
			Validate: func(input string) error {
				if input != "y" && input != "n" {
					return errors.New("Please enter y or n")
				}
				return nil
			},
		}
		result, err := prompt.Run()
		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}
		if result == "n" {
			fmt.Println("Exiting...")
			return
		}
	}

	// Write ini config file
	ini := `[settings]
username = %s
api_key = %s
`
	ini = fmt.Sprintf(ini, username, apiKey)
	err := writeFile(configFile, ini)
	if err != nil {
		fmt.Printf("Error writing config file: %v\n", err)
		return
	}

	fmt.Printf("Username: %s\n", username, apiKey) // nolint: errcheck
}

func writeFile(filename string, data string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(data)
	if err != nil {
		return err
	}

	return nil
}

func init() {
	rootCmd.AddCommand(setupCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// setupCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// setupCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	setupCmd.Flags().StringVarP(&configFile, "config", "c", "", "Config file")
}
