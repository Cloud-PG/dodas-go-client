/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

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
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

// rebootCmd represents the reboot command
var rebootCmd = &cobra.Command{
	Use:   "reboot <infID> <vmID>",
	Args:  cobra.MinimumNArgs(2),
	Short: "reboot a vm in the cluster",
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("status called")

		client := &http.Client{
			Timeout: 300 * time.Second,
		}
		fmt.Println("Submitting request to  : ", clientConf.Im.Host)

		req, err := http.NewRequest("PUT", string(clientConf.Im.Host)+"/"+string(args[0])+"/vms/"+string(args[1])+"/reboot", nil)

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
			fmt.Println("Command received correctly. VM reboot will start soon.")
		} else {
			fmt.Println("ERROR:\n", string(body))
			return
		}

	},
}

func init() {
	rootCmd.AddCommand(rebootCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// rebootCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// rebootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
