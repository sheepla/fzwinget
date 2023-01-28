package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/sheepla/fzwinget/api"
	"github.com/sheepla/fzwinget/ui"
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

const (
	wingetDefaultFileName = "winget.exe"
)

type exitCode int

const (
	exitCodeOK exitCode = iota
	exitCodeErrArgs
	exitCodeErrAPI
	exitCodeErrFuzzyFinder
	exitCodeErrWinget
	exitCodeErrOpenURL
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
			Hidden:    false,
		},
		{
			Name:      "show",
			Aliases:   []string{"s", "view"},
			Usage:     "find packages and show detailed informations",
			UsageText: "show [OPTIONS] QUERY...",
			Flags:     searchParamFlags,
			Action:    runShowCommand,
			Hidden:    false,
		},
		{
			Name:      "open",
			Aliases:   []string{"o"},
			Usage:     "find packages and open the page of the selected software(s)",
			UsageText: "open [OPTIONS] QUERY...",
			Flags:     searchParamFlags,
			Hidden:    false,
			Action:    runOpenCommand,
		},
	}

	//app.Flags = []cli.Flag{
	//	&cli.StringFlag{
	//		Name:        "winget",
	//		Aliases:     []string{},
	//		Usage:       "specify path to winget executable file",
	//		DefaultText: "--winget path/to/winget.exe",
	//		EnvVars:     []string{"WINGET"},
	//	},
	//}

	app.Action = runRootCommand

	return app
}

func buildSearchParam(ctx *cli.Context) *api.SearchParam {
	return api.NewSearchParam(
		strings.Join(ctx.Args().Slice(), " "),
		ctx.String("id"),
		ctx.String("name"),
		ctx.String("publisher"),
		ctx.String("description"),
		ctx.String("tags"),
		50,
	)
}

func runRootCommand(ctx *cli.Context) error {
	if ctx.NArg() == 0 {
		return cli.Exit("must require arguments", exitCodeErrArgs.Int())
	}
	return cli.Exit(fmt.Sprintf("unknown command: %v", ctx.Args().Slice()), exitCodeErrArgs.Int())
}

func runInstallCommand(ctx *cli.Context) error {
	result, err := searchPackages(ctx)
	if err != nil {
		return cli.Exit(err, exitCodeOK.Int())
	}

	selected, err := ui.FindPackage(*result)
	if err != nil {
		if errors.Is(err, io.EOF) {
			return cli.Exit("", exitCodeOK.Int())
		}

		return cli.Exit(err, exitCodeErrFuzzyFinder.Int())
	}

	for _, idx := range selected {
		cmd := newWingetInstallCommand(
			ctx.App.Writer,
			ctx.App.ErrWriter,
			(*result)[idx].ID,
		)

		if err := cmd.Run(); err != nil {
			return cli.Exit(err, exitCodeErrWinget.Int())
		}
	}

	return cli.Exit("", exitCodeErrAPI.Int())
}

func runShowCommand(ctx *cli.Context) error {
	result, err := searchPackages(ctx)
	if err != nil {
		return cli.Exit(err, exitCodeOK.Int())
	}

	selected, err := ui.FindPackage(*result)
	if err != nil {
		if errors.Is(err, io.EOF) {
			return cli.Exit("", exitCodeOK.Int())
		}

		return cli.Exit(err, exitCodeErrFuzzyFinder.Int())
	}

	for _, idx := range selected {
		cmd := newWingetShowCommand(
			ctx.App.Writer,
			ctx.App.ErrWriter,
			(*result)[idx].ID,
		)

		if err := cmd.Run(); err != nil {
			return cli.Exit(err, exitCodeErrWinget.Int())
		}
	}

	return cli.Exit(err, exitCodeOK.Int())
}

func runOpenCommand(ctx *cli.Context) error {
	result, err := searchPackages(ctx)
	if err != nil {
		return cli.Exit(err, exitCodeOK.Int())
	}

	selected, err := ui.FindPackage(*result)
	if err != nil {
		if errors.Is(err, io.EOF) {
			return cli.Exit("", exitCodeOK.Int())
		}

		return cli.Exit(err, exitCodeErrFuzzyFinder.Int())
	}

	for _, idx := range selected {
		cmd := newOpenURLCommand(
			ctx.App.Writer,
			ctx.App.ErrWriter,
			(*result)[idx].Homepage,
		)

		if err := cmd.Run(); err != nil {
			return cli.Exit(err, exitCodeErrOpenURL.Int())
		}
	}

	return cli.Exit(err, exitCodeOK.Int())
}

//func runDebugCommand(ctx *cli.Context) error {
//	param := buildSearchParam(ctx)
//	fmt.Fprintf(ctx.App.Writer, "%v", param)
//
//	return nil
//}

func searchPackages(ctx *cli.Context) (*api.SearchResult, error) {
	param := buildSearchParam(ctx)
	result, err := api.Search(param)
	if err != nil {
		return nil, cli.Exit(err, exitCodeErrAPI.Int())
	}

	if len([]api.Package(*result)) == 0 {
		return nil, errors.New("no results found")
	}

	return result, nil
}

func newWingetShowCommand(stdout, stderr io.Writer, id string) *exec.Cmd {
	cmd := exec.Command("cmd.exe", "/q", "/k", "winget.exe", "show", "--id", id)
	cmd.Stdout = bufio.NewWriter(stdout)
	cmd.Stderr = bufio.NewWriter(stderr)

	return cmd
}

func newWingetInstallCommand(stdout, stderr io.Writer, id string) *exec.Cmd {
	cmd := exec.Command("cmd.exe", "/q", "/k", "winget.exe", "install", "--id", id)
	cmd.Stdout = bufio.NewWriter(stdout)
	cmd.Stderr = bufio.NewWriter(stderr)

	return cmd
}

func newOpenURLCommand(stdout, stderr io.Writer, url string) *exec.Cmd {
	cmd := exec.Command(
		"cmd.exe", "/q", "/c",
		"start", "/b", url,
	)

	cmd.Stdout = bufio.NewWriter(stdout)
	cmd.Stderr = bufio.NewWriter(stderr)

	return cmd
}
