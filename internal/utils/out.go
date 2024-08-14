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
	"fmt"
	"sync"
	"time"
)

func ProgressTicker(headline string, wg *sync.WaitGroup, ch <-chan bool) {
	wg.Add(1)
	defer wg.Done()

	ticker := time.NewTicker(
		5 * time.Second,
	) // TODO: this falsely lengthens then wait block by 5 seconds minimum
	defer ticker.Stop()

	startTime := time.Now()

	for {
		select {
		case <-ch:
			return
		case <-ticker.C:
			elapsed := time.Since(startTime)
			minutes := int(elapsed.Minutes())
			seconds := int(elapsed.Seconds()) - minutes*60

			if minutes == 0 {
				fmt.Printf(
					"Still working on %s for %d seconds\n",
					headline,
					seconds,
				)
			} else {
				fmt.Printf(
					"Still working on %s for %d minutes and %d seconds\n",
					headline,
					minutes,
					seconds,
				)
			}
		}
	}
}
