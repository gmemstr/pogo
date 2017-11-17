# Pogo
## Podcast RSS feed generator and CMS in Go.

[![Build Status](https://travis-ci.org/gmemstr/pogo.svg?branch=master)](https://travis-ci.org/gmemstr/pogo) [![gitgalaxy](https://img.shields.io/badge/website-gitgalaxy.com-blue.svg)](https://gitgalaxy.com) [![live branch](https://img.shields.io/badge/live-podcast.gitgalaxy.com-green.svg)](https://podcast.gitgalaxy.com) [![follow](https://img.shields.io/twitter/follow/gitgalaxy.svg?style=social&label=Follow)](https://twitter.com/gitgalaxy)

## Goal

To produce a product that is easy to deploy and easier to use when hosting a podcast from ones own servers. 

## Features

 * Auto-generate rss feed
 * Basic frontend for listening to episodes
 * Flat-file directory structure
 * Human readable files
 * Self publishing interface w/ password protection
 * Custom CSS and themeing capabilities
 * JSON feed generation for easier parsing
 * Docker support

## Requirements

[github.com/gorilla/feeds](https://github.com/gorilla/feeds)

[github.com/fsnotify/fsnotify](https://github.com/fsnotify/fsnotify)

[github.com/gorilla/mux](https://github.com/gorilla/mux)

## Building

```
godep restore
go build
# Set environment variable
export POGO_SECRET=secret
# Windows
# set POGO_SECRET=secret
./podcast
```

## File format

Pogo uses a flat file structure for managing podcast episodes. As such, files have a special naming convention.

For podcast audio files, filenames take the form of YEAR-MONTH-DAY followed by the title. The two values are
separated by underscores (`YYYY-MM-DD_TITLE.mp3`).

"Shownote" files are markdown formatted and simply append `_SHOWNOTES.md` to the existing filename (sans .mp3 of course). 
