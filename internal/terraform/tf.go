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

package terraform // TODO: figure out how to set stdout and stderr in tfexec and log lines

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/hashicorp/terraform-exec/tfexec"
)

func TfInit(dir string, m bool) error {
	//var w string
	var tfInitOptions []tfexec.InitOption

	ctx := context.Background()

	tf, _ := initializeTerraformClient(dir)

	// If true, we need to migrate the state to a remote location
	if m {
		tfInitOptions = append(tfInitOptions, tfexec.ForceCopy(m))
	}

	err := tf.Init(ctx, tfInitOptions...)

	return err
}

func TfPlan(dir string, varFiles []string, vars []Vars) PlanResult {
	var tfPlanOptions []tfexec.PlanOption
	var result PlanResult

	// Find the binary and setup the client
	ctx := context.Background()
	p, err := findBinary()

	if err != nil {
		result.Err = err
		return result
	}

	tf, err := buildClient(dir, p)

	if err != nil {
		result.Err = err
		return result
	}

	// Create the plan file coordinates
	tmpDir, err := os.MkdirTemp(dir, "pastures")
	if err != nil {
		result.Err = err
		return result
	}

	defer os.RemoveAll(tmpDir)

	planPath := fmt.Sprintf("%s/%s-%v", tmpDir, "pastureplan", time.Now().Unix())

	// Set the plan out target
	tfPlanOptions = append(tfPlanOptions, tfexec.Out(planPath))

	// include var files if they're provided
	for _, v := range varFiles {
		tfPlanOptions = append(tfPlanOptions, tfexec.VarFile(v))
	}

	// load up vars if they're supplied
	for _, v := range vars {
		assignment := v.Key + "=" + v.Value
		tfPlanOptions = append(tfPlanOptions, tfexec.Var(assignment))
	}

	// Run the plan
	_, err = tf.Plan(ctx, tfPlanOptions...) // TODO: tf validate before
	if err != nil {
		result.Err = err
		return result
	}

	// Return the plan, or an error
	plan, err := tf.ShowPlanFileRaw(ctx, planPath)
	if err != nil {
		result.Err = err
		return result
	}

	result.Plan = plan
	return result
}

func TfApply(dir string, varFiles []string, vars []*Vars, targets []string) error {
	var tfApplyOptions []tfexec.ApplyOption

	// find the binary and setup the client
	ctx := context.Background()

	tf, _ := initializeTerraformClient(dir)

	// include a var files if they're provided
	for _, v := range varFiles {
		tfApplyOptions = append(tfApplyOptions, tfexec.VarFile(v))
	}

	// load up vars if they're provided
	for _, v := range vars {
		assignment := v.Key + "=" + v.Value
		tfApplyOptions = append(tfApplyOptions, tfexec.Var(assignment))
	}

	// an edge case, but account for specific targets if they're supplied
	for _, t := range targets {
		tfApplyOptions = append(tfApplyOptions, tfexec.Target(t))
	}

	// do what we came here to do
	err := tf.Apply(ctx, tfApplyOptions...) // TODO: tf validate before

	if err != nil {
		return err
	}

	// TODO: return output in case the caller is interested

	return nil
}

func TfDestroy(dir string, varFiles []string, vars []*Vars, targets []string) error {
	var tfDestroyOptions []tfexec.DestroyOption

	// find the binary and setup the client
	ctx := context.Background()

	tf, _ := initializeTerraformClient(dir)

	// include a var files if they're provided
	for _, v := range varFiles {
		tfDestroyOptions = append(tfDestroyOptions, tfexec.VarFile(v))
	}

	// load up vars if they're provided
	for _, v := range vars {
		assignment := v.Key + "=" + v.Value
		tfDestroyOptions = append(tfDestroyOptions, tfexec.Var(assignment))
	}

	// an edge case, but account for specific targets if they're supplied
	for _, t := range targets {
		tfDestroyOptions = append(tfDestroyOptions, tfexec.Target(t))
	}

	// do what we came here to do
	err := tf.Destroy(ctx, tfDestroyOptions...) // TODO: tf validate before

	if err != nil {
		return err
	}

	return nil
}

func TfOutput(dir string, outputVar string) (string, error) {
	var output string

	ctx := context.Background()

	tf, _ := initializeTerraformClient(dir)

	outputs, err := tf.Output(ctx)

	if err != nil {
		return "", err
	}

	if outputVar != "" {
		raw := string(outputs[outputVar].Value)
		output = strings.Trim(raw, `"`)

		// Asked for a value not found in outputs
		if output == "" {
			err := errors.New("output value not found")
			return "", err
		}
	} else {
		bytes, _ := json.Marshal(outputs)
		output = string(bytes)
	}

	return output, nil
}

func TfShow(dir string) (string, error) {
	ctx := context.Background()

	tf, _ := initializeTerraformClient(dir)

	_, err := tf.Show(ctx) // TODO: actually catch state if no error

	if err != nil {
		return "", err
	}

	return "", err // TODO: actually export state
}

func TfPull(dir string) (string, error) {
	ctx := context.Background()

	tf, _ := initializeTerraformClient(dir)

	s, err := tf.StatePull(ctx)

	if err != nil {
		return "", err
	}

	return s, nil
}

func NewVars() *Vars {
	return &Vars{}
}

func AddVar(k string, v string) *Vars {
	return &Vars{
		Key:   k,
		Value: v,
	}
}

func buildClient(d string, p string) (*tfexec.Terraform, error) {
	tf, err := tfexec.NewTerraform(d, p)

	if err != nil {
		return nil, err // error building Terraform client
	}

	return tf, nil
}

func findBinary() (string, error) {
	execPath, err := exec.LookPath("terraform")

	if err != nil {
		return "", err // unable to find the terraform binary in PATH
	}

	return execPath, nil
}

func initializeTerraformClient(dir string) (*tfexec.Terraform, error) {
	p, err := findBinary()
	if err != nil {
		return nil, err
	}
	tf, err := buildClient(dir, p)
	if err != nil {
		return nil, err
	}
	return tf, nil
}
