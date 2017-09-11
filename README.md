# white rabbit

[![gitgalaxy](https://img.shields.io/badge/website-gitgalaxy.com-blue.svg)](https://gitgalaxy.com) [![shield](https://img.shields.io/badge/live-podcast.gitgalaxy.com-green.svg)](https://podcast.gitgalaxy.com) [![follow](https://img.shields.io/twitter/follow/gitgalaxy.svg?style=social&label=Follow)](https://twitter.com/gitgalaxy)


podcast rss generator and cms in golang

## goal

to produce a product that is easy to deploy and easier to use when hosting a podcast from ones own servers. 

## features

 * auto-generate rss feed
 * flat-file directory structure
 * human readable files
 * self publishing interface w/ password protection
 * basic frontend for listening to episodes
 * custom css and themeing capabilities
 * json feed generation for easier parsing
 * docker support

## requirements

[github.com/gmemstr/feeds](https://github.com/gmemstr/feeds) _this branch contains some fixes for "podcast specific" tags_

[github.com/fsnotify/fsnotify](https://github.com/fsnotify/fsnotify)

[github.com/spf13/viper](https://github.com/spf13/viper)

[github.com/gorilla/mux](https://github.com/gorilla/mux)

## building

```
make install
make
./webserver
```


**non-make**
```
go get github.com/gmemstr/feeds
go get github.com/fsnotify/fsnotify
go get github.com/spf13/viper
go get github.com/gorilla/mux
go build webserver.go generate_rss.go admin.go
./webserver
```

## file format

white rabbit uses a flat file structure for managing podcast episodes. as such, files have a special naming convention.

for podcast audio files, filenames take the form of YEAR-MONTH-DAY followed by the title. The two values are
seperated by underscores (`YYYY-MM-DD_TITLE.mp3`).

shownote fils are markdown formatted and simply append `_SHOWNOTES.md` to the existing filename (sans .mp3 of course). 

### Makefile

there are several commands in the Makefile, for various reasons. (commands are preceded by the `make` command)

 * `all` - also works by just running `make`, compiles go code to executable
 * `windows` - creates named compiled .exe
 * `linux` - creates named compiled binary
 * `install` - installs go dependencies 
 * `docker` - build docker image for running elsewhere
 * `and run` - build and run the executable (remove .exe in file for \*nix)