// Copyright 2016-2023, Pulumi Corporation.
// Copyright 2024, Dmitry Kisler.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"log"

	neon "github.com/kislerdm/pulumi-neon/provider"
	p "github.com/pulumi/pulumi-go-provider"
)

// Serve the provider against Pulumi's Provider protocol.
func main() {
	if err := p.RunProvider(neon.Name, neon.Version, neon.Provider()); err != nil {
		log.Fatal(err)
	}
}
