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
	"path/filepath"

	"github.com/GoogleCloudPlatform/pastures-poc-toolkit/internal/google"
)

const (
	varFileName string = "pasture-fast.tfvars.json"
	varDirName  string = "tfvars"
	varSuffix   string = ".auto.tfvars.json"
)

func LoadVarsFile(configPath string, prefix string) *VarsFile {
	var bkt string

	if prefix != "" {
		bkt = bktName(prefix)
	}

	return &VarsFile{
		Name:       varFileName,
		LocalPath:  filepath.Join(configPath, varFileName),
		RemotePath: varDirName + "/" + varFileName,
		Bucket:     bkt,
	}
}

func (v *VarsFile) SetBucket(prefix string) {
	v.Bucket = bktName(prefix)
}

func (v *VarsFile) AddConfig(config ConfigValues) {
	v.Config = config
}

func (v *VarsFile) UploadFile() error {
	if err := google.UploadObject(v.Bucket, v.RemotePath, v.LocalPath); err != nil {
		return err
	}

	return nil
}

func (v *VarsFile) DownloadFile() error {
	if err := google.DownloadObject(v.Bucket, v.LocalPath, v.RemotePath); err != nil { // TODO
		return err
	}

	return nil
}

func resmanDependencies(name string, stage string, prefix string, configPath string) *VarsFile {
	vars := &VarsFile{
		Name:       name,
		LocalPath:  filepath.Join(configPath, foundationDir, stage, name+varSuffix),
		RemotePath: varDirName + "/" + name + varSuffix,
		Bucket:     bktName(prefix),
	}

	return vars
}
