FROM golang:stretch as builder

WORKDIR /build/pogo
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN GOOS=linux go build -ldflags "-s -w"


FROM debian:stretch-slim

RUN apt update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

WORKDIR /app
EXPOSE 3000

COPY --from=builder /build/pogo/pogo pogo

ENTRYPOINT ["/app/pogo"]
