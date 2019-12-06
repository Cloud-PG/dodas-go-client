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

var decodeFields = map[string]string{
	"ID":            "id",
	"Type":          "type",
	"Username":      "username",
	"Password":      "password",
	"Token":         "token",
	"Host":          "host",
	"Tenant":        "tenant",
	"AuthVersion":   "auth_version",
	"AuthUrl":       "auth_url",
	"Domain":        "domain",
	"ServiceRegion": "service_region",
}

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
	fmt.Println("Submitting request to  : ", clientConf.Im.Host)

	req, err := http.NewRequest("POST", string(clientConf.Im.Host), bytes.NewBuffer(template))

	req.Header.Set("Content-Type", "text/yaml")

	authHeader := PrepareAuthHeaders()

	req.Header.Set("Authorization", authHeader)

	var request []string
	for name, headers := range req.Header {
		name = strings.ToLower(name)
		for _, h := range headers {
			request = append(request, fmt.Sprintf("%v: %v", name, h))
		}
	}

	request = append(request, fmt.Sprint("\n"))
	//fmt.Printf(strings.Join(request, "\n"))

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	body, _ := ioutil.ReadAll(resp.Body)

	stringSplit := strings.Split(string(body), "/")

	if resp.StatusCode == 200 {
		fmt.Println("InfrastructureID: ", stringSplit[len(stringSplit)-1])
	} else {
		fmt.Println("ERROR:\n", string(body))
		return
	}

	// TODO: create .dodas dir and save infID

}

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create <Template1> ... <TemplateN>",
	Args:  cobra.MinimumNArgs(1),
	Short: "Create a cluster from a TOSCA template",
	Long: `
`,
	Run: func(cmd *cobra.Command, args []string) {
		for _, temp := range args {
			templateFile = temp
			err := Validate()
			if err != nil {
				panic(err)
			}
			sendRequest()
		}
	},
}

func init() {

	rootCmd.AddCommand(createCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "verbose local command")
}
