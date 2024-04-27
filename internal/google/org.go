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
	"errors"
	"fmt"
	"strconv"
	"strings"

	resourcemanager "cloud.google.com/go/resourcemanager/apiv3"
	"cloud.google.com/go/resourcemanager/apiv3/resourcemanagerpb"
	"google.golang.org/api/iterator"
)

func GetOrganization(domain string) (*Organization, error) { // TODO: move this to a method for an org struct
	fmt.Println("\nGetting organization details...")
	ctx := context.Background()

	orgList := []Organization{}

	c, err := resourcemanager.NewOrganizationsClient(ctx)

	if err != nil {
		return nil, err
	}

	defer c.Close()

	req := &resourcemanagerpb.SearchOrganizationsRequest{
		Query: fmt.Sprintf(`domain:%s`, domain),
	}

	it := c.SearchOrganizations(ctx, req)
	for {
		o := Organization{}

		resp, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		o.Domain = domain

		oid := strings.Split(resp.Name, "/")
		orgId, err := strconv.Atoi(oid[1])

		if err != nil {
			return nil, err
		}

		o.Id = orgId
		o.CustomerId = resp.GetDirectoryCustomerId()

		orgList = append(orgList, o)
	}

	// Check if multiple orgs (or no orgs) were found
	if len(orgList) > 1 {
		err := errors.New("too many orgs found - multiple orgs not supported")
		return nil, err
	} else if len(orgList) == 0 {
		err := errors.New("no org found")
		return nil, err
	} else {
		return &orgList[0], err
	}
}
