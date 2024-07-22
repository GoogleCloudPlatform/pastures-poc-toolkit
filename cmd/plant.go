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

package cmd

import (
	"github.com/GoogleCloudPlatform/pastures-poc-toolkit/cmd/seeds/dataCloud"
	"github.com/GoogleCloudPlatform/pastures-poc-toolkit/cmd/seeds/foundation"
	"github.com/spf13/cobra"
)

// plantCmd represents the plant command
var plantCmd = &cobra.Command{
	Use:   "plant SEED",
	Short: "Creates a POC environment from a template",
	Long: `Plant creates a POC environment in a FAST foundation sandbox using
a seed template (e.g. data-cloud). Example:
	
	pasture plant data-cloud --region us-central1 --pasture-size small
	
A list of seed templates is shown by running:
	
	pasture plant --help`,
	Args: cobra.NoArgs,
	// Run: func(cmd *cobra.Command, args []string) {
	// 	// always print the help command if invoked without subcommand
	// 	cmd.Help()
	// 	os.Exit(0)
	// },
}

func init() {
	// Add the plant command to the root command
	RootCmd.AddCommand(plantCmd)

	// Define persistent flags for all seeds
	plantCmd.PersistentFlags().Bool("skip-foundation", false, "Prevents the Fabric FAST landing zone deployment")
	plantCmd.PersistentFlags().Bool("dry-run", false, "Displays the desired state of the POC")
	plantCmd.PersistentFlags().Bool("local-only", false, "Skip migrating state and vars to remote backend")

	// Hide these flags
	if err := plantCmd.PersistentFlags().MarkHidden("skip-foundation"); err != nil {
		cobra.CheckErr(err)
	}

	if err := plantCmd.PersistentFlags().MarkHidden("local-only"); err != nil {
		cobra.CheckErr(err)
	}

	// Load seed command palettes
	addSeedToPlant()
}

func addSeedToPlant() {
	plantCmd.AddCommand(&dataCloud.DataCloudPlant)
	plantCmd.AddCommand(&foundation.FoundationPlant)
}
