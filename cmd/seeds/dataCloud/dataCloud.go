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
	DataCloudCreate  = *DataCloudCmd
	DataCloudDestroy = *DataCloudCmd
)

// DataCloudCmd represents the dataCloud command
var DataCloudCmd = &cobra.Command{
	Use:   "data-cloud",
	Short: "Deploy a Data Cloud pasture with blueprints",
	Long: "\nCreates a data-cloud landing zone in a FAST foundation " +
		"sandbox. Blueprints are deployed as features into the landing zone. " +
		"An example of how to use this pasture:\n\n\t" +
		"pasture create data-cloud --region us-central1 --pasture-size small",
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		// Check if Google ADC is valid
		checkGoogleADCValidity()

		// Construct path for the config
		configPath := getConfigPath()

		// Get persistent flags from parent
		setPersistentFlagsFromParent(cmd)

		// Hydrate the configuration
		varFile, varData := hydrateConfig(configPath)

		// Load foundation stages
		stages := fabric.InitializeFoundationStages(
			configPath,
			varData.Prefix,
			varFile,
		)

		// Seed stage
		seed := fabric.NewSeedStage(configPath)
		seed.HydrateSeed(cmd.Use, varData.Prefix, configPath)
		seed.AddVarFile(varFile)

		// Merge foundation and seed stages
		stages = append(stages, seed)

		// Execute the stages
		processStages(cmd, stages, varFile)
	},
}

func hydrateConfig(configPath string) (*fabric.VarsFile, *fabric.FastConfig) {
	varsFile := fabric.LoadVarsFile(configPath, "")
	varData := fabric.NewFastConfig()

	if err := varData.ReadConfig(varsFile.LocalPath); err != nil {
		fmt.Println(
			"Unable to read var file. Try running pasture configure --rehydrate",
		)
		cobra.CheckErr(err)
	}

	varsFile.AddConfig(varData)
	varsFile.SetBucket(
		varData.Prefix,
	) // TODO: this can be optimized by splitting deps and stage vars

	return varsFile, varData
}

func checkGoogleADCValidity() {
	if _, err := google.AppDefaultCredentials(); err != nil {
		cobra.CheckErr(err)
	}
}

func getConfigPath() string {
	configPath, err := utils.ConfigPath()
	if err != nil {
		fmt.Println("Unable to set configuration path")
		cobra.CheckErr(err)
	}

	return configPath
}

func setPersistentFlagsFromParent(cmd *cobra.Command) {
	dryRun, _ = cmd.Flags().GetBool("dry-run")
	skipFast, _ = cmd.Flags().GetBool("skip-foundation")
	verbose, _ = cmd.Flags().GetBool("verbose")
	isInternal, _ = cmd.Flags().GetBool("internal")
}

func processStages(
	cmd *cobra.Command,
	stages []*fabric.Stage,
	varFile *fabric.VarsFile,
) {
	// Loop through the ordered collection of
	// FAST foundation and seed stages
	for _, s := range stages {

		// Get seed template parameters
		seedVars := getSeedVars(s)

		// Skip over foundation stages if destroying a pasture
		// or the `skipFast` flag has been set
		if shouldSkipStage(cmd, s) {
			continue
		}

		// Smoke test if FAST can be deployed to the current org
		if dryRun && s.Name == "0-bootstrap" {
			handleDryRun(s)
			break // exit the loop early
		}

		// Determine if we've run Pastures on this terminal before
		firstRun := handleFirstRun(s)

		// Initialize the stage
		fmt.Println("Initializing", s.Name)
		if err := s.Init(verbose); err != nil {
			fmt.Println("Failed to migrate state to remote backend")
			cobra.CheckErr(err)
		}
		fmt.Println("Configuration complete")

		// Begin stage execution
		if cmd.Parent().Name() == "destroy" {
			fmt.Println("Destroying stage:", s.Name)
			destroyStage(s, seedVars)
		} else {
			fmt.Println("Deploying stage:", s.Name)
			applyStage(s, seedVars, varFile, firstRun)
		}

		fmt.Println("Stage complete:", s.Name)

		if s.Type == "seed" && cmd.Parent().Name() == "create" {
			handleSeedStage(s)
		}
	}
}

func shouldSkipStage(cmd *cobra.Command, s *fabric.Stage) bool {
	if cmd.Parent().Name() == "destroy" && s.Type == "foundation" {
		fmt.Println("Skipping foundation stage:", s.Name)
		return true
	}

	if skipFast && s.Type == "foundation" {
		fmt.Println("Skipping foundation stage:", s.Name)
		return true
	}

	return false
}

func handleDryRun(s *fabric.Stage) {
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
}

func getSeedVars(s *fabric.Stage) []*terraform.Vars {
	seedVars := make([]*terraform.Vars, 0)

	if s.Type == "seed" {
		seedVars = append(seedVars, terraform.AddVar("region", region))
		seedVars = append(
			seedVars,
			terraform.AddVar("state_bucket", s.ProviderFile.Bucket),
		)
		seedVars = append(
			seedVars,
			terraform.AddVar(
				"state_dir",
				strings.Split(s.ProviderFile.RemotePath, "/")[0],
			),
		)
		seedVars = append(seedVars, terraform.AddVar("pasture_size", size))

		if isInternal {
			seedVars = append(
				seedVars,
				terraform.AddVar("internal_env", "true"),
			)
		}
	}

	return seedVars
}

func handleFirstRun(s *fabric.Stage) bool {
	if err := s.DiscoverFiles(); err != nil {
		fmt.Println("Pastures first run detected - running with local state")
		return true
	}
	return false
}

func destroyStage(s *fabric.Stage, seedVars []*terraform.Vars) {
	fmt.Println("Starting destroy:", s.Name)
	if err := s.Destroy(seedVars, verbose); err != nil {
		fmt.Println("Stage failed to destroy:", s.Name)
		cobra.CheckErr(err)
	}
	fmt.Println("Successfully destroyed stage:", s.Name)
}

func applyStage(
	s *fabric.Stage,
	seedVars []*terraform.Vars,
	varFile *fabric.VarsFile,
	firstRun bool,
) {
	fmt.Println("Starting apply:", s.Name)
	if err := s.Apply(seedVars, verbose); err != nil {
		fmt.Println("Stage failed to deploy:", s.Name)
		cobra.CheckErr(err)
	}
	fmt.Println("Successfully applied stage:", s.Name)

	if s.Name == "0-bootstrap" {
		fmt.Println("Uploading pasture vars to GCS bucket")
		if err := varFile.UploadFile(); err != nil {
			fmt.Println("Failed to upload pasture var file")
			cobra.CheckErr(err)
		}
	}

	if firstRun {
		if err := s.DiscoverFiles(); err != nil {
			fmt.Println("Unable to retrieve stage dependencies for:", s.Name)
			cobra.CheckErr(err)
		}

		if err := s.Init(verbose); err != nil {
			fmt.Println("Failed to migrate state to remote backend")
			cobra.CheckErr(err)
		}
	}
}

func handleSeedStage(s *fabric.Stage) {
	ep, err := terraform.TfOutput(s.Path, "datafusion_endpoint", verbose)
	if err != nil {
		fmt.Println("Stage complete:", s.Name)
	}
	fmt.Println(
		"Navigate to your Data Fusion endpoint to begin data "+
			"ingestion and integration:",
		ep,
	)
}

func init() {
	// Define and add flags for the seed
	DataCloudCreate.Flags().
		StringVarP(
			&region, "region", "r", "us-central1",
			"Region for GCP resources to be deployed",
		)
	DataCloudCreate.Flags().
		StringVarP(
			&size, "pasture-size", "s", "",
			"Size of pasture environment - must be 'big' or 'small'",
		)

	// TODO: is there a better way to do this in Cobra?
	DataCloudDestroy.Flags().
		StringVarP(
			&region, "region", "r", "us-central1",
			"Region for GCP resources to be deployed",
		)
	DataCloudDestroy.Flags().
		StringVarP(&size, "pasture-size", "s", "",
			"Size of pasture environment - must be 'big' or 'small'",
		)

	// Required flags
	if err := DataCloudCreate.MarkFlagRequired("region"); err != nil {
		cobra.CheckErr(err)
	}

	if err := DataCloudDestroy.MarkFlagRequired("pasture-size"); err != nil {
		cobra.CheckErr(err)
	}
}
