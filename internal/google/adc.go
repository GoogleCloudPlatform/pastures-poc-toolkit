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
	"errors"

	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/GoogleCloudPlatform/pastures-poc-toolkit/internal/utils"
	"golang.org/x/oauth2/google"
)

type UserClaim struct {
	jwt.RegisteredClaims
	ID    int
	Email string
	Name  string
}

const (
	googleKeys  string = "https://www.googleapis.com/oauth2/v3/certs"
	adcJsonName string = "application_default_credentials.json"
)

var (
	userClaim UserClaim
)

func findToken(jsonPath string) (string, error) {
	ctx := context.Background()

	// locate application default credentials
	// read in JSON from path
	jsonData, _ := utils.ReadFile(jsonPath)
	adc, err := google.CredentialsFromJSON(ctx, jsonData)

	if err != nil {
		return "", err // can't find adc
	}

	rawToken, err := adc.TokenSource.Token()

	if err != nil {
		return "", err // found adc, but token is likely expired
	}

	// get google signing keys
	fetchedKeys, err := jwk.Fetch(ctx, googleKeys)

	if err != nil {
		return "", errors.New("unable to fetch signing keys from google")
	}

	// Assuming rawToken is of a type that has an Extra method returning an interface{}
	idTokenInterface := rawToken.Extra("id_token")
	if idTokenInterface == nil {
		// Handle the error: id_token is missing or nil
		return "", errors.New("id_token is missing or nil")
	}

	// Safely assert idTokenInterface to a string
	idToken, ok := idTokenInterface.(string)
	if !ok {
		// Handle the error: idTokenInterface is not a string
		return "", errors.New("idTokenInterface is not a string")
	}

	// parse the token
	jwt.ParseWithClaims(idToken, &userClaim, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return "", errors.New("unexpected token signature method")
		}

		kid, ok := token.Header["kid"].(string)
		if !ok {
			return "", errors.New("could not find key id in token header")
		}

		keys, ok := fetchedKeys.LookupKeyID(kid)
		if !ok {
			return "", errors.New("no keys found matching key id in token header")
		}

		var empty interface{}
		return empty, keys.Raw(&empty)
	})

	return userClaim.Email, nil
}

func AppDefaultCredentials() (string, error) {
	adcPath, err := exec.Command(
		"gcloud", "info", "--format", "value(config.paths.global_config_dir)",
	).Output()

	filePath := strings.Trim(string(adcPath), "\n")
	filePath += "/" + adcJsonName

	if err != nil {
		return "", err
	}

	email, err := findToken(filePath)

	if err != nil {
		fmt.Println("\nNo default credentials found - authorizing with Google")
		authCmd := exec.Command(
			"gcloud", "auth", "application-default", "login",
			"--no-launch-browser",
		)
		authCmd.Stdout = os.Stdout
		authCmd.Stderr = os.Stderr
		authCmd.Stdin = os.Stdin

		if err := authCmd.Run(); err != nil {
			fmt.Printf("Error starting gcloud command: %v\n", err)
			return "", err
		}

		// need to refresh adc context to get email now that we're auth'd
		email, _ = findToken(filePath)
	} else {
		fmt.Println("\nFound default Google credentials - skipping")
	}

	return email, nil
}
