package app

import (
	"encoding/json"
	"os"

	"github.com/algolia/algoliasearch-client-go/algoliasearch"
)

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

func ClearIndex(index algoliasearch.Index) error {
	params := algoliasearch.Map{}
	return index.DeleteByQuery("", params)
}
