SOURCEFILES = webserver.go admin.go generate_rss.go setup.go configreader.go

all:	
	go build $(SOURCEFILES)

windows:
	go build -o pogoapp.exe $(SOURCEFILES)

linux:
	go build -o pogoapp $(SOURCEFILES)

install:
	go get github.com/gmemstr/feeds
	go get github.com/fsnotify/fsnotify
	go get github.com/gorilla/mux

docker:
	docker build .
