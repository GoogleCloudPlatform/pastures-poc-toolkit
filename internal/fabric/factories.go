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
	"os"

	"github.com/williamsmt/pastures/internal/utils"
	"gopkg.in/yaml.v3"
)

func NewRoleFactory(path string) *RoleFactory {
	return &RoleFactory{
		Name: "custom-roles",
		Path: path + "/data/custom-roles",
	}
}

func (f *RoleFactory) ApplyFactory(prefix string) error {
	factoryFiles, err := os.ReadDir(f.Path)

	if err != nil {
		return err
	}

	for _, file := range factoryFiles {
		var r yaml.Node

		// read the file
		p := f.Path + "/" + file.Name()
		inBytes, err := utils.ReadFile(p)

		if err != nil {
			return err
		}

		// load file data into struct
		if err := yaml.Unmarshal(inBytes, &r); err != nil {
			return err
		}

		// update name property
		for i, d := range r.Content {
			// loop thru the nodes
			for x := 0; x < len(d.Content); x += 2 {
				// find the name property
				if d.Content[x].Value == "name" {
					// set the new name in the top level struct
					oldValue := d.Content[x+1].Value
					r.Content[i].Content[x+1].Value = prefix + "_" + oldValue
				}
			}
		}

		// marshal new data to bytes
		outBytes, err := yaml.Marshal(&r)

		if err != nil {
			return err
		}

		// write the updated yaml file
		if err := utils.CreateFile(p, outBytes, true); err != nil {
			return err
		}
	}

	return nil
}
