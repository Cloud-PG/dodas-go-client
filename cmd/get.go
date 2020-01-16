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
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Args:  cobra.MinimumNArgs(1),
	Short: "Wrapper command for get operations",
	Long: `Wrapper command for get operations.
dodas get -h for possible commands`,
	Run: func(cmd *cobra.Command, args []string) {
	},
}

// outputCmd represents the output command
var outputCmd = &cobra.Command{
	Use:   "output <infID>",
	Args:  cobra.MinimumNArgs(1),
	Short: "Get deployment output",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("status called")

		client := &http.Client{
			Timeout: 300 * time.Second,
		}
		fmt.Println("Submitting request to  : ", clientConf.Im.Host)

		req, err := http.NewRequest("GET", string(clientConf.Im.Host)+"/"+string(args[0])+"/outputs", nil)

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

		fmt.Print("Deployment output:\n")

		if resp.StatusCode == 200 {
			fmt.Println(string(body))
		} else {
			fmt.Println("ERROR:\n", string(body))
			return
		}

	},
}

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status <infID>",
	Args:  cobra.MinimumNArgs(1),
	Short: "Get deployment status",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("status called")

		client := &http.Client{
			Timeout: 300 * time.Second,
		}
		fmt.Println("Submitting request to  : ", clientConf.Im.Host)

		req, err := http.NewRequest("GET", string(clientConf.Im.Host)+"/"+string(args[0])+"/contmsg", nil)

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

		fmt.Print("Deployment status:\n")

		if resp.StatusCode == 200 {
			fmt.Println(string(body))
		} else {
			fmt.Println("ERROR:\n", string(body))
			return
		}

	},
}

// vmstatusCmd represents the status command
var vmstatusCmd = &cobra.Command{
	Use:   "vm <infID> <vmID>",
	Args:  cobra.MinimumNArgs(2),
	Short: "Get VM deployment logs",
	Long: `
`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("vmstatus called")
		listVMs, err := GetVMs(string(args[0]))
		if err != nil {
			panic(err)
		}

		vmN := string(args[1])
		vmID, err := strconv.Atoi(vmN)
		if err != nil {
			panic(err)
		}
		vm := listVMs[vmID]

		clientHTTP := &http.Client{
			Timeout: 300 * time.Second,
		}
		fmt.Println("Submitting request to  : ", vm)

		req, err := http.NewRequest("GET", vm+"/contmsg", nil)
		if err != nil {
			panic(err)
		}

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

		resp, err := clientHTTP.Do(req)
		if err != nil {
			panic(err)
		}

		body, _ := ioutil.ReadAll(resp.Body)

		fmt.Printf("Deployment status for vm %v:\n", vmID)

		if resp.StatusCode == 200 {
			fmt.Println(string(body))
		} else {
			fmt.Println("ERROR:\n", string(body))
			return
		}

	},
}

// vmCmd represents the vm command
var vmCmd = &cobra.Command{
	Use:   "vm <infID> <vmID>",
	Args:  cobra.MinimumNArgs(2),
	Short: "Get VM details",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("vm called")
		listVMs, err := GetVMs(string(args[0]))
		if err != nil {
			panic(err)
		}

		vmN := string(args[1])
		vmID, err := strconv.Atoi(vmN)
		if err != nil {
			panic(err)
		}
		vm := listVMs[vmID]
		clientHTTP := &http.Client{
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

		resp, err := clientHTTP.Do(req)
		if err != nil {
			panic(err)
		}

		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println(string(body))
	},
}

func init() {
	rootCmd.AddCommand(getCmd)

	getCmd.AddCommand(outputCmd)
	getCmd.AddCommand(statusCmd)
	statusCmd.AddCommand(vmstatusCmd)

	getCmd.AddCommand(vmCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
