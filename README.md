<img src="https://cdn.rawgit.com/gmemstr/pogo/ddc9ff3a/assets/web/static/logo-sm.png" alt="Pogo logo" align="right">

## Pogo

[![Build Status](https://travis-ci.org/gmemstr/pogo.svg?branch=master)](https://travis-ci.org/gmemstr/pogo) [![Demo](https://img.shields.io/badge/demo-demo.pogoapp.net-green.svg)](http://demo.pogoapp.net)
	
Podcast RSS feed generator and CMS in Go.

## Getting Started

There are a couple options for getting Pogo up and running.

- [Download the latest release](#running)
- [Clone the repo and build](#building)

## Features 

- Automatic RSS and JSON feed generation
- Frontend for listening and publishing episodes
- Multiple user support
- Custom CSS themes
- Simple deployment
- SPA frontend

## Running

[![asciicast](https://asciinema.org/a/GTsDD2nJeLktUAbrQq4hyziOt.png)](https://asciinema.org/a/GTsDD2nJeLktUAbrQq4hyziOt)

1. [Download the latest release](https://github.com/gmemstr/pogo/releases/latest)
2. Run `pogo`

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
