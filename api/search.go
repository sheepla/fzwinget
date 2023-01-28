package api

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/url"
	"strings"
	"time"

	"github.com/tidwall/gjson"
)

const (
	host = "api.winget.run"
)

type SearchParam struct {
	// parameters
	Query       string
	Name        string
	ID          string
	Publisher   string
	Description string
	Tags        string
	// options
	//EnsureContains bool
	//PartialMatch bool
	Take int
	// Order        int
}

func NewSeachParam(query string) *SearchParam {
	return &SearchParam{
		Query: query,
	}
}

func (param *SearchParam) toURL() *url.URL {
	u := &url.URL{
		Scheme: "https",
		Host:   host,
	}

	q := u.Query()

	if strings.TrimSpace(param.Query) != "" {
		q.Add("query", param.Query)
	}

	if strings.TrimSpace(param.Name) != "" {
		q.Add("name", param.Name)
	}

	if strings.TrimSpace(param.ID) != "" {
		q.Add("id", param.ID)
	}

	if strings.TrimSpace(param.Publisher) != "" {
		q.Add("publisher", param.Publisher)
	}

	if strings.TrimSpace(param.Description) != "" {
		q.Add("description", param.Description)
	}

	if strings.TrimSpace(param.Tags) != "" {
		q.Add("tags", param.Tags)
	}

	q.Add("ensureContains", "true")
	q.Add("partialMatch", "true")

	u.RawQuery = q.Encode()

	return u
}

type SearchResult []Package

type Package struct {
	Banner      string
	Publisher   string
	Description string
	Featured    bool
	Homepage    string
	ID          string
	IconURL     string
	License     string
	LicenseURL  string
	Logo        string
	Name        string
	Tags        []string
	Versions    []string
	UpdatedAt   time.Time
	CreatedAt   time.Time
}

func parseSearchResult(contentReader io.Reader) (*SearchResult, error) {
	buf := new(bytes.Buffer)
	if _, err := io.Copy(buf, contentReader); err != nil {
		return nil, fmt.Errorf("failed to load content: %w", err)
	}

	content := buf.Bytes()

	if !gjson.ValidBytes(content) {
		return nil, errors.New("invalid JSON format")
	}

	var result SearchResult
	var pkg Package

	gjson.GetBytes(content, "Packages").ForEach(func(_, packages gjson.Result) bool {
		pkg.Banner = packages.Get("Banner").String()
		pkg.Featured = packages.Get("Featured").Bool()
		pkg.ID = packages.Get("Id").String()
		pkg.IconURL = packages.Get("IconUrl").String()
		pkg.Logo = packages.Get("Logo").String()
		pkg.UpdatedAt = packages.Get("UpdatedAt").Time()
		pkg.CreatedAt = packages.Get("CreatedAt").Time()

		pkg.Versions = gjsonToStringSlice(packages.Get("Versions"))

		packages.Get("Latest").ForEach(func(_, latest gjson.Result) bool {
			pkg.Name = latest.Get("Name").String()
			pkg.Publisher = latest.Get("Publisher").String()
			pkg.Description = latest.Get("Description").String()
			pkg.Homepage = latest.Get("Homepage").String()
			pkg.License = latest.Get("License").String()
			pkg.LicenseURL = latest.Get("LicenseUrl").String()
			pkg.Tags = gjsonToStringSlice(latest.Get("Tags"))

			// continue iteration
			return true
		})

		result = append(result, pkg)

		// continue iteration
		return true
	})

	return &result, nil
}

func gjsonToStringSlice(data gjson.Result) []string {
	var arr []string

	for _, item := range data.Array() {
		arr = append(arr, item.String())
	}

	return arr
}
