# Copyright 2024 Google LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

all: build clean

build:
	echo "Building..."
	if [ ! -d "dist" ]; then mkdir dist; fi
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -tags netgo -ldflags '-extldflags "-static"' -o dist/linux/pasture main.go
	@if [ "$$(uname)" = "Darwin" ]; then \
		GOOS=darwin GOARCH=arm64 go build -o dist/mac/pasture-darwin-arm64 main.go; \
		GOOS=darwin GOARCH=amd64 go build -o dist/mac/pasture-darwin-amd64 main.go; \
		lipo -create -output dist/mac/pasture dist/mac/pasture-darwin-amd64 dist/mac/pasture-darwin-arm64; \
	fi

clean:
	echo "Cleaning up..."
	@if [ "$$(uname)" = "Darwin" ]; then \
		rm -f ./dist/mac/pastures-darwin-amd64 ./dist/mac/pastures-darwin-arm64; \
	fi
