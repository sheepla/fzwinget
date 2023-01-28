package api

import (
	"path/filepath"
	"testing"

	"github.com/sheepla/tdloader"
)

func TestParseSearchResult(t *testing.T) {
	f := tdloader.MustGetFile(filepath.Join("_testdata", "search.json"))
	defer f.Close()

	result, err := parseSearchResult(f)
	if err != nil {
		t.Fatal(err)
	}

	if result == nil {
		t.Fatal("result is nil")
	}

	for _, pkg := range *result {
		t.Log("ID:", pkg.ID)
		t.Log("Name: ", pkg.Name)
		t.Log("Versions: ", pkg.Versions)
		t.Log("Tags: ", pkg.Tags)
		t.Log("Description: ", pkg.Description)
		t.Log("Publisher: ", pkg.Publisher)
		t.Log("Featured: ", pkg.Featured)
		t.Log("Banner: ", pkg.Banner)
		t.Log("UpdatedAt: ", pkg.UpdatedAt)
	}
}
