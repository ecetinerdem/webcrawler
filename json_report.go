package main

import (
	"encoding/json"
	"os"
	"sort"
)

func writeJSONReport(pages map[string]PageData, filename string) error {

	var pagesSlice []PageData
	var pageKeys []string

	for k, _ := range pages {
		pageKeys = append(pageKeys, k)
	}

	sort.Strings(pageKeys)

	for _, key := range pageKeys {
		pagesSlice = append(pagesSlice, pages[key])
	}

	data, err := json.MarshalIndent(pagesSlice, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(filename, data, 0644)
	if err != nil {
		return err
	}

	return nil
}
