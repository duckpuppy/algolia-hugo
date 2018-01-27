package app

import "github.com/algolia/algoliasearch-client-go/algoliasearch"

type Config struct {
	AlgoliaAPIKey    string `mapstructure:"algolia_api_key"`
	AlgoliaAppID     string `mapstructure:"algolia_app_id"`
	AlgoliaIndexName string `mapstructure:"algolia_index_name"`
	UploadFile       string
}

func (c *Config) GetIndex() algoliasearch.Index {
	client := algoliasearch.NewClient(c.AlgoliaAppID, c.AlgoliaAPIKey)
	index := client.InitIndex(c.AlgoliaIndexName)
	return index
}
