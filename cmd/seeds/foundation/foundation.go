/*
Copyright © 2024 Google LLC

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

package foundation

import (
	"fmt"

	"github.com/GoogleCloudPlatform/pastures-poc-toolkit/internal/fabric"
	"github.com/GoogleCloudPlatform/pastures-poc-toolkit/internal/google"
	"github.com/GoogleCloudPlatform/pastures-poc-toolkit/internal/utils"
	"github.com/spf13/cobra"
)

var (
	dryRun  bool
	verbose bool
)

// Set up pointers to support multiple distinct parents
var (
	FoundationPlant = *FoundationCmd
	FoundationBurn  = *FoundationCmd
)

// FoundationCmd represents the foundation command
var FoundationCmd = &cobra.Command{
	Use:   "foundation",
	Short: "Deploy a foundation-only pasture with no blueprints",
	Long: `Creates a foundation landing zone from the FAST framework.
Projects can optionally be deployed as features into the landing zone. An
example of how to use this pasture:
	
	pasture plant foundation`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		// Construct path for the config
		p, err := utils.ConfigPath()
		if err != nil {
			fmt.Println("Unable to set configuration path")
			cobra.CheckErr(err)
		}

		// Check if Google ADC is valid
		if _, err := google.AppDefaultCredentials(); err != nil {
			cobra.CheckErr(err)
		}

		// Get persistent flags from parent
		dryRun, _ = cmd.Flags().GetBool("dry-run")
		verbose, _ = cmd.Flags().GetBool("verbose")

		// Hydrate configuration
		varFile := fabric.LoadVarsFile(p, "")
		varData := fabric.NewFastConfig()

		if err := varData.ReadConfig(varFile.LocalPath); err != nil {
			fmt.Println("Unable to read var file. Try running pasture plow --rehydrate")
			cobra.CheckErr(err)
		}

		varFile.AddConfig(varData)
		varFile.SetBucket(varData.Prefix) // TODO: this can be optimized by splitting deps and stage vars

		// Load foundation stages
		stages := fabric.InitializeStages(p, varData.Prefix, varFile)

		// Do things with the stages
		for _, s := range stages {
			var firstRun bool = false

			// burn not supported for foundation stage
			if cmd.Parent().Name() == "burn" && s.Type == "foundation" {
				fmt.Println("Skipping foundation stage:", s.Name)
				continue
			}

			// dry run bootstrap stage
			if dryRun && s.Name == "0-bootstrap" {
				fmt.Println("Testing if foundation can be applied to GCP organization")

				if err := s.Init(verbose); err != nil {
					fmt.Println("Cannot initialize stage for dry run")
					cobra.CheckErr(err)
				}

				if err := s.Plan(verbose); err != nil {
					fmt.Println("Foundation cannot be applied to GCP organization")
					cobra.CheckErr(err)
				}

				fmt.Println("Foundation can be applied to GCP organization")
				break // Don't do anything else
			}

			// do what we came here to do
			if cmd.Parent().Name() == "burn" {
				fmt.Println("Destroying stage:", s.Name)
			} else {
				fmt.Println("Deploying stage:", s.Name)
			}

			// try fetching dependency files
			if err := s.DiscoverFiles(); err != nil {
				fmt.Println("Pastures first run detected - running with local state")
				firstRun = true
			}

			// check if state needs to be migrated
			fmt.Println("Initializing", s.Name)
			if err := s.Init(verbose); err != nil {
				fmt.Println("Failed to migrate state to remote backend")
				cobra.CheckErr(err)
			}

			fmt.Println("Configuration complete")

			if cmd.Parent().Name() == "burn" {
				// destroy the stage
				fmt.Println("Starting destroy:", s.Name)
				if err := s.Destroy(nil, verbose); err != nil {
					fmt.Println("Stage failed to destroy:", s.Name)
					cobra.CheckErr(err)
				}

				fmt.Println("Successfully destroyed stage:", s.Name)
			} else {
				// apply stage
				fmt.Println("Starting apply:", s.Name)
				if err := s.Apply(nil, verbose); err != nil {
					fmt.Println("Stage failed to deploy:", s.Name)
					cobra.CheckErr(err)
				}

				fmt.Println("Successfully applied stage:", s.Name)

				// move pasture vars to bucket
				if s.Name == "0-bootstrap" {
					fmt.Println("Uploading pasture vars to GCS bucket")

					if err := varFile.UploadFile(); err != nil {
						fmt.Println("Failed to upload pasture var file")
						cobra.CheckErr(err)
					}
				}

				// first run was detected - move things to cloud
				if firstRun {
					// try fetching dependency files
					if err := s.DiscoverFiles(); err != nil {
						fmt.Println("Unable to retrieve stage dependencies for:", s.Name)
						cobra.CheckErr(err)
					}

					// migrate the state
					if err := s.Init(verbose); err != nil {
						fmt.Println("Failed to migrate state to remote backend")
						cobra.CheckErr(err)
					}
				}
			}

			fmt.Println("Stage complete:", s.Name)
		}

		fmt.Println("Navigate to the Google Cloud Console to deploy your first workload:", "https://console.cloud.google.com/welcome")
	},
}

func init() {}
