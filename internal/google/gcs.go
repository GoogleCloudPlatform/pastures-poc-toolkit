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

package google

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"cloud.google.com/go/storage"
)

func DownloadObject(bucketName string, savePath string, object string) error {
	ctx := context.Background()

	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to create storage client: %w", err)
	}

	bkt := client.Bucket(bucketName)
	obj := bkt.Object(object)

	r, err := obj.NewReader(ctx)
	if err != nil {
		return fmt.Errorf("failed to create reader for object %s: %w", object, err)
	}
	defer r.Close()

	// fileNameParts := strings.Split(object, "/")
	// fileName := fileNameParts[len(fileNameParts)-1]

	// filePath := filepath.Join(savePath, fileName)

	if _, err := os.Stat(savePath); err == nil {
		if err := os.Remove(savePath); err != nil {
			return fmt.Errorf("failed to remove existing file %s: %w", savePath, err)
		}
	}

	createdFile, err := os.Create(savePath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", savePath, err)
	}
	defer createdFile.Close()

	if _, err := io.Copy(createdFile, r); err != nil {
		return fmt.Errorf("failed to copy object %s to file %s: %w", object, savePath, err)
	}

	return nil
}

func UploadObject(bucketName string, objectPath string, localPath string) error {
	ctx := context.Background()

	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to create storage client: %w", err)
	}
	defer client.Close()

	// Open local file
	f, err := os.Open(localPath)
	if err != nil {
		return fmt.Errorf("os.Open: %w", err)
	}
	defer f.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	o := client.Bucket(bucketName).Object(objectPath)

	wc := o.NewWriter(ctx)
	if _, err = io.Copy(wc, f); err != nil {
		return fmt.Errorf("io.Copy: %w", err)
	}
	if err := wc.Close(); err != nil {
		return fmt.Errorf("Writer.Close: %w", err)
	}

	return nil
}
