// Copyright (c) 2016 Pani Networks
// All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/romana/core/romana/util"

	cli "github.com/spf13/cobra"
	config "github.com/spf13/viper"
)

// policyCmd represents the policy commands
var policyCmd = &cli.Command{
	Use:   "policy [add|remove|list]",
	Short: "Add, Remove or List a policy.",
	Long: `Add, Remove or List a policy.

For more information, please check http://romana.io
`,
}

func init() {
	policyCmd.AddCommand(policyAddCmd)
	policyCmd.AddCommand(policyRemoveCmd)
	policyCmd.AddCommand(policyListCmd)
}

var policyAddCmd = &cli.Command{
	Use:          "add [tenantName][policyFile]",
	Short:        "Add a new policy.",
	Long:         `Add a new policy.`,
	RunE:         policyAdd,
	SilenceUsage: true,
}

var policyRemoveCmd = &cli.Command{
	Use:          "remove [tenantName][policyName]",
	Short:        "Remove a specific policy.",
	Long:         `Remove a specific policy.`,
	RunE:         policyRemove,
	SilenceUsage: true,
}

var policyListCmd = &cli.Command{
	Use:          "list [tenantName]",
	Short:        "List policy for a specific tenant.",
	Long:         `List policy for a specific tenant.`,
	RunE:         policyList,
	SilenceUsage: true,
}

// policyAdd adds kubernetes policy for a specific tenant
// using the policyFile provided.
func policyAdd(cmd *cli.Command, args []string) error {
	if len(args) != 2 {
		return util.UsageError(cmd, "TENANT and POLICY FILE name should be provided.")
	}

	tenantName := args[0]
	policyFile := args[1]
	// Tenant check once adaptor add supports for it.
	/*
		if !adaptor.TenantExists(tnt) {
			return errors.New("Tenant doesn't exists: " + tnt)
		}
	*/

	f, err := os.Open(policyFile)
	if err != nil {
		return errors.New("Couldn't open Policy file: " + policyFile)
	}
	defer f.Close()

	baseURL := config.GetString("BaseURL")
	kubeURL := baseURL + fmt.Sprintf(":8080/apis/romana.io/demo/v1")
	kubeURL = kubeURL + fmt.Sprintf("/namespaces/%s/networkpolicys", tenantName)

	req, err := http.NewRequest("POST", kubeURL, f)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 201 {
		if config.GetString("Format") == "json" {
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return err
			}
			fmt.Printf(util.JSONIndent(string(body)))
		} else {
			fmt.Printf("Policy (%s) for Tenant (%s) successfully created.\n", policyFile, tenantName)
		}
		return nil
	} else {
		if config.GetString("Format") == "json" {
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return err
			}
			fmt.Printf(util.JSONIndent(string(body)))
			return nil
		} else {
			return fmt.Errorf("Error creating Policy (%s) for Tenant (%s).", policyFile, tenantName)
		}
	}
}

// policyRemove removes kubernetes policy for a specific tenant
// using the policyName provided.
func policyRemove(cmd *cli.Command, args []string) error {
	if len(args) < 2 {
		return util.UsageError(cmd, "TENANT and POLICY name should be provided.")
	}

	tenantName := args[0]
	policyNames := args[1:]
	// Tenant check once adaptor add supports for it.
	/*
		if !adaptor.TenantExists(tnt) {
			return errors.New("Tenant doesn't exists: " + tnt)
		}
	*/

	var errs []error
	baseURL := config.GetString("BaseURL")
	kubeURL := baseURL + fmt.Sprintf(":8080/apis/romana.io/demo/v1")
	for _, policyName := range policyNames {
		k := kubeURL + fmt.Sprintf("/namespaces/%s/networkpolicys/%s/", tenantName, policyName)

		req, err := http.NewRequest("DELETE", k, nil)
		if err != nil {
			return err
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode == 200 {
			if config.GetString("Format") == "json" {
				body, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					return err
				}
				fmt.Printf(util.JSONIndent(string(body)))
			} else {
				fmt.Printf("Policy (%s) for Tenant (%s) successfully deleted.\n", policyName, tenantName)
			}
		} else {
			if config.GetString("Format") == "json" {
				body, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					return err
				}
				fmt.Printf(util.JSONIndent(string(body)))
			} else {
				errs = append(errs, fmt.Errorf("Error deleting Policy (%s) for Tenant (%s).", policyName, tenantName))
			}
		}
	}

	if len(errs) == 0 {
		return nil
	} else {
		lastErr := errs[len(errs)-1]
		errs := errs[:len(errs)-1]
		for _, e := range errs {
			fmt.Printf("%v\n", e)
		}
		return lastErr
	}
}

// policyList lists kubernetes policies for a specific tenant.
func policyList(cmd *cli.Command, args []string) error {
	fmt.Println("Unimplemented: List policies for a specific tenant.")
	return nil
}
