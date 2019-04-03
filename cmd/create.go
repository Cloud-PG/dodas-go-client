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
	"reflect"
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
	fmt.Println("submitting to : ", clientConf.Im.Host)

	req, err := http.NewRequest("POST", string(clientConf.Im.Host), bytes.NewBuffer(template))

	req.Header.Set("Content-Type", "text/yaml")

	var authHeaderCloudList []string

	fields := reflect.TypeOf(clientConf.Cloud)
	values := reflect.ValueOf(clientConf.Cloud)

	for i := 0; i < fields.NumField(); i++ {
		field := fields.Field(i)
		value := values.Field(i)

		keyTemp := fmt.Sprintf("%v = %v", decodeFields[field.Name], value)
		authHeaderCloudList = append(authHeaderCloudList, keyTemp)
	}

	authHeaderCloud := strings.Join(authHeaderCloudList, ";")

	var authHeaderIMList []string

	fields = reflect.TypeOf(clientConf.Im)
	values = reflect.ValueOf(clientConf.Im)

	for i := 0; i < fields.NumField(); i++ {
		field := fields.Field(i)
		value := values.Field(i)
		keyTemp := fmt.Sprintf("%v = %v", decodeFields[field.Name], value.Interface())
		authHeaderIMList = append(authHeaderIMList, keyTemp)
	}

	authHeaderIM := strings.Join(authHeaderIMList, ";")

	authHeader := authHeaderCloud + "\\n" + authHeaderIM

	req.Header.Set("Authorization", authHeader)

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

	stringSplit := strings.Split(string(body), "/")

	fmt.Println("InfrastructureID: ", stringSplit[len(stringSplit)-1])

}

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a cluster from a TOSCA template",
	Long: `
`,
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
	createCmd.PersistentFlags().StringVar(&templateFile, "template", "", "Path to TOSCA template file")
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "verbose local command")
}
