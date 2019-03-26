// Copyright © 2019 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var verbose bool

func sendRequest() {

	fmt.Printf("Template: %v \n", string(templateFile))
	template, err := ioutil.ReadFile(templateFile)
	if err != nil {
		panic(err)
	}

	///fmt.Println(string(template))
	client := &http.Client{
		Timeout: 300 * time.Second,
	}
	fmt.Println("submitting to : ", clientConf.Im.Host)

	req, err := http.NewRequest("POST", string(clientConf.Im.Host), bytes.NewBuffer(template))

	req.Header.Set("Content-Type", "text/yaml")

	// TODO: chain of string and get from config
	req.Header.Set("Authorization", "id = os; type = OpenStack; host = ; username = ; tenant = DODAS; password = ; service_region = recas-cloud;\\nid = im; type = InfrastructureManager; username = ; password = ")

	var request []string
	for name, headers := range req.Header {
		name = strings.ToLower(name)
		for _, h := range headers {
			request = append(request, fmt.Sprintf("%v: %v", name, h))
		}
	}

	request = append(request, fmt.Sprint("\n"))
	fmt.Printf(strings.Join(request, "\n"))

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	//fmt.Println(resp)

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))

}

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		Validate()
		sendRequest()
	},
}

func init() {

	rootCmd.AddCommand(createCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	createCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "verbose local command")
}
