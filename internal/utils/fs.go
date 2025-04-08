/*
Copyright © 2025 Google LLC

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

package utils

import (
	"io"
	"os"
	"path/filepath"
)

func ConfigPath(d string) (string, error) {
	// set user home directory
	home, err := os.UserHomeDir()

	// user home directory not found
	if err != nil {
		return "", err
	}

	// join config directory name with home directory path
	p := filepath.Join(home, d)

	// return fully qualified path of config directory
	return p, nil
}

func (f *File) Read() error {
	// construct fully qualified path of file
	p := filepath.Join(f.Path, f.Name)

	// attempt to open the file
	file, err := os.Open(p)

	// could not open the file
	if err != nil {
		return err
	}

	// close the file once we're done here
	defer file.Close()

	// capture bytes from file read
	bytes, err := io.ReadAll(file)

	// couldnt read the file
	if err != nil {
		return err
	}

	// set bytes to structure content field
	f.Content = bytes

	// all done here
	return nil
}