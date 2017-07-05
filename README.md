# whiterabbit
Podcast RSS/CMS in Golang

## requirements

[github.com/gmemstr/feeds](https://github.com/gmemstr/feeds)
[github.com/fsnotify/fsnotify](https://github.com/fsnotify/fsnotify)
[github.com/spf13/viper](https://github.com/spf13/viper)
[github.com/gorilla/mux](https://github.com/gorilla/mux)

## building

```
go get github.com/gmemstr/feeds
go get github.com/fsnotify/fsnotify
go get github.com/spf13/viper
go get github.com/gorilla/mux
go build webserver.go generate_rss.go
./webserver
```
