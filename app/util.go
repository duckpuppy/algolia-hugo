package app

import (
	"encoding/json"
	"os"

	"github.com/algolia/algoliasearch-client-go/algoliasearch"
	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
)

func init() {
	log.SetHandler(cli.Default)
	log.SetLevel(log.DebugLevel)
}

// LoadObjectFile loads a JSON file of search terms and returns a slice of algoliasearch.Objects
func LoadObjectFile(file string) ([]algoliasearch.Object, error) {
	jsonfile, err := os.Open(file)
	if err != nil {
		return []algoliasearch.Object{}, err
	}

	dec := json.NewDecoder(jsonfile)
	var objects []algoliasearch.Object
	err = dec.Decode(&objects)
	if err != nil {
		return []algoliasearch.Object{}, err
	}

	return objects, nil
}

// ClearIndex will clear the search index
func ClearIndex(index algoliasearch.Index) error {
	_, err := index.Clear()
	return err
}
