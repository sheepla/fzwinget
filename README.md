<div align="right">

[![golangci-lint](https://github.com/sheepla/fzwinget/actions/workflows/ci.yml/badge.svg)](https://github.com/sheepla/fzwinget/actions/workflows/ci.yml)
[![Release](https://github.com/sheepla/fzwinget/actions/workflows/release.yml/badge.svg)](https://github.com/sheepla/fzwinget/actions/workflows/release.yml)

</div>


<div align="center">

# fzwinget

*winget ❤ fuzzy-finder*

**fzwinget** is a wrapper command for [Windows Package Manager CLI](https://github.com/microsoft/winget-cli) a.k.a. **winget** with built-in fuzzy-finder interactive interface

</div>

<div align="center">
  <img src="https://repository-images.githubusercontent.com/594500449/fadbdca9-f764-437c-ae91-bb417cfb6d07" href="https://github.com/sheepla/fzwinget/edit/master/README.md" alt="screenshot width="60%">
</div>


## Features

- [x] built-in fuzzy-finder UI for fast selection by filtering by package ID and name
- [x] install the selected packages instantly: `fzwinget install QUERY...`
- [x] able to open the selected package link in the default web browser: `fzwinget open QUERY...`

## Usage

```
NAME:
   fzwinget - COMMAND QUERY...

USAGE:
   [OPTIONS] COMMAND QUERY...

VERSION:
   vUNKNOWN-revUNKNOWN

DESCRIPTION:
   a winget wrapper command with built-in fuzzyfiner interface

COMMANDS:
   install, i     find packages and run `winget install`
   show, s, view  find packages and show detailed informations
   open, o        find packages and open the page of the selected software(s)
   help, h        Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help (default: false)
   --version, -v  print the version (default: false)
```

## Installation

```
go install github.com/sheepla/fzwinget@latest
```

## Thanks

- [winget](https://github.com/microsoft/winget-cli) - Official Windows Package Manager CLI repository
- [winget.run](https://winget.run/) - The website finding winget packages made easy
- [winget.run API](https://github.com/winget-run/api) - The REST API behind winget.run, allowing users to search, discover. This tool utilise for this API.

## License

MIT

## Author

> [sheepla](https://github.com/sheepla/)

