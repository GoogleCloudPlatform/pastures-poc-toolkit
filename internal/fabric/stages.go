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
	"sync"

	"github.com/GoogleCloudPlatform/pastures-poc-toolkit/internal/terraform"
	"github.com/GoogleCloudPlatform/pastures-poc-toolkit/internal/utils"
)

const ( // TODO: move all to common file for library
	fabricDst          = "fast"
	fastSrc            = "fast/stages"
	foundationDir      = "foundations"
	seedDst            = "pastures"
	seedSrc            = "terraform"
	seedDir            = "seeds"
	fabricRepo         = "https://github.com/GoogleCloudPlatform/cloud-foundation-fabric.git"
	seedRepo           = "https://github.com/GoogleCloudPlatform/pastures-poc-toolkit"
	outputBucketSuffix = "-prod-iac-core-outputs-0"
)

var (
	fastStages = []string{"0-bootstrap", "1-resman"}
	resmanVars = []string{"0-globals", "0-bootstrap"}
)

func InitializeStages(configPath string, prefix string, vars ...*VarsFile) []*Stage {
	stages := make([]*Stage, 0)
	deps := make([]*VarsFile, 0)

	for _, s := range fastStages {

		if s == "1-resman" {
			for _, r := range resmanVars {
				deps = append(deps, resmanDependencies(r, s, prefix, configPath))
			}
		} else {
			deps = vars
		}

		repo := utils.NewRepo()
		repo.SetURL(fabricRepo)
		repo.SetDestination(filepath.Join(configPath, fabricDst))
		repo.SetLink(
			filepath.Join(configPath, foundationDir),
			filepath.Join(configPath, fabricDst, fastSrc),
		)

		stages = append(stages, &Stage{
			Name:         s,
			Type:         "foundation",
			Path:         filepath.Join(configPath, foundationDir, s),
			Repository:   repo,
			ProviderFile: NewProviderFile(s, prefix, filepath.Join(configPath, foundationDir)),
			StageVars:    deps,
		})
	}

	return stages
}

func NewSeedStage(configPath string) *Stage {
	repo := utils.NewRepo()
	repo.SetURL(seedRepo)
	repo.SetDestination(filepath.Join(configPath, seedDst))
	repo.SetLink(
		filepath.Join(configPath, seedDir),
		filepath.Join(configPath, seedDst, seedSrc),
	)

	return &Stage{
		Type:       "seed",
		Repository: repo,
	}
}

func (s *Stage) HydrateSeed(name string, prefix string, configPath string) {
	s.Name = name
	s.Path = filepath.Join(configPath, seedDir, name)
	s.ProviderFile = NewProviderFile(name, prefix, filepath.Join(configPath, seedDir))
}

func (s *Stage) AddVarFile(file *VarsFile) {
	s.StageVars = append(s.StageVars, file)
}

func (s *Stage) SetFactory(factory FabricFactory) {
	s.Factories = append(s.Factories, factory)
}

func (s *Stage) DiscoverFiles() error {
	files := make([]ConfigFile, 0)

	// grab all of the var files
	for _, v := range s.StageVars {
		files = append(files, v)
	}

	// add the provider file
	files = append(files, s.ProviderFile)

	// try to download them all
	for _, f := range files {
		if err := f.DownloadFile(); err != nil {
			return err
		}
	}

	return nil
}

func (s *Stage) Init() error {
	var migrate bool = false

	// test if module initialized
	_, err := terraform.TfPull(s.Path)

	// try to initialize
	if err != nil {
		if err := terraform.TfInit(s.Path, false); err != nil {
			migrate = true
		}
	}

	// try one more time, but migrate the state
	if migrate {
		if err := terraform.TfInit(s.Path, true); err != nil {
			return err
		}
	}

	return nil
}

func (s *Stage) Plan() error {
	var wg sync.WaitGroup
	var files []string

	// make some channels
	done := make(chan bool)
	result := make(chan terraform.PlanResult)

	// extract var files from stage
	for _, f := range s.StageVars {
		files = append(files, f.LocalPath)
	}

	// start an overwatch
	go utils.ProgressTicker(s.Type, &wg, done)

	// do what we came here to do
	go func() {
		planResult := terraform.TfPlan(s.Path, files, nil)
		done <- true
		result <- planResult
	}()

	// catch the plan result
	planResult := <-result

	// wait for stuff to finish
	wg.Wait()

	// turndown the channels
	close(done)
	close(result)

	if planResult.Err != nil {
		return planResult.Err
	}

	return nil
}

func (s *Stage) Apply(vars []*terraform.Vars) error {
	var wg sync.WaitGroup
	var files []string

	// make some channels
	done := make(chan bool)
	err := make(chan error)

	// extract var files from stage
	for _, f := range s.StageVars {
		files = append(files, f.LocalPath)
	}

	// start an overwatch
	go utils.ProgressTicker(s.Name, &wg, done)

	// do what we came here to do
	go func() {
		applyErr := terraform.TfApply(s.Path, files, vars, nil)
		done <- true
		err <- applyErr
	}()

	// catch any errors
	tfError := <-err

	// wait for stuff to finish
	wg.Wait()

	// turndown the channels
	close(done)
	close(err)

	if tfError != nil {
		return tfError
	}

	return nil
}

func (s *Stage) Destroy(vars []*terraform.Vars) error {
	var wg sync.WaitGroup
	var files []string

	// make some channels
	done := make(chan bool)
	err := make(chan error)

	// extract var files from stage
	for _, f := range s.StageVars {
		files = append(files, f.LocalPath)
	}

	// start an overwatch
	go utils.ProgressTicker(s.Name, &wg, done)

	// do what we came here to do
	go func() {
		destroyErr := terraform.TfDestroy(s.Path, files, vars, nil)
		done <- true
		err <- destroyErr
	}()

	// catch any errors
	tfError := <-err

	// wait for stuff to finish
	wg.Wait()

	// turndown the channels
	close(done)
	close(err)

	if tfError != nil {
		return tfError
	}

	return nil
}

func bktName(prefix string) string {
	return prefix + outputBucketSuffix
}
