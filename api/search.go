package api

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/url"
	"path"
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
		Path:   path.Join("v2", "packages"),
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

func Search(param *SearchParam) (*SearchResult, error) {
	body, err := fetch(param.toURL())
	if err != nil {
		return nil, err
	}

	defer body.Close()

	result, err := parseSearchResult(body)
	if err != nil {
		return nil, err
	}

	return result, nil
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

	gjson.GetBytes(content, "Packages").ForEach(func(_, value gjson.Result) bool {
		pkg.Banner = value.Get("Banner").String()
		pkg.Featured = value.Get("Featured").Bool()
		pkg.ID = value.Get("Id").String()
		pkg.IconURL = value.Get("IconUrl").String()
		pkg.Logo = value.Get("Logo").String()
		pkg.UpdatedAt = value.Get("UpdatedAt").Time()
		pkg.CreatedAt = value.Get("CreatedAt").Time()

		pkg.Versions = gjsonToStringSlice(value.Get("Versions"))

		pkg.Name = value.Get("Latest.Name").String()
		pkg.Publisher = value.Get("Latest.Publisher").String()
		pkg.Description = value.Get("Latest.Description").String()
		pkg.Homepage = value.Get("Latest.Homepage").String()
		pkg.License = value.Get("Latest.License").String()
		pkg.LicenseURL = value.Get("Latest.LicenseUrl").String()
		pkg.Tags = gjsonToStringSlice(value.Get("Latest.Tags"))

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
