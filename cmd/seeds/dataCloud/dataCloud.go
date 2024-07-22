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

package dataCloud

import (
	"fmt"
	"strings"

	"github.com/GoogleCloudPlatform/pastures-poc-toolkit/internal/fabric"
	"github.com/GoogleCloudPlatform/pastures-poc-toolkit/internal/google"
	"github.com/GoogleCloudPlatform/pastures-poc-toolkit/internal/terraform"
	"github.com/GoogleCloudPlatform/pastures-poc-toolkit/internal/utils"
	"github.com/spf13/cobra"
)

var (
	dryRun     bool
	skipFast   bool
	region     string
	size       string
	verbose    bool
	isInternal bool
)

// Set up pointers to support multiple distinct parents
var (
	DataCloudPlant = *DataCloudCmd
	DataCloudBurn  = *DataCloudCmd
)

// DataCloudCmd represents the dataCloud command
var DataCloudCmd = &cobra.Command{
	Use:   "data-cloud",
	Short: "Deploy a Data Cloud pasture with blueprints",
	Long: `Creates a data-cloud landing zone in a FAST foundation sandbox.
Blueprints are deployed as features into the landing zone. An
example of how to use this pasture:
	
	pasture plant data-cloud --region us-central1 --pasture-size small`,
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
		skipFast, _ = cmd.Flags().GetBool("skip-foundation")
		verbose, _ = cmd.Flags().GetBool("verbose")
		isInternal, _ = cmd.Flags().GetBool("internal")

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

		// Seed stage
		seed := fabric.NewSeedStage(p)
		seed.HydrateSeed(cmd.Use, varData.Prefix, p)
		seed.AddVarFile(varFile)

		// All at once
		stages = append(stages, seed)

		// Do things with the stages
		for _, s := range stages {
			var firstRun bool = false

			seedVars := make([]*terraform.Vars, 0)

			// burn not supported for foundation stage
			if cmd.Parent().Name() == "burn" && s.Type == "foundation" {
				fmt.Println("Skipping foundation stage:", s.Name)
				continue
			}

			// skip foundation stages
			if skipFast && s.Type == "foundation" {
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

			// Seed variables
			if s.Type == "seed" {
				seedVars = append(seedVars, terraform.AddVar("region", region))
				seedVars = append(seedVars, terraform.AddVar("state_bucket", s.ProviderFile.Bucket))
				seedVars = append(seedVars, terraform.AddVar("state_dir", strings.Split(s.ProviderFile.RemotePath, "/")[0]))
				seedVars = append(seedVars, terraform.AddVar("pasture_size", size))

				if isInternal {
					seedVars = append(seedVars, terraform.AddVar("internal_env", "true"))
				}
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
				if err := s.Destroy(seedVars, verbose); err != nil {
					fmt.Println("Stage failed to destroy:", s.Name)
					cobra.CheckErr(err)
				}

				fmt.Println("Successfully destroyed stage:", s.Name)
			} else {
				// apply stage
				fmt.Println("Starting apply:", s.Name)
				if err := s.Apply(seedVars, verbose); err != nil {
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

			if s.Type == "seed" && cmd.Parent().Name() == "plant" {
				ep, err := terraform.TfOutput(s.Path, "datafusion_endpoint", verbose)

				if err != nil {
					fmt.Println("Stage complete:", s.Name)
				}

				fmt.Println("Navigate to your Data Fusion endpoint to begin data ingestion and integration:", ep)
			}
		}
	},
}

func init() {
	// Define and add flags for the seed
	DataCloudPlant.Flags().StringVarP(&region, "region", "r", "us-central1", "Region for GCP resources to be deployed")
	DataCloudPlant.Flags().StringVarP(&size, "pasture-size", "s", "", "Size of pasture environment - must be 'big' or 'small'")

	// TODO: is there a better way to do this in Cobra?
	DataCloudBurn.Flags().StringVarP(&region, "region", "r", "us-central1", "Region for GCP resources to be deployed")
	DataCloudBurn.Flags().StringVarP(&size, "pasture-size", "s", "", "Size of pasture environment - must be 'big' or 'small'")

	// Required flags
	if err := DataCloudPlant.MarkFlagRequired("region"); err != nil {
		cobra.CheckErr(err)
	}

	if err := DataCloudBurn.MarkFlagRequired("pasture-size"); err != nil {
		cobra.CheckErr(err)
	}
}
