all:	
	go build webserver.go admin.go generate_rss.go setup.go configreader.go

windows: admin.go webserver.go generate_rss.go
	go build -o pogoapp.exe webserver.go admin.go generate_rss.go setup.go

linux: admin.go webserver.go generate_rss.go
	go build -o pogoapp webserver.go admin.go generate_rss.go setup.go

install:
	go get github.com/gmemstr/feeds
	go get github.com/fsnotify/fsnotify
	go get github.com/spf13/viper
	go get github.com/gorilla/mux

docker:
	docker build .

and run:
	go build webserver.go admin.go generate_rss.go setup.go
	./pogoapp.exe