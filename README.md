<img src="https://cdn.rawgit.com/gmemstr/pogo/ddc9ff3a/assets/web/static/logo-sm.png" alt="Pogo logo" align="right">

## Pogo
	
Podcast RSS feed generator and CMS in Go.

## Getting Started

There are a couple options for getting Pogo up and running.

- [Download the latest release](https://github.com/gmemstr/pogo/releases/latest)
- [Clone the repo and build](#building)

## Status

[![Build Status](https://travis-ci.org/gmemstr/pogo.svg?branch=master)](https://travis-ci.org/gmemstr/pogo) [![gitgalaxy](https://img.shields.io/badge/website-gitgalaxy.com-blue.svg)](https://gitgalaxy.com) [![live branch](https://img.shields.io/badge/live-podcast.gitgalaxy.com-green.svg)](https://podcast.gitgalaxy.com) [![follow](https://img.shields.io/twitter/follow/gitgalaxy.svg?style=social&label=Follow)](https://twitter.com/gitgalaxy)

## Features 

- Automatic RSS and JSON feed generation
- Frontend for listening and publishing episodes
- Multiple user support
- Custom CSS themes
- Docker support

## Running

1. [Download the latest release](https://github.com/gmemstr/pogo/releases/latest)
2. Unzip somewhere safe
3. [Edit the config](https://github.com/gmemstr/pogo/wiki/Configuration)
4. Run `pogo`
5. Navigate to your instance (`localhost:3000` by default)
6. Login to the admin interface (default: **admin**, **password1**)
7. **CHANGE YOUR PASSWORD**

## Building

_Note: [This requires a valid Go enviornment setup!](https://golang.org/doc/install)_

```
# Go get the repository
go get github.com/gmemstr/pogo

# Go to directory
cd $GOPATH/src/github.com/gmemstr/pogo

# Get godep
go get github.com/tools/godep

# Install Go dependencies
godep restore

# Build
go build

# Run
./pogo
```

## Credits

Pogo depends on several other open source projects to function.

 - [Golang](https://golang.org/)
 - [gorilla/mux](http://github.com/gorilla/mux)
 - [gorilla/feeds](http://github.com/gorilla/feeds)
 - [fsnotify/fsnotify](http://github.com/fsnotify/fsnotify)
 - [mattn/go-sqlite3](http://github.com/mattn/go-sqlite3)
