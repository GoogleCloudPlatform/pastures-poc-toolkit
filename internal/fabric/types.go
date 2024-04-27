/*
Copyright © 2024 Google Inc.

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

package fabric

import (
	"github.com/williamsmt/pastures/internal/google"
	"github.com/williamsmt/pastures/internal/utils"
)

type ConfigValues interface {
	WriteConfig(filePath string) error
	ReadConfig(filePath string) error
}

type FastConfig struct {
	Organization   *google.Organization `json:"organization"`
	BillingAccount *BillingAccount      `json:"billing_account"`
	BootstrapUser  string               `json:"bootstrap_user"`
	FastFeatures   *FastFeatures        `json:"fast_features"`
	Locations      *Locations           `json:"locations"`
	Prefix         string               `json:"prefix"`
	LogSinks       LogSinks             `json:"log_sinks"`
	Groups         *Groups              `json:"groups"`
	Iam            IamAuthoritative     `json:"iam"`
	IamAdditive    IamAdditives         `json:"iam_bindings_additive"`
}

type BillingAccount struct {
	Id           string `json:"id"`
	Is_Org_Level bool   `json:"is_org_level"`
	No_Iam       bool   `json:"no_iam"`
}

type FastFeatures struct {
	Sandbox bool `json:"sandbox"`
}

type Locations struct {
	Bq      string   `json:"bq"`
	Gcs     string   `json:"gcs"`
	Logging string   `json:"logging"`
	PubSub  []string `json:"pubsub"`
}

type Groups struct {
	Gcp_Billing_Admins      string `json:"gcp-billing-admins"`
	Gcp_Devops              string `json:"gcp-devops"`
	Gcp_Network_Admins      string `json:"gcp-network-admins"`
	Gcp_Organization_Admins string `json:"gcp-organization-admins"`
	Gcp_Security_Admins     string `json:"gcp-security-admins"`
	Gcp_Support             string `json:"gcp-support"`
}

type IamAdditive struct {
	Member string `json:"member"`
	Role   string `json:"role"`
}

type IamAuthoritative map[string][]string

type IamAdditives map[string]IamAdditive

type LogSinks map[string]LogSink

type LogSink map[string]string

type Stage struct {
	Name         string
	Type         string
	Repository   *utils.Repo
	Path         string
	ProviderFile *ProviderFile
	StageVars    []*VarsFile
	Factories    []FabricFactory
}

type FabricFactory interface {
	ApplyFactory(prefix string) error
}

type RoleFactory struct {
	Name string
	Path string
}

type ConfigFile interface {
	UploadFile() error
	DownloadFile() error
}

type VarsFile struct {
	Name       string
	LocalPath  string
	RemotePath string
	Bucket     string
	Config     ConfigValues
}

type ProviderFile struct {
	Name       string
	LocalPath  string
	RemotePath string
	Bucket     string
}
