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

package google

import (
	"context"
	"fmt"
	"strconv"

	"cloud.google.com/go/iam/apiv1/iampb"
	resourcemanager "cloud.google.com/go/resourcemanager/apiv3"
)

func SetRequiredOrgIAMRoles(org *Organization, u string) error {
	ctx := context.Background()
	c, err := resourcemanager.NewOrganizationsClient(ctx)

	if err != nil {
		return err
	}

	defer c.Close()

	roles := []string{
		"roles/billing.admin",
		"roles/logging.admin",
		"roles/iam.organizationRoleAdmin",
		"roles/resourcemanager.projectCreator",
		"roles/resourcemanager.organizationAdmin",
		"roles/resourcemanager.tagAdmin",
		"roles/resourcemanager.folderAdmin",
		"roles/owner",
	}

	// Retrieve the current IAM policy
	getPolicyReq := &iampb.GetIamPolicyRequest{
		Resource: fmt.Sprintf("organizations/%s", strconv.Itoa(org.Id)),
	}
	currentPolicy, err := c.GetIamPolicy(ctx, getPolicyReq)
	if err != nil {
		return err
	}

	// Merge the new roles with the existing ones TODO: rather than merge, is there a 'member' operation that is gracefully additive?
	for _, role := range roles {
		found := false
		for _, binding := range currentPolicy.Bindings {
			if binding.Role == role {
				// If the role already exists, append the new member
				binding.Members = append(binding.Members, fmt.Sprintf("user:%s", u))
				found = true
				break
			}
		}
		if !found {
			// If the role doesn't exist, create a new binding
			currentPolicy.Bindings = append(currentPolicy.Bindings, &iampb.Binding{
				Role:    role,
				Members: []string{fmt.Sprintf("user:%s", u)},
			})
		}
	}

	// Set the updated IAM policy
	setPolicyReq := &iampb.SetIamPolicyRequest{
		Resource: fmt.Sprintf("organizations/%s", strconv.Itoa(org.Id)),
		Policy:   currentPolicy,
	}
	_, err = c.SetIamPolicy(ctx, setPolicyReq)
	if err != nil {
		return err
	}

	return nil
}
