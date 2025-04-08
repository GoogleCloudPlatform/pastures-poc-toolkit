/*
Copyright © 2025 Google LLC

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
	"github.com/GoogleCloudPlatform/pastures-poc-toolkit/internal/utils"
	"github.com/spf13/cobra"
)

var (
	reqBinaries = map[string]string{
		"gcloud":    "version",
		"terraform": "version",
	}
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initializes environment configuration",
	Long: `Creates the necessary configuration file to establish a
secure organization in Google Cloud.

A valid GCP organization domain name mus tbe supplied as
the argument to this command`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		// Check prereqs
		inform.Print("Running preflight checks")

		if err := checkPreReq(reqBinaries); err != nil {
			// something wrong with a required binary
			inform.Error(err)
		} else {
			// binaries present and behaving as expected
			inform.Complete()
		}
	},
}

func checkPreReq(binaries map[string]string) error {
	// Check the required binaries are available
	for k, v := range binaries {
		if err := utils.CheckInstalled(k, v); err != nil {
			return err
		}
	}

	// binaries are present and functional
	return nil
}

func init() {
	RootCmd.AddCommand(initCmd)
}
