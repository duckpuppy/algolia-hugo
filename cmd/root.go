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

	"github.com/duckpuppy/algolia-hugo/app"
	"github.com/kyoh86/xdg"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var config app.Config

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
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $XDG_CONFIG_HOME/algolia-hugo/algolia-hugo.yaml)")
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
			viper.AddConfigPath(dir)
		}
		viper.AddConfigPath(fmt.Sprintf("%s/algolia-hugo", xdg.ConfigHome()))
		viper.SetConfigName("algolia-hugo")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

	viper.Unmarshal(&config)
}

func setDefaults() {
	viper.SetDefault("UploadFile", "public/index.json")
}
