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
	"log"
	"os"
	"path/filepath"

	"github.com/manifoldco/promptui"
	toml "github.com/pelletier/go-toml/v2"
	"github.com/spf13/cobra"
	"k8s.io/client-go/util/homedir"
)

// Config contains the relevant fields to create and configure a toml file.
type Config struct {
	// Username defines the username of the wakatime user.
	Username string `toml:"username,multiline,omitempty"`
	// APIKey is the key associated with a wakatime account.
	APIKey string `toml:"apikey,multiline,omitempty"`
	// FilePath describes the path of the toml file.
	FilePath string `toml:"filepath,multiline,omitempty"`
}

// setupCmd represents the setup command
var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Setup the wakatime configration in ~/.config/wakatime-view.toml",
	Run: func(cmd *cobra.Command, args []string) {
		config := Config{}
		config.Prompt(cmd, args)
	},
}

func (config *Config) PromptUserName() (err error) {
	validate := func(input string) error {
		if len(input) < 3 {
			return errors.New("username must have more than 3 characters")
		}
		return nil
	}

	prompt := promptui.Prompt{
		Label:    "Please enter your username",
		Validate: validate,
		Default:  "",
	}

	res, err := prompt.Run()
	if err != nil {
		return fmt.Errorf("Prompt for username failed %v", err)
	}

	config.Username = res

	return nil
}

func (config *Config) PromptAPIKey() (err error) {
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

	res, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return fmt.Errorf("Unable to get APIKey %w", err)
	}

	config.APIKey = res

	return nil
}

func (config *Config) Prompt(cmd *cobra.Command, args []string) {
	// At this point I think we just assume we know where users want to store their config file.
	defaultConfigFile := filepath.Join(homedir.HomeDir(), ".config", "wakatime-view.toml")
	err := config.PromptUserName()
	if err != nil {
		log.Fatalf("Prompting the user failed %e", err)
	}

	err = config.PromptAPIKey()
	if err != nil {
		log.Fatalf("Prompting the user failed %e", err)
	}

	if _, err := os.Stat(defaultConfigFile); err == nil {
		fmt.Printf("%s already exists. Overwrite? (y/n) ", defaultConfigFile)
		prompt := promptui.Prompt{
			Label: "The wakatime-view.toml already exists. Overwrite? (y/n)",
			Validate: func(input string) error {
				if input != "y" && input != "n" {
					return errors.New("Please enter y or n")
				}
				return nil
			},
		}
		result, err := prompt.Run()
		if err != nil {
			fmt.Printf("Unable to create config file %v", err)
			return
		}

		if result == "n" {
			fmt.Println("Exiting...")
			return
		}
	} else {

		t, err := toml.Marshal(config)
		if err != nil {
			fmt.Printf("Unable to marshal into config file %v", err)
			return
		}

		err = writeFile(defaultConfigFile, t)
		if err != nil {
			fmt.Printf("Error writing config file: %v\n", err)
			return
		}
	}
}

func writeFile(filename string, data []byte) error {
	// If the file already exists, we'll just overwrite it.
	f, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}

	defer f.Close()
	f, err = os.OpenFile("notes.txt", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		log.Fatal(err)
	}

	_, err = f.WriteString("[wakatime]\n")
	if err != nil {
		return fmt.Errorf("failed to write to file: %v", err)
	}

	_, err = f.Write(data)
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
}
