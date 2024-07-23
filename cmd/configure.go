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
	"fmt"
	"time"

	"github.com/GoogleCloudPlatform/pastures-poc-toolkit/internal/fabric"
	"github.com/GoogleCloudPlatform/pastures-poc-toolkit/internal/google"
	"github.com/GoogleCloudPlatform/pastures-poc-toolkit/internal/utils"
	"github.com/spf13/cobra"
)

const (
	pastureVer = "v1.1.1" // x-release-please-version
)

var (
	// Global variables for command-line flags
	orgDomain        string
	billingAccountId string
	location         string
	isInternal       bool
	fabricVer        string
	prefix           string
	group            string
	orgAdminSa       string
	rehydrate        bool
	seedVer          string
	skipSeed         bool

	// static variables for prerequisites, etc
	reqBinaries = map[string]string{
		"gcloud":    "version",
		"terraform": "version",
	}

	groupIamRoles = []string{
		"roles/billing.admin",
		"roles/logging.admin",
		"roles/iam.organizationRoleAdmin",
		"roles/resourcemanager.projectCreator",
		"roles/resourcemanager.organizationAdmin",
		"roles/resourcemanager.tagAdmin",
		"roles/resourcemanager.folderAdmin",
		"roles/owner",
	}

	// global vars for other things TODO: these defaults likely belong somewhere else
	gIamRoles         = []string{"roles/resourcemanager.organizationAdmin"}
	gIamAdditiveRoles = []string{"roles/orgpolicy.policyAdmin"}

	// patch FAST not making these unique with prefixes
	logSinks = fabric.LogSinks{
		"audit-logs": fabric.LogSink{
			"filter": "logName:\"/logs/cloudaudit.googleapis.com%2Factivity\" OR logName:\"/logs/cloudaudit.googleapis.com%2Fsystem_event\"",
			"type":   "logging",
		},
		"vpc-sc": fabric.LogSink{
			"filter": "protoPayload.metadata.@type=\"type.googleapis.com/google.cloud.audit.VpcServiceControlAuditMetadata\"",
			"type":   "logging",
		},
	}
)

// configureCmd represents the configure command
var configureCmd = &cobra.Command{
	Use:   "configure",
	Short: "Initializes environment configuration",
	Long: "This command will create an environment and define its " +
		"properties in a pasture configuration file, which is " +
		"located by default at $HOME/.pastures/pasture.yaml.",
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {

		// Check if prereqs are in place
		fmt.Println("Running preflight checks")

		for k, v := range reqBinaries {
			if err := utils.CheckInstalled(k, v); err != nil {
				fmt.Println("Preflight checks failed")
				cobra.CheckErr(err)
			}
		}

		// Construct path for all config
		path, err := utils.ConfigPath()

		if err != nil {
			cobra.CheckErr(err)
		}

		fmt.Println("Creating config directory at:", path)
		if err := utils.CreateDir(path); err != nil {
			fmt.Println("Config directory already exists at:", path)
		}

		// Create a new variable file instance
		vars := fabric.LoadVarsFile(path, prefix)

		// TODO: if requested, print the foundations directory path (to be enhanced with harvest feature)

		// Authorize with Google and get current user
		email, err := google.AppDefaultCredentials()

		if err != nil {
			fmt.Println("Unable to authorize with Google")
			cobra.CheckErr(err)
		}

		// Establish a tfvars file from somewhere
		if rehydrate {
			fmt.Println("Sourcing existing configuration from GCS bucket")

			// download the tfvars file
			if err := vars.DownloadFile(); err != nil {
				fmt.Println("Cannot download existing pastures configuration")
				cobra.CheckErr(err)
			}
		} else {
			fmt.Println("Building a new configuration file")

			if err := vars.GetFileMetadata(); err == nil {
				err := fmt.Errorf("existing pasture for prefix %s found - try running configure with --rehydrate flag", prefix)
				cobra.CheckErr(err)
			}

			// Build fastConfig struct
			fastConfig := fabric.NewFastConfig()

			if err := fastConfig.SetOrg(orgDomain); err != nil {
				cobra.CheckErr(err)
			}
			fastConfig.SetBilling(billingAccountId, isInternal)
			fastConfig.SetUser(email)

			// Enable sandbox for seeds
			if skipSeed {
				fastConfig.SetFeatures(false)
			} else {
				fastConfig.SetFeatures(true)
			}

			fastConfig.SetLocations(location)

			if err := fastConfig.SetPrefix(prefix); err != nil {
				fmt.Println("Prefix must be less than 10 characters")
				cobra.CheckErr(err)
			}

			fastConfig.SetGroups(group)

			// Add IAM policies to vars struct
			if isInternal {
				var adds []*fabric.IamAdditive

				for _, r := range gIamRoles { // Authoritative bindings
					err := fastConfig.AddIamBinding(r, []string{"serviceAccount:" + orgAdminSa})

					if err != nil {
						cobra.CheckErr(err)
					}
				}

				for _, r := range gIamAdditiveRoles { // Nonauthoritative bindings
					adds = append(adds, &fabric.IamAdditive{
						Role:   r,
						Member: "serviceAccount:" + orgAdminSa,
					})
				}

				if err := fastConfig.AddIamMember(adds); err != nil {
					fmt.Println("Unable to set IAM additive policy")
					cobra.CheckErr(err)
				}

				// Customize log sinks
				fastConfig.SetLogSinks(prefix, logSinks) // TODO: refactor to simple slice and iterate in method
			}

			fmt.Println("Applying prerequisite roles to group:", group)

			if err := google.SetRequiredOrgIAMRoles(fastConfig.Organization, group, groupIamRoles); err != nil {
				fmt.Println("Unable to apply prerequisite roles to group:", group)
				cobra.CheckErr(err)
			}

			fmt.Println("Waiting for role assignment propagation")
			time.Sleep(10 * time.Second) // TODO: 10 seconds may or may not be enough for propagation

			// Write the tfvars file
			fmt.Println("Writing configuration file to path:", vars.LocalPath)

			vars.AddConfig(fastConfig)

			if err := vars.Config.WriteConfig(vars.LocalPath); err != nil {
				fmt.Println("Unable to write config file to path")
				cobra.CheckErr(err)
			}
		}

		// Init FAST stages
		stages := fabric.InitializeStages(path, prefix, vars)

		// Create seed stage shell and append to foundations
		if !skipSeed {
			stages = append(stages, fabric.NewSeedStage(path))
		}

		for i, s := range stages {
			// clone repositories
			if s.Type == "foundation" {
				if i > 0 { // we only need to deal with foundation once
					continue
				} else {
					fmt.Printf("Using %s tag for Fabric FAST \n", fabricVer)
					s.Repository.SetRef("refs/tags/" + fabricVer)
				}
			} else if s.Type == "seed" {
				fmt.Printf("Using %s tag for the Pasture seed %s \n", seedVer, s.Name) // TODO: we don't have a seed name here; just a shell
				s.Repository.SetRef("refs/tags/" + seedVer)
			}

			fmt.Println("Cloning repository for", s.Type)
			if err := s.Repository.Clone(false); err != nil {
				fmt.Println("Unable to clone repository")
				cobra.CheckErr(err)
			}

			// symlink relevant subdirs
			if err := s.Repository.Link.Link(); err != nil {
				fmt.Println("Unable to link repository target to directory")
				cobra.CheckErr(err)
			}

			// configure stage factories
			if s.Name == "0-bootstrap" {
				fmt.Println("Updating custom role names in custom role factory")
				roleFactory := fabric.NewRoleFactory(s.Path)

				s.SetFactory(roleFactory)
				for _, f := range s.Factories {
					f.ApplyFactory(prefix)
				}
			}

			// download providers file if rehydrating
			// if rehydrate {
			// 	fmt.Println("Sourcing provider file from GCS bucket for stage:", s.Name)
			// 	if err := s.ProviderFile.DownloadFile(); err != nil {
			// 		fmt.Println("Unable to source provider file for stage:", s.Name)
			// 		cobra.CheckErr(err)
			// 	}
			// }

			// initialize tf directory
			// fmt.Println("Initializing FAST stage:", s.Name)
			// if err := terraform.TfInit(s.Path, false); err != nil {
			// 	fmt.Printf("Unable to intialize FAST stage: %s", s.Name)
			// 	cobra.CheckErr(err)
			// }
		}

		fmt.Println("\nPasture configure complete! configuration hydrated...")
	},
}

func init() {
	// Add the configure command to the root command
	RootCmd.AddCommand(configureCmd)

	// Define and add flags for the configure command
	configureCmd.Flags().StringVarP(&orgDomain, "domain", "d", "", "GCP organization domain name")
	configureCmd.Flags().StringVarP(&billingAccountId, "billing-account", "b", "", "GCP billing account ID")
	configureCmd.Flags().StringVarP(&location, "location", "l", "US", "GCP multi-region location code")
	configureCmd.Flags().StringVar(&fabricVer, "fabric-version", "v29.0.0", "Cloud Foundation Fabric FAST version")
	configureCmd.Flags().BoolVarP(&isInternal, "internal", "G", false, "Internal use only")
	configureCmd.Flags().StringVarP(&prefix, "prefix", "p", "", "Prefix for resources with unique names (max 9 characters)")
	configureCmd.Flags().StringVarP(&group, "group-owner", "g", "", "Name of Cloud Identity group that owns the pastures")
	configureCmd.Flags().StringVar(&orgAdminSa, "org-admin-sa", "", "Service account email of the internal environment administrator")
	configureCmd.Flags().BoolVar(&rehydrate, "rehydrate", false, "Restore previous Pastures configuration from saved version in GCS bucket")
	configureCmd.Flags().StringVar(&seedVer, "seed-version", pastureVer, "Version of pasture seed terraform modules to use")
	configureCmd.Flags().BoolVar(&skipSeed, "skip-seed", false, "Limits deployment to FAST foundation only")

	// One of these flags is required
	configureCmd.MarkFlagsOneRequired("domain", "rehydrate")
	configureCmd.MarkFlagsMutuallyExclusive("domain", "rehydrate")

	// New config flag group
	configureCmd.MarkFlagsRequiredTogether("domain", "billing-account", "group-owner")

	// Internal environment flag group
	configureCmd.MarkFlagsRequiredTogether("internal", "org-admin-sa")
	configureCmd.MarkFlagsMutuallyExclusive("rehydrate", "internal")

	// These flags are always required
	if err := configureCmd.MarkFlagRequired("prefix"); err != nil {
		cobra.CheckErr(err)
	}

	// Hide the internal flags
	if err := configureCmd.Flags().MarkHidden("internal"); err != nil {
		cobra.CheckErr(err)
	}

	if err := configureCmd.Flags().MarkHidden("org-admin-sa"); err != nil {
		cobra.CheckErr(err)
	}

	if err := configureCmd.Flags().MarkHidden("skip-seed"); err != nil {
		cobra.CheckErr(err)
	}

	// if err := configureCmd.Flags().MarkHidden("fabric-version"); err != nil {
	// 	cobra.CheckErr(err)
	// }
}
