all:	
	go build src/webserver.go src/admin.go src/generate_rss.go

windows: src/admin.go src/webserver.go src/generate_rss.go
	go build -o whiterabbit.exe src/webserver.go src/admin.go src/generate_rss.go

liunx: src/admin.go src/webserver.go src/generate_rss.go
	go build -o whiterabbit src/webserver.go src/admin.go src/generate_rss.go

install:
	go get github.com/gmemstr/feeds
	go get github.com/fsnotify/fsnotify
	go get github.com/spf13/viper
	go get github.com/gorilla/mux