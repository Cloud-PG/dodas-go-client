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

func destroyInf() {

	client := &http.Client{
		Timeout: 300 * time.Second,
	}
	fmt.Println("Submitting request to  : ", clientConf.Im.Host)

	req, err := http.NewRequest("DELETE", string(clientConf.Im.Host)+"/"+infID, nil)

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
		fmt.Println("Removed infrastracture ID: ", infID)
	} else {
		fmt.Println("ERROR:\n", string(body))
		return
	}
}

// destroyCmd represents the destroy command
var destroyCmd = &cobra.Command{
	Use:   "destroy <infID1> ... <infIDN>",
	Args:  cobra.MinimumNArgs(1),
	Short: "Destroy infrastructure with this InfID",
	Long: `
`,
	Run: func(cmd *cobra.Command, args []string) {
		for _, inf := range args {
			infID = inf
			destroyInf()
		}
	},
}

func init() {
	rootCmd.AddCommand(destroyCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// destroyCmd.PersistentFlags().String("foo", "", "A help for foo")
	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// destroyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
