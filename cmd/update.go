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
	"log"

	"github.com/duckpuppy/algolia-hugo/app"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
		index := config.GetIndex()

		// Open the upload file and unmarshal it before going further
		objects, err := app.LoadObjectFile(config.UploadFile)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Using index: %s: \n", config.AlgoliaIndexName)

		// First, delete all existing content in the index
		fmt.Println("Deleting existing objects")
		if err = app.ClearIndex(index); err != nil {
			log.Fatal(err)
		}

		// Next, add the new objects to the search records
		fmt.Printf("Uploading objects from %s\n", config.UploadFile)
		if _, err = index.AddObjects(objects); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
	updateCmd.Flags().StringVarP(&config.UploadFile, "file", "f", "public/index.json", "The file to upload")
	viper.BindPFlag("UploadFile", updateCmd.Flags().Lookup("file"))
}
