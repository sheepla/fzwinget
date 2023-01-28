package main

import (
	"fmt"
	"os"

	"github.com/sheepla/fzwinget/api"
	cli "github.com/urfave/cli/v3"
)

var (
	appName        = "fzwinget"
	appUsage       = "COMMAND QUERY..."
	appDescription = "a winget wrapper command with built-in fuzzyfiner interface"
	appUsageText   = "[OPTIONS] COMMAND QUERY..."
	appVersion     = "UNKNOWN"
	appRevision    = "UNKNOWN"
)

const (
	searchParamCategoryName = "search parameters"
)

type exitCode int

const (
	exitCodeOK exitCode = iota
	exitCodeErrArgs
	exitCodeErrAPI
)

func (e exitCode) Int() int {
	return int(e)
}

func main() {
	if err := initApp().Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

func initApp() *cli.App {
	app := &cli.App{
		Name:        appName,
		Usage:       appUsage,
		UsageText:   appUsageText,
		Description: appDescription,
		Version:     fmt.Sprintf("v%v-rev%v", appVersion, appRevision),
	}

	searchParamFlags := []cli.Flag{
		&cli.StringFlag{
			Name:     "id",
			Aliases:  []string{"i"},
			Usage:    "specify package ID for search parameters",
			Category: searchParamCategoryName,
		},
		&cli.StringFlag{
			Name:     "name",
			Aliases:  []string{"n"},
			Usage:    "specify package name for search parameters",
			Category: searchParamCategoryName,
		},
		&cli.StringFlag{
			Name:     "description",
			Aliases:  []string{"d", "desc"},
			Usage:    "specify package description for search parameters",
			Category: searchParamCategoryName,
		},
		&cli.StringFlag{
			Name:     "publisher",
			Aliases:  []string{"p", "pub"},
			Usage:    "specify package publisher for search parameters",
			Category: searchParamCategoryName,
		},
		&cli.StringFlag{
			Name:     "tags",
			Aliases:  []string{"t"},
			Usage:    "specify package tags for search parameters separated by commas e.g. `foo,bar,baz`",
			Category: searchParamCategoryName,
		},
	}

	app.Commands = []*cli.Command{
		{
			Name:      "install",
			Aliases:   []string{"i"},
			Usage:     "find packages and run `winget install`",
			UsageText: "install [OPTIONS] QUERY...",
			Flags:     searchParamFlags,
			Action:    runInstallCommand,
		},
		{
			Name:      "show",
			Aliases:   []string{"s", "view"},
			Usage:     "find packages and show detailed informations",
			UsageText: "show QUERY...",
			Flags:     searchParamFlags,
			Action:    runShowCommand,
		},
		{
			Name:      "open",
			Aliases:   []string{"o"},
			Usage:     "find packages and open the page of the selected software(s)",
			UsageText: "open QUERY...",
			Flags:     searchParamFlags,
		},
	}

	app.Action = runRootCommand

	return app
}

func buildSearchParam(ctx *cli.Context) *api.SearchParam {
	return &api.SearchParam{
		ID:          ctx.String("id"),
		Name:        ctx.String("name"),
		Description: ctx.String("description"),
		Publisher:   ctx.String("publisher"),
		Tags:        ctx.String("tags"),
	}

}

func runRootCommand(ctx *cli.Context) error {
	if ctx.NArg() == 0 {
		return cli.Exit("must require arguments", exitCodeErrArgs.Int())
	}
	return cli.Exit(fmt.Sprintf("unknown command: %v", ctx.Args().Slice()), exitCodeOK.Int())
}

func runInstallCommand(ctx *cli.Context) error {
	result, err := api.Search(buildSearchParam(ctx))
	if err != nil {
		return cli.Exit(err, exitCodeErrAPI.Int())
	}

	fmt.Println(result)

	return nil
}

func runShowCommand(ctx *cli.Context) error {

	return nil
}
