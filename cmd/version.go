// Copyright © 2018 Patrick Aikens
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

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	Version string
	Build   string
	Branch  string
	Commit  string
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display application version",
	Run: func(cmd *cobra.Command, args []string) {
		var version string
		if Branch != "master" {
			version = fmt.Sprintf("%s-%s-%s", Version, Branch, Commit)
		} else {
			version = Version
		}

		fmt.Printf("Version: %s\nBuilt:   %s\n", version, Build)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
