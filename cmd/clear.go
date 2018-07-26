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

	"github.com/apex/log"

	"github.com/spf13/cobra"
)

// clearCmd represents the clear command
var clearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Clear all the contents of the configured index",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Clearing index: %s\n", config.AlgoliaIndexName)
		if err := config.ClearIndex(); err != nil {
			log.WithError(err).Fatal("Failed to clear index")
		}
	},
}

func init() {
	rootCmd.AddCommand(clearCmd)
}
