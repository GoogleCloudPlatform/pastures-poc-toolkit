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

// burnCmd represents the burn command
var burnCmd = &cobra.Command{
	Use:   "burn SEED",
	Short: "Removes the POC resources created by a seed.",
	Long: `Removes the POC resources created by a seed in a previous
run of the plant command. Example:

	pasture burn data-cloud --jumpstart data-warehouse
	
A list of seed templates is shown by running:
	
	pasture burn --help`,
	Args: cobra.ExactArgs(1),
	// Run: func(cmd *cobra.Command, args []string) {
	// 	// always print the help command if invoked without subcommand
	// 	cmd.Help()
	// 	os.Exit(0)
	// },
}

func init() {
	RootCmd.AddCommand(burnCmd)

	// Define persistent flags for all seeds
	burnCmd.PersistentFlags().Bool("skip-foundation", false, "Prevents the Fabric FAST landing zone deployment")
	burnCmd.PersistentFlags().Bool("dry-run", false, "Displays the desired state of the POC")
	burnCmd.PersistentFlags().Bool("local-only", false, "Skip migrating state and vars to remote backend")

	// Hide these flags
	if err := burnCmd.PersistentFlags().MarkHidden("skip-foundation"); err != nil {
		cobra.CheckErr(err)
	}

	if err := burnCmd.PersistentFlags().MarkHidden("local-only"); err != nil {
		cobra.CheckErr(err)
	}

	addSeedToBurn()
}

func addSeedToBurn() {
	burnCmd.AddCommand(&dataCloud.DataCloudBurn)
	burnCmd.AddCommand(&foundation.FoundationBurn)
}
