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

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create SEED",
	Short: "Creates a POC environment from a template",
	Long: `Create instantiates a POC environment in a FAST foundation sandbox using
a seed template (e.g. data-cloud). Example:
	
	pasture create data-cloud --region us-central1 --pasture-size small
	
A list of seed templates is shown by running:
	
	pasture create --help`,
	Args: cobra.NoArgs,
	// Run: func(cmd *cobra.Command, args []string) {
	// 	// always print the help command if invoked without subcommand
	// 	cmd.Help()
	// 	os.Exit(0)
	// },
}

func init() {
	// Add the create command to the root command
	RootCmd.AddCommand(createCmd)

	// Define persistent flags for all seeds
	createCmd.PersistentFlags().Bool("skip-foundation", false, "Prevents the Fabric FAST landing zone deployment")
	createCmd.PersistentFlags().Bool("dry-run", false, "Displays the desired state of the POC")
	createCmd.PersistentFlags().Bool("local-only", false, "Skip migrating state and vars to remote backend")
	createCmd.PersistentFlags().BoolP("internal", "G", false, "Internal use only")

	// Hide these flags
	if err := createCmd.PersistentFlags().MarkHidden("skip-foundation"); err != nil {
		cobra.CheckErr(err)
	}

	if err := createCmd.PersistentFlags().MarkHidden("local-only"); err != nil {
		cobra.CheckErr(err)
	}

	if err := createCmd.PersistentFlags().MarkHidden("internal"); err != nil {
		cobra.CheckErr(err)
	}

	// Load seed command palettes
	addSeedToCreate()
}

func addSeedToCreate() {
	createCmd.AddCommand(&dataCloud.DataCloudCreate)
	createCmd.AddCommand(&foundation.FoundationCreate)
}
