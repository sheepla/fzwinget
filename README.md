# fzwinget

<div align="center">

winget ❤ fuzzyfinder

*fzwinget* is a wrapper command for [winget](microsoft/winget-cli) with built-in fuzzy-finder interactive interface.

</div>


## Features

- [x] built-in fuzzy-finder UI for fast selection by filtering by package ID and name
- [x] Instantly install selected packages: `fzwinget install QUERY...`
- [x] able to open the selected package link in the default web browser: `fzwinget open QUERY`

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
