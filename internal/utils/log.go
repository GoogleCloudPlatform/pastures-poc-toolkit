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
	"fmt"
	"os"
)

func NewMesssage() *Message {
	return &Message{}
}

func (m *Message) Print(msg string) {
	// check msg is not empty
	if msg == "" {
		messageNull()
	}
	
	// set the message value
	m.Message = msg

	// print the message without line return
	fmt.Print(m.Message + "\r")
}

func (m *Message) Complete() {
	// check msg is not empty
	if m.Message == "" {
		messageNull()
	}

	// append complete to previous message
	fmt.Println(m.Message, "- Complete!")
}

func (m *Message) Error(e error) {
	// check msg is not empty
	if m.Message == "" {
		messageNull()
	}

	// wrap error in previous message and exit
	fmt.Fprintln(os.Stderr, m.Message, "-", e)
	os.Exit(1)
}

func messageNull() {
	// msg was null - can't perform logging without a message
	fmt.Fprintln(os.Stderr, "Error: Message cannot be null")
	os.Exit(1)
}