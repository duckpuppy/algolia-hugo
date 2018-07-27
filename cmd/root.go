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
	"os"

	"github.com/apex/log"
	"github.com/duckpuppy/algolia-hugo/app"
	"github.com/kyoh86/xdg"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var config app.Config
var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "algolia-hugo",
	Short: "Easily manage your search index on Algolia for your Hugo site",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().
		StringVarP(&cfgFile, "config", "c", "", "config file (default is $XDG_CONFIG_HOME/algolia-hugo/algolia-hugo.yaml)")

	rootCmd.PersistentFlags().BoolVarP(&config.Verbose, "verbose", "v", false, "Display verbose output")
	_ = viper.BindPFlag("Verbose", rootCmd.PersistentFlags().Lookup("verbose"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	setDefaults()

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Search config in $XDG_CONFIG_HOME/algolia-hugo directory with name "algolia-hugo" (without extension).
		for _, dir := range xdg.ConfigDirs() {
			viper.AddConfigPath(fmt.Sprintf("%s/algolia-hugo", dir))
		}
		viper.AddConfigPath(fmt.Sprintf("%s/algolia-hugo", xdg.ConfigHome()))
		viper.SetConfigName("algolia-hugo")
	}
	// If a config file is found, read it in.
	_ = viper.ReadInConfig()

	viper.AutomaticEnv() // bind to environment variables that match key names

	// Unmarshal the config into the config variable
	_ = viper.Unmarshal(&config)

	// If the environment variables exist, we need to get the values and set the config
	config.AlgoliaAPIKey = viper.GetString("algolia_api_key")
	config.AlgoliaAppID = viper.GetString("algolia_app_id")
	config.AlgoliaIndexName = viper.GetString("algolia_index_name")

	if config.Verbose {
		log.WithField("config", viper.ConfigFileUsed()).Info("Loaded config")
	}
}

func setDefaults() {
	viper.SetDefault("UploadFile", "public/index.json")
	viper.SetDefault("Verbose", false)
}
