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
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

// loginCmd represets the login command

var loginCmd = &cobra.Command{
	Use:   "login <infID> <vmID>",
	Args:  cobra.MinimumNArgs(2),
	Short: "WIP feature: ssh login into a deployed vm",
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

		strBody := string(body)

		asciiBody := strBody
		//strconv.QuoteToASCII(strBody)

		fmt.Printf(asciiBody)
		fmt.Printf("\n\n")

		stringList := strings.Split(asciiBody, "-----END RSA PRIVATE KEY-----")[0]

		fmt.Printf(stringList)
		fmt.Printf("\n\n")

		vmkeyTmp := strings.Split(stringList, "-----BEGIN RSA PRIVATE KEY-----")[1]

		fmt.Printf(vmkeyTmp)
		fmt.Printf("\n\n")

		vmkey := "-----BEGIN RSA PRIVATE KEY-----" + vmkeyTmp + "-----END RSA PRIVATE KEY-----"
		fmt.Printf(vmkey)

		ioutil.WriteFile("/tmp/dat2", []byte(vmkey), 0600)

		// exec.Command("openssl", "rsa", "-in", "/tmp/dat2", "-outform", "PEM", "-out", "/tmp/dat3")
		// if err != nil {
		// 	panic(err)
		// }

		cd := exec.Command("ssh", "-i", "/tmp/dat2", "-lcloudadm", "193.204.89.72", "-o", "UserKnownHostsFile=/dev/null", "-o", "StrictHostKeyChecking=no", "bash")

		grepIn, _ := cd.StdinPipe()
		grepOut, _ := cd.StdoutPipe()
		grepErr, _ := cd.StderrPipe()
		err = cd.Start()
		if err != nil {
			log.Fatal(err)
		}

		go io.Copy(os.Stdout, grepOut)
		go io.Copy(os.Stderr, grepErr)
		go io.Copy(grepIn, os.Stdin)

		log.Printf("Waiting for command to finish...")
		err = cd.Wait()
		log.Printf("Command finished with error: %v", err)

		return

		// rs, err := ioutil.ReadFile("/tmp/dat3")
		// if err != nil {
		// 	return
		// }

		// LOGIN

		//pem.Encode(block)
		// Create the Signer for this private key.
		// signer, err := ssh.ParsePrivateKey(rs) //buffer)
		// if err != nil {
		// 	log.Fatalf("unable to parse private key: %v", err)
		// }

		// config := &ssh.ClientConfig{
		// 	User:            "cloudadm",
		// 	HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		// 	Auth: []ssh.AuthMethod{
		// 		// Use the PublicKeys method for remote authentication.
		// 		ssh.PublicKeys(signer),
		// 	},
		// }

		// client, err := ssh.Dial("tcp", "193.204.89.72:22", config)
		// if err != nil {
		// 	log.Fatal("Failed to dial: ", err)
		// }

		// session, err := client.NewSession()
		// if err != nil {
		// 	log.Fatal("Failed to create session: ", err)
		// }
		// defer session.Close()

		// modes := ssh.TerminalModes{
		// 	ssh.ECHO:          0,     // disable echoing
		// 	ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		// 	ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
		// }

		// if err := session.RequestPty("vt220", 80, 40, modes); err != nil {
		// 	session.Close()
		// 	fmt.Printf("request for pseudo terminal failed: %s", err)
		// }

		// stdin, err := session.StdinPipe()
		// if err != nil {
		// 	fmt.Printf("Unable to setup stdin for session: %v", err)
		// }
		// go io.Copy(stdin, os.Stdin)

		// stdout, err := session.StdoutPipe()
		// if err != nil {
		// 	fmt.Printf("Unable to setup stdout for session: %v", err)
		// }
		// go io.Copy(os.Stdout, stdout)

		// stderr, err := session.StderrPipe()
		// if err != nil {
		// 	fmt.Printf("Unable to setup stderr for session: %v", err)
		// }
		// go io.Copy(os.Stderr, stderr)

		// err = session.Run("bash")
		// if err != nil {
		// 	panic(err)
		// }

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
