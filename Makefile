all:	
	go build -o whiterabbit.exe src/webserver.go src/admin.go src/generate_rss.go

windows: src/admin.go src/webserver.go src/generate_rss.go
	go build -o whiterabbit.exe src/webserver.go src/admin.go src/generate_rss.go

liunx: src/admin.go src/webserver.go src/generate_rss.go
	go build -o whiterabbit src/webserver.go src/admin.go src/generate_rss.go
