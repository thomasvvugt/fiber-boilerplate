FROM golang:1.15

WORKDIR /go/src/app
COPY . .

RUN go get -d -v ./...

ENTRYPOINT ["/bin/bash", "-c", "go run main.go"]
