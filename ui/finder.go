package ui

import (
	"fmt"

	fzf "github.com/ktr0731/go-fuzzyfinder"
	"github.com/sheepla/fzwinget/api"
)

func FindPackage(result api.SearchResult) ([]int, error) {
	return fzf.FindMulti(
		result,
		func(i int) string {
			return fmt.Sprintf("[%v] %v", result[i].ID, result[i].Name)
		},

		fzf.WithMode(fzf.ModeCaseInsensitive),
		fzf.WithPreviewWindow(func(i, width, height int) string {
			return fmt.Sprintf(
				"[%v]\n%v\n\n───────────────────\n\npublished by %v\n\nlicenced under the %v\n\ntags: %v\nversions: %v\n\nURL: %v\nLicense URL: %v",
				result[i].ID,
				result[i].Name,
				result[i].Publisher,
				result[i].License,
				result[i].Tags,
				result[i].Versions,
				//humanize.Time(result[i].CreatedAt),
				//humanize.Time(result[i].UpdatedAt),
				result[i].Homepage,
				result[i].LicenseURL,
			)
		}),
	)
}
