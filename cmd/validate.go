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

	"github.com/dciangot/toscalib"
	"github.com/spf13/cobra"
)

// Validate TOSCA template
func Validate() {
	fmt.Println("validate called")
	var t toscalib.ServiceTemplateDefinition
	template, err := ioutil.ReadFile(templateFile)
	if err != nil {
		panic(err)
	}

	err = t.Parse(bytes.NewBuffer(template))
	if err != nil {
		panic(err)
	}
	// t.TopologyTemplate.NodeTemplates

	//t.TopologyTemplate.NodeTemplates["Type"]

	//typeList := make(map[string][]string)

	inputs := make(map[string][]string)
	templs := make(map[string][]string)

	for name := range t.TopologyTemplate.NodeTemplates {
		//fmt.Println(name)

		for templ := range t.TopologyTemplate.NodeTemplates[name].Properties {
			if t.TopologyTemplate.NodeTemplates[name].Properties[templ].Value != "" && t.TopologyTemplate.NodeTemplates[name].Properties[templ].Value != nil {
				//fmt.Println(t.TopologyTemplate.NodeTemplates[name].Properties[templ].Value)
				templs[name] = append(templs[name], templ)
			}
			//fmt.Println(templ)
		}

		//fmt.Print("-----\n")
		derived := t.NodeTypes[t.TopologyTemplate.NodeTemplates[name].Type].DerivedFrom
		for derived != "" {
			for interf := range t.NodeTypes[derived].Properties {
				//fmt.Println(interf)
				inputs[name] = append(inputs[name], interf)
			}
			//fmt.Println(derived)
			derived = t.NodeTypes[derived].DerivedFrom
		}

		for interf := range t.NodeTypes[t.TopologyTemplate.NodeTemplates[name].Type].Properties {
			inputs[name] = append(inputs[name], interf)
		}

	}
	//fmt.Println(inputs)
	//fmt.Println(templs)

	for node := range templs {
		//fmt.Println(node)
		for nodeParam := range templs[node] {
			isPresent := false
			for param := range inputs[node] {
				if inputs[node][param] == templs[node][nodeParam] {
					isPresent = true
				}
			}
			//fmt.Printf("%v %v\n", templs[node][nodeParam], isPresent)
			if !isPresent {
				fmt.Printf("%v not defined in type %v \n", templs[node][nodeParam], t.TopologyTemplate.NodeTemplates[node].Type)
				fmt.Printf("ERROR: Invalid template for %v", node)
				return
			}
		}
		//fmt.Print("-----\n")
	}

	fmt.Print("Template OK\n")
}

// validateCmd represents the validate command
var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate your tosca template",
	Long: `Example:
dodas validate --template my_tosca_template.yml`,
	Run: func(cmd *cobra.Command, args []string) {
		Validate()
	},
}

func init() {
	rootCmd.AddCommand(validateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// validateCmd.PersistentFlags().String("foo", "", "A help for foo")
	validateCmd.PersistentFlags().StringVar(&templateFile, "template", "", "Path to TOSCA template file")
	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// validateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
