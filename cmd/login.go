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
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"

	"github.com/spf13/cobra"
)

// getKey extract key from returned string
func getKey(asciiBody string) (string, error) {

	stringList := strings.Split(asciiBody, "-----END RSA PRIVATE KEY-----")[0]

	if len(stringList) < 2 {
		return "", fmt.Errorf("Cannot find private key for this machine")
	}

	vmkeyTmp := strings.Split(stringList, "-----BEGIN RSA PRIVATE KEY-----")[1]

	vmkey := "-----BEGIN RSA PRIVATE KEY-----" + vmkeyTmp + "-----END RSA PRIVATE KEY-----"

	return vmkey, nil
}

// getPubIP from extracted string
func getPubIP(asciiBody string) (string, error) {

	stringList := strings.Split(asciiBody, "net_interface.1.ip = '")

	if len(stringList) < 2 {
		return "", fmt.Errorf("Cannot find a public IP for this machine")
	}

	partialString := stringList[1]

	pubIPTmp := strings.Split(partialString, "' and")[0]

	return pubIPTmp, nil
}


// loginCmd represets the login command
var loginCmd = &cobra.Command{
	Use:   "login <infID> <vmID>",
	Args:  cobra.MinimumNArgs(2),
	Short: "ssh login into a deployed vm",
	Long:  ``,

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("login called")

		// GET KEYS
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

		asciiBody := string(body)

		vmkey, err := getKey(asciiBody)
		if err != nil {
			panic(err)
		}

		pubIP, err := getPubIP(asciiBody)
		if err != nil {
			panic(err)
		}

		ioutil.WriteFile("/tmp/data", []byte(vmkey), 0600)

		rs, err := ioutil.ReadFile("/tmp/data")
		if err != nil {
			return
		}

		// Create the Signer for this private key.
		signer, err := ssh.ParsePrivateKey(rs) //buffer)
		if err != nil {
			log.Fatalf("unable to parse private key: %v", err)
		}

		config := &ssh.ClientConfig{
			User:            "cloudadm",
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
			Auth: []ssh.AuthMethod{
		 		// Use the PublicKeys method for remote authentication.
				ssh.PublicKeys(signer),
			},
		}

		client, err := ssh.Dial("tcp", pubIP+":22", config)
		if err != nil {
			log.Fatal("Failed to dial: ", err)
		}

		session, err := client.NewSession()
		if err != nil {
			log.Fatal("Failed to create session: ", err)
		}
		defer session.Close()

		fd := int(os.Stdin.Fd())
		state, err := terminal.MakeRaw(fd)
		if err != nil {
			fmt.Printf("terminal make raw: %s", err)
		}
		defer terminal.Restore(fd, state)
	
		w, h, err := terminal.GetSize(fd)
		if err != nil {
			fmt.Printf("terminal get size: %s", err)
		}
	
		modes := ssh.TerminalModes{
			ssh.ECHO:          1,
			ssh.TTY_OP_ISPEED: 14400,
			ssh.TTY_OP_OSPEED: 14400,
		}
	
		term := os.Getenv("TERM")
		if term == "" {
			term = "xterm-256color"
		}
		if err := session.RequestPty(term, h, w, modes); err != nil {
			fmt.Printf("session xterm: %s", err)
		}
	
		session.Stdout = os.Stdout
		session.Stderr = os.Stderr
		session.Stdin = os.Stdin
	
		if err := session.Shell(); err != nil {
			fmt.Printf("session shell: %s", err)
		}
	
		if err := session.Wait(); err != nil {
			if e, ok := err.(*ssh.ExitError); ok {
				switch e.ExitStatus() {
				case 130:
					fmt.Printf("ssh: %s", err)
				}
			}
			fmt.Printf("ssh: %s", err)
		}

	},
}

func init() {
	rootCmd.AddCommand(loginCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// loginCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// loginCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
