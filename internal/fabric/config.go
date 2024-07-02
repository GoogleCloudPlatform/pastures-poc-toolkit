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

package fabric

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/GoogleCloudPlatform/pastures-poc-toolkit/internal/google"
	"github.com/GoogleCloudPlatform/pastures-poc-toolkit/internal/utils"
)

func NewFastConfig() *FastConfig {
	return &FastConfig{}
}

func (f *FastConfig) SetOrg(d string) error {
	org, err := google.GetOrganization(d)

	if err != nil {
		return err
	}

	f.Organization = org

	return nil
}

func (f *FastConfig) SetBilling(a string, x bool) {
	var b BillingAccount

	b.Id = a

	if x {
		b.Is_Org_Level = false
		b.No_Iam = true
	}

	f.BillingAccount = &b
}

func (f *FastConfig) SetUser(e string) {
	f.BootstrapUser = e
}

func (f *FastConfig) SetFeatures() { // TODO: move this default somewhere else
	var a FastFeatures

	a.Sandbox = true
	f.FastFeatures = &a
}

func (f *FastConfig) SetLocations(l string) { // TODO: update to support various field types in struct
	var loc Locations

	// TODO: this should be field targeted with reflectino loop
	loc.Bq = l
	loc.Gcs = l

	// TODO: move these defaults somewhere else
	loc.Logging = "global"
	loc.PubSub = []string{}

	f.Locations = &loc
}

func (f *FastConfig) SetPrefix(p string) error {
	if i := strings.Count(p, ""); i > 10 {
		err := errors.New("too many characters in prefix")
		return err
	}

	f.Prefix = p

	return nil
}

func (f *FastConfig) SetLogSinks(prefix string, sinks LogSinks) {
	updatedLogSinks := make(LogSinks)

	for k, v := range sinks {
		updatedLogSink := make(LogSink)
		for m, n := range v {
			updatedLogSink[m] = n
		}
		updatedLogSinks[prefix+k] = updatedLogSink
	}

	f.LogSinks = updatedLogSinks
}

func (f *FastConfig) SetGroups(g string) { // TODO: update to support various field types in struct
	var grp Groups

	// TODO: check if group exists first

	// TODO: this should be targeted to each field with reflection loop
	grp.Gcp_Billing_Admins = g
	grp.Gcp_Devops = g
	grp.Gcp_Network_Admins = g
	grp.Gcp_Organization_Admins = g
	grp.Gcp_Security_Admins = g
	grp.Gcp_Support = g

	f.Groups = &grp
}

func (f *FastConfig) AddIamBinding(k string, v []string) error { // TODO: the input paramter should be typed to IamPolicy
	m := make(map[string][]string) // Need to initialize

	if _, kExists := f.Iam[k]; kExists {
		err := errors.New("duplicate key detected. roles must be unique")
		return err
	}

	m[k] = v
	f.Iam = m

	return nil
}

func (f *FastConfig) AddIamMember(policies []*IamAdditive) error {
	m := make(map[string]IamAdditive) // Need to initialize

	for _, policy := range policies {
		if _, kExists := f.IamAdditive[policy.Role]; kExists {
			err := errors.New("duplicate key detected. roles must be unique")
			return err
		}

		m[policy.Role] = IamAdditive{
			Role:   policy.Role,
			Member: policy.Member,
		}
	}

	f.IamAdditive = m

	return nil
}

func (f *FastConfig) WriteConfig(filePath string) error {
	j, err := json.MarshalIndent(f, "", "    ")

	if err != nil {
		return err
	}

	if err := utils.CreateFile(filePath, j, true); err != nil {
		return err
	}

	return nil
}

func (f *FastConfig) ReadConfig(filePath string) error {
	bytes, err := utils.ReadFile(filePath)

	if err != nil {
		return err
	}

	if err := json.Unmarshal(bytes, &f); err != nil {
		return err
	}

	return nil
}
