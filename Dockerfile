# Use latest golang image
FROM golang:latest

RUN  mkdir -p /go/src \
  && mkdir -p /go/bin \
  && mkdir -p /go/pkg
ENV GOPATH=/go
ENV PATH=$GOPATH/bin:$PATH

# Set working directory
WORKDIR /go/src/POGO

# Add source to container so we can build
ADD . /go/src/POGO

# 1. Install make and dependencies
# 2. Install Go dependencies
# 3. Build named Linux binary and allow execution
# 4. List directory structure (for debugging really)\
# 5. Make empty podcast direcory
# 6. Create empty feed files
RUN go get godep && \
	godep restore && \
	make linux && chmod +x pogoapp && \
	ls -al && \
	mkdir podcasts && \
	touch assets/web/feed.rss assets/web/feed.json && \
	echo '{}' >assets/web/feed.json && \
	echo '{}' >assets/config/users.json && \
	echo '{}' >assets/config/config.json

EXPOSE 8000

CMD ./pogoapp