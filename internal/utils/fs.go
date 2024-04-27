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

package utils

import (
	"errors"
	"io"
	"os"
	"path/filepath"
)

const (
	configDir = ".pastures"
)

func ConfigPath() (string, error) {
	home, err := os.UserHomeDir()

	if err != nil {
		return "", err
	}

	p := filepath.Join(home, configDir)

	return p, nil
}

func CreateDir(p string) error {
	if _, err := os.Stat(p); os.IsNotExist(err) {
		os.Mkdir(p, 0755)
	} else {
		return errors.New("configuration directory already exists")
	}
	return nil
}

func CreateFile(p string, d []byte, o bool) error {
	if !o {
		if _, err := os.Stat(p); !os.IsNotExist(err) {
			return err
		}
	}

	if err := os.WriteFile(p, d, 0644); err != nil { // TODO: don't assume 0644 perms
		return err
	}

	return nil
}

func ReadFile(filePath string) ([]byte, error) {
	file, err := os.Open(filePath)

	if err != nil {
		return nil, err
	}

	defer file.Close()

	bytes, _ := io.ReadAll(file)

	return bytes, nil
}

func RemoveDir(p string) error {
	if _, err := os.Stat(p); !os.IsNotExist(err) {
		if err := os.RemoveAll(p); err != nil {
			return err
		}
	}
	return nil
}

func CreateSymlink(symlinks ...*Symlink) error {

	for _, s := range symlinks {

		// Check if the symlink already exists
		if _, err := os.Lstat(s.Source); err == nil {
			// If the symlink exists, remove it
			if err := os.Remove(s.Source); err != nil {
				return err
			}
		}

		// Create the symlink
		if err := os.Symlink(s.Target, s.Source); err != nil {
			return err
		}

	}

	return nil
}

func FileExists(p string) error { // TODO: This should be used in all above functions for test if exists
	if _, err := os.Stat(p); !os.IsNotExist(err) {
		return err
	}

	return nil
}

func NewSymlink(src string, tgt string) *Symlink {
	return &Symlink{
		Source: src,
		Target: tgt,
	}
}

func (s *Symlink) SetSource(source string) {
	s.Source = source
}

func (s *Symlink) SetTarget(target string) {
	s.Target = target
}

func (s *Symlink) Link() error {
	// Check if the symlink already exists
	if _, err := os.Lstat(s.Source); err == nil {
		// If the symlink exists, remove it
		if err := os.Remove(s.Source); err != nil {
			return err
		}
	}

	// Create the symlink
	if err := os.Symlink(s.Target, s.Source); err != nil {
		return err
	}

	return nil
}
