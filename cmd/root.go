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
	"os"
	"reflect"
	"strings"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

var version bool
var cfgFile string
var templateFile string
var infID string

// PrepareAuthHeaders ..
func PrepareAuthHeaders() string {

	var authHeaderCloudList []string

	fields := reflect.TypeOf(clientConf.Cloud)
	values := reflect.ValueOf(clientConf.Cloud)

	for i := 0; i < fields.NumField(); i++ {
		field := fields.Field(i)
		value := values.Field(i)

		if value.Interface() != "" {
			keyTemp := fmt.Sprintf("%v = %v", decodeFields[field.Name], value)
			authHeaderCloudList = append(authHeaderCloudList, keyTemp)
		}
	}

	authHeaderCloud := strings.Join(authHeaderCloudList, ";")

	var authHeaderIMList []string

	fields = reflect.TypeOf(clientConf.Im)
	values = reflect.ValueOf(clientConf.Im)

	for i := 0; i < fields.NumField(); i++ {
		field := fields.Field(i)
		if decodeFields[field.Name] != "host" {
			value := values.Field(i)
			if value.Interface() != "" {
				keyTemp := fmt.Sprintf("%v = %v", decodeFields[field.Name], value.Interface())
				authHeaderIMList = append(authHeaderIMList, keyTemp)
			}
		}
	}

	authHeaderIM := strings.Join(authHeaderIMList, ";")

	authHeader := authHeaderCloud + "\\n" + authHeaderIM

	//fmt.Printf(authHeader)

	return authHeader
}

type confCloud struct {
	ID            string `yaml:"id"`
	Type          string `yaml:"type"`
	Username      string `yaml:"username"`
	Password      string `yaml:"password"`
	Host          string `yaml:"host"`
	Tenant        string `yaml:"tenant"`
	AuthUrl       string `yaml:"auth_url"`
	AuthVersion   string `yaml:"auth_version"`
	Domain        string `yaml:"domain"`
	ServiceRegion string `yaml:"service_region"`
}

type confIM struct {
	ID       string `yaml:"id"`
	Type     string `yaml:"type"`
	Host     string `yaml:"host"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Token    string `yaml:"token"`
}

type conf struct {
	Im    confIM    `yaml:"im"`
	Cloud confCloud `yaml:"cloud"`
}

var clientConf conf

func (c *conf) getConf(path string) *conf {

	f, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	err = yaml.UnmarshalStrict(f, c)
	if err != nil {
		panic(err)
	}

	//fmt.Printf("--- c.im:\n%v\n\n", string(c.Im.Host))

	return c
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "dodas",
	Short: "A self-sufficient client for DODAS deployments",
	Long: `A self-sufficient client for DODAS deployments.
Default configuration file searched in $HOME/.dodas.yaml

Usage examples:
"""
# CREATE A CLUSTER FROM TEMPLATE
dodas create --template my_tosca_template.yml

# VALIDATE TOSCA TEMPLATE
dodas validate --template my_tosca_template.yml
"""`,

	Run: func(cmd *cobra.Command, args []string) {
	},
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// VersionString ..
var VersionString string

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(version string) {
	VersionString = version
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// BuildDoc ...
func BuildDoc() {
	err := doc.GenMarkdownTree(rootCmd, "docs")
	if err != nil {
		log.Fatal(err)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "DODAS config file (default is $HOME/.dodas.yaml)")
	rootCmd.Flags().BoolVarP(&version, "version", "v", false, "DODAS client version")
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		//fmt.Println(home)

		// Search config in home directory with name ".dodas_go_client" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".dodas")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		//fmt.Println("Using config file:", viper.ConfigFileUsed())
		clientConf.getConf(viper.ConfigFileUsed())
		//if clientConf.im.Password == "" {
		//	fmt.Println("No password")
		//}
	}
}
