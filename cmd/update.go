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
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/algolia/algoliasearch-client-go/algoliasearch"
	"github.com/spf13/cobra"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		client := algoliasearch.NewClient(config.AlgoliaAppID, config.AlgoliaAPIKey)
		index := client.InitIndex(config.AlgoliaIndexName)

		// First, delete all existing content in the index
		fmt.Printf("Using index: %s: \n", config.AlgoliaIndexName)

		fmt.Println("Deleting existing objects")
		params := algoliasearch.Map{}
		err := index.DeleteByQuery("", params)
		if err != nil {
			log.Fatal(err)
		}

		// Next, upload the new index
		jsonfile, err := os.Open(config.UploadFile)
		if err != nil {
			log.Fatal(err)
		}
		dec := json.NewDecoder(jsonfile)
		var objects []algoliasearch.Object
		err = dec.Decode(&objects)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Uploading objects from %s\n", config.UploadFile)
		res, err := index.AddObjects(objects)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(res)
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
	updateCmd.Flags().StringVarP(&config.UploadFile, "file", "f", "public/index.json", "The file to upload")
}
