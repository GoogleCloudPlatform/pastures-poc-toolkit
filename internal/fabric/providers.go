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
	providerSuffix  = "-providers.tf"
	providerDirName = "providers"
)

func NewProviderFile(stage string, prefix string, path string) *ProviderFile {
	return &ProviderFile{
		Name:       stage,
		LocalPath:  filepath.Join(path, stage, stage+providerSuffix),
		RemotePath: providerDirName + "/" + stage + providerSuffix,
		Bucket:     bktName(prefix),
	}
}

func (v *ProviderFile) UploadFile() error {
	if err := google.UploadObject(v.Bucket, v.RemotePath, v.LocalPath); err != nil {
		return err
	}

	return nil
}

func (v *ProviderFile) DownloadFile() error {
	if err := google.DownloadObject(v.Bucket, v.LocalPath, v.RemotePath); err != nil {
		return err
	}

	return nil
}
