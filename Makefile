build:
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w"

run:
	go build && ./scylla-go-sample

init:
	go mod init github.com/ragul28/scylla-go-sample
	go get -u

mod:
	go mod tidy
	go mod verify
	go mod vendor
