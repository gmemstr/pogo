# Use latest golang image
FROM golang:latest

RUN  mkdir -p /go/src \
  && mkdir -p /go/bin \
  && mkdir -p /go/pkg
ENV GOPATH=/go
ENV PATH=$GOPATH/bin:$PATH

# Set working directory
WORKDIR $GOPATH/src/github.com/gmemstr/pogo

# Add source to container so we can build
ADD . $GOPATH/src/github.com/gmemstr/pogo

# 1. Install make and dependencies
# 2. Install Go dependencies
# 3. Build named Linux binary and allow execution
# 4. List directory structure (for debugging really)\
# 5. Make empty podcast direcory
# 6. Create empty feed files
RUN go get github.com/tools/godep && \
	godep restore && \
	go build -o pogoapp && chmod +x pogoapp

EXPOSE 3000

CMD ./pogoapp
