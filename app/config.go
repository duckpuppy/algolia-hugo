package app

type Config struct {
	AlgoliaAPIKey    string `mapstructure:"algolia_api_key"`
	AlgoliaAppID     string `mapstructure:"algolia_app_id"`
	AlgoliaIndexName string `mapstructure:"algolia_index_name"`
	UploadFile       string
}
