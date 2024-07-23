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

// destroyCmd represents the destroy command
var destroyCmd = &cobra.Command{
	Use:   "destroy SEED",
	Short: "Removes the POC resources created by a seed.",
	Long: `Removes the POC resources created by a seed in a previous
run of the plant command. Example:

	pasture destroy data-cloud --jumpstart data-warehouse
	
A list of seed templates is shown by running:
	
	pasture destroy --help`,
	Args: cobra.ExactArgs(1),
	// Run: func(cmd *cobra.Command, args []string) {
	// 	// always print the help command if invoked without subcommand
	// 	cmd.Help()
	// 	os.Exit(0)
	// },
}

func init() {
	RootCmd.AddCommand(destroyCmd)

	// Define persistent flags for all seeds
	destroyCmd.PersistentFlags().Bool("skip-foundation", false, "Prevents the Fabric FAST landing zone deployment")
	destroyCmd.PersistentFlags().Bool("dry-run", false, "Displays the desired state of the POC")
	destroyCmd.PersistentFlags().Bool("local-only", false, "Skip migrating state and vars to remote backend")
	destroyCmd.PersistentFlags().BoolP("internal", "G", false, "Internal use only")

	// Hide these flags
	if err := destroyCmd.PersistentFlags().MarkHidden("skip-foundation"); err != nil {
		cobra.CheckErr(err)
	}

	if err := destroyCmd.PersistentFlags().MarkHidden("local-only"); err != nil {
		cobra.CheckErr(err)
	}

	if err := destroyCmd.PersistentFlags().MarkHidden("internal"); err != nil {
		cobra.CheckErr(err)
	}

	// Load seed command palettes
	addSeedToDestroy()
}

func addSeedToDestroy() {
	destroyCmd.AddCommand(&dataCloud.DataCloudDestroy)
	destroyCmd.AddCommand(&foundation.FoundationDestroy)
}
