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
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

// GetVMs ..
func GetVMs(infID string) ([]string, error) {
	///fmt.Println(string(template))
	client := &http.Client{
		Timeout: 300 * time.Second,
	}
	fmt.Println("Submitting request to  : ", clientConf.Im.Host)

	req, err := http.NewRequest("GET", string(clientConf.Im.Host)+"/"+infID, nil)

	req.Header.Set("Content-Type", "application/json")

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

	if resp.StatusCode == 200 {
		fmt.Println("Available Infrastructure VMs:\n", string(body))
	} else {
		return nil, fmt.Errorf(string(body))
	}
	return strings.Split(string(body), "\n"), nil
}

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Args:  cobra.MinimumNArgs(1),
	Short: "Wrapper function for list operations",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
	},
}

// infIDsCmd represents the infIDs command
var infIDsCmd = &cobra.Command{
	Use:   "infIDs",
	Short: "List Infrastructure IDs owned by the current client",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("infIDs called")
		///fmt.Println(string(template))
		client := &http.Client{
			Timeout: 300 * time.Second,
		}
		fmt.Println("Submitting request to  : ", clientConf.Im.Host)

		req, err := http.NewRequest("GET", string(clientConf.Im.Host), nil)

		req.Header.Set("Content-Type", "application/json")

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

		stringList := strings.Split(string(body), "\n")

		fmt.Print("Infrastructure IDs:\n")
		for _, str := range stringList {
			stringSplit := strings.Split(str, "/")

			if resp.StatusCode == 200 {
				fmt.Println(stringSplit[len(stringSplit)-1])
			} else {
				fmt.Println("ERROR:\n", string(body))
				return
			}
		}
	},
}

// infosCmd represents the infos command
var infosCmd = &cobra.Command{
	Use:   "vms-info <infID>",
	Args:  cobra.MinimumNArgs(1),
	Short: "List detailed information for all the machine of an infID",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("infos called")
		listVMs, err := GetVMs(string(args[0]))
		if err != nil {
			panic(err)
		}

		for _, vm := range listVMs {
			client := &http.Client{
				Timeout: 300 * time.Second,
			}
			fmt.Println("Submitting request to  : ", vm)

			req, err := http.NewRequest("GET", vm, nil)

			req.Header.Set("Content-Type", "application/json")

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

			resp, err := client.Do(req)
			if err != nil {
				panic(err)
			}

			body, _ := ioutil.ReadAll(resp.Body)
			fmt.Println(string(body))
		}

	},
}

// vmsCmd represents the vms command
var vmsCmd = &cobra.Command{
	Use:   "vms <infID>",
	Args:  cobra.MinimumNArgs(1),
	Short: "List all the machine of an infID",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("vms called")
		GetVMs(string(args[0]))

	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	listCmd.AddCommand(infIDsCmd)
	listCmd.AddCommand(infosCmd)
	listCmd.AddCommand(vmsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
