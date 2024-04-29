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

package utils

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

func NewRepo() *Repo {
	return &Repo{}
}

func (r *Repo) SetURL(url string) {
	r.Url = url
}

func (r *Repo) SetDestination(dst string) {
	r.Dst = dst
}

func (r *Repo) SetRef(ref string) {
	r.Ref = ref
}

func (r *Repo) SetLink(source string, target string) {
	r.Link.SetSource(source)
	r.Link.SetTarget(target)
}

func (r *Repo) Clone(force bool) error {
	var options git.CloneOptions
	var clone bool = false

	// set options
	options.URL = r.Url

	if r.Ref != "" {
		options.ReferenceName = plumbing.ReferenceName(r.Ref)
	}

	// see if clone already exists
	if err := FileExists(r.Dst); err != nil {
		if force {
			clone = true
		}
	} else {
		clone = true
	}

	if clone {
		// remove existing clone
		if err := RemoveDir(r.Dst); err != nil {
			return err
		}

		// do what we came here to do
		if _, err := git.PlainClone(r.Dst, false, &options); err != nil {
			return err
		}
	}
	return nil
}
