
all:	
	go build -o pogoapp

windows:
	go build -o pogoapp.exe 

linux:
	go build -o pogoapp

install:
	go get github.com/gmemstr/feeds
	go get github.com/fsnotify/fsnotify
	go get github.com/gorilla/mux

docker:
	docker build .
