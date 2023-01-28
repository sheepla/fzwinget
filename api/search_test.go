package api

import (
	"path/filepath"
	"testing"

	"github.com/sheepla/tdloader"
)

type ParamTestCase []struct {
	SearchParam SearchParam
	WantURL     string
}

func TestSearchParamToURL(t *testing.T) {
	testcase := ParamTestCase{
		{
			SearchParam: SearchParam{
				ID: "Neovim.Neovim",
			},
			WantURL: "https://api.winget.run/v2/packages?ensureContains=true&id=Neovim.Neovim&partialMatch=true",
		},
		{
			SearchParam: SearchParam{
				Query: "web browser",
			},
			WantURL: "https://api.winget.run/v2/packages?ensureContains=true&partialMatch=true&query=web+browser",
		},
		{
			SearchParam: SearchParam{
				Publisher: "Microsoft",
			},
			WantURL: "https://api.winget.run/v2/packages?ensureContains=true&partialMatch=true&publisher=Microsoft",
		},
	}

	for _, v := range testcase {
		have := v.SearchParam.toURL().String()
		want := v.WantURL
		if have != want {
			t.Fatal("have:", have)
			t.Fatal("want:", want)
		}
	}

}

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
