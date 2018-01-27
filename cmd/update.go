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
	"github.com/apex/log"
	"github.com/duckpuppy/algolia-hugo/app"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Upload a new set of search index objects from a JSON file",
	Run: func(cmd *cobra.Command, args []string) {
		index := config.GetIndex()
		ctx := log.WithFields(log.Fields{
			"file":  config.UploadFile,
			"index": config.AlgoliaIndexName,
			"cmd":   "update",
		})

		// Open the upload file and unmarshal it before going further
		objects, err := app.LoadObjectFile(config.UploadFile)
		if err != nil {
			ctx.WithError(err).Fatal("Failed to load the upload file")
		}

		// First, delete all existing content in the index
		ctx.Info("Deleting existing objects")
		if err = app.ClearIndex(index); err != nil {
			ctx.WithError(err).Fatal("Failed to delete existing objects")
		}

		// Next, add the new objects to the search records
		ctx.Info("Uploading objects")
		if _, err = index.AddObjects(objects); err != nil {
			ctx.WithError(err).Fatal("Failed to upload new objects")
		}
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
	updateCmd.Flags().StringVarP(&config.UploadFile, "file", "f", "public/index.json", "The file to upload")
	viper.BindPFlag("UploadFile", updateCmd.Flags().Lookup("file"))
}
