.PHONY: default
default: all ;

all: clean test build

build:
	go build -v github.com/jkusniar/lara/cmd/lara
	go build -v github.com/jkusniar/lara/cmd/lara-ctl

install:
	go install github.com/jkusniar/lara/cmd/lara
	go install github.com/jkusniar/lara/cmd/lara-ctl

test:
	go test -v ./...

clean:
	go clean
	rm -fr lara
	rm -fr lara-ctl