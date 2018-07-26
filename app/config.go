package app

import (
	"github.com/algolia/algoliasearch-client-go/algoliasearch"
	"github.com/apex/log"
)

type Config struct {
	AlgoliaAPIKey    string `mapstructure:"algolia_api_key"`
	AlgoliaAppID     string `mapstructure:"algolia_app_id"`
	AlgoliaIndexName string `mapstructure:"algolia_index_name"`
	UploadFile       string `mapstructure:"upload_file"`
	Verbose          bool
}

func (c *Config) GetIndex() algoliasearch.Index {
	client := algoliasearch.NewClient(c.AlgoliaAppID, c.AlgoliaAPIKey)
	index := client.InitIndex(c.AlgoliaIndexName)
	return index
}

// LoadUploadFile loads the configured JSON file of search terms and returns a slice of algoliasearch.Objects
func (c *Config) LoadUploadFile() ([]algoliasearch.Object, error) {
	return LoadObjectFile(c.UploadFile)
}

// ClearIndex will clear the search index
func (c *Config) ClearIndex() error {
	return ClearIndex(c.GetIndex())
}

func (c *Config) UploadIndex() error {
	// Open the upload file and unmarshal it before going further
	objects, err := c.LoadUploadFile()
	if err != nil {
		log.WithError(err).WithField("file", c.UploadFile).Fatal("Failed to load the upload file")
		return err
	}

	// First, delete all existing content in the index
	log.Info("Deleting existing objects")
	if err = c.ClearIndex(); err != nil {
		log.WithError(err).Fatal("Failed to delete existing objects")
		return err
	}

	// Next, add the new objects to the search records
	log.Info("Uploading objects")
	if _, err = c.GetIndex().AddObjects(objects); err != nil {
		log.WithError(err).Fatal("Failed to upload new objects")
		return err
	}
	return nil
}
