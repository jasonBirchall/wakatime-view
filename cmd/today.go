/*
Copyright © 2022 jason birchall.

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
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
	"github.com/spf13/cobra"
	"k8s.io/client-go/util/homedir"
)

type Client struct {
	url   string
	token string
}

type Response struct {
	Data struct {
		Decimal    string `json:"decimal"`
		Digital    string `json:"digital"`
		IsUpToDate bool   `json:"is_up_to_date"`
	} `json:"data"`
}

func NewClient(url, token string) Client {
	return Client{
		url:   url,
		token: token,
	}
}

// todayCmd represents the today command
var todayCmd = &cobra.Command{
	Use:   "today",
	Short: "Waka today gives you gives you wakatime data for today",
	Run: func(cmd *cobra.Command, args []string) {
		today(cmd, args)
	},
}

func getToken() (string, error) {
	// Get token from user dir config
	defaultConfigFile := filepath.Join(homedir.HomeDir(), ".config", "wakatime-view", "wakatime-view.toml")
	if _, err := os.Stat(defaultConfigFile); os.IsNotExist(err) {
		return "", fmt.Errorf("no wakatime config file found")
	}

	// read defaultConfigFile
	data, err := os.ReadFile(defaultConfigFile)
	if err != nil {
		return "", err
	}

	// parse toml file
	var c Config
	err = toml.Unmarshal(data, &c)
	if err != nil {
		return "", err
	}

	return c.APIKey, nil
}

func today(cmd *cobra.Command, args []string) {
	// Get token from config file or from argument
	token, err := getToken()
	if err != nil {
		log.Fatal(err)
	}

	// Create client
	client := NewClient("https://wakatime.com/api/v1", token)

	// Get data from wakatime

	res, err := client.getWakaData("/users/current/summaries?range=today")
	if err != nil {
		return
	}

	fmt.Println(res)
	// // Print data
	// printWakaData(data)
}

func (c Client) getWakaData(endpoint string) (string, error) {
	res, err := http.Get(c.url + "/" + endpoint + "?api_key=" + c.token)
	if err != nil {
		return "", fmt.Errorf("error getting wakatime data: %w", err)
	}

	defer res.Body.Close()

	// j := Response{}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("error reading wakatime data: %w", err)
	}

	// if err := json.Unmarshal(body, &j); err != nil {
	// 	return "", fmt.Errorf("error unmarshalling wakatime data: %w", err)
	// }

	fmt.Println(string(body))
	return "", fmt.Errorf("wakatime data is not up to date")
}

func init() {
	rootCmd.AddCommand(todayCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// todayCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// todayCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
