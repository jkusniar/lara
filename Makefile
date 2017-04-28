.PHONY: default
default: test ;

build:
	go build github.com/jkusniar/lara/cmd/lara
	go build github.com/jkusniar/lara/cmd/lara-ctl

install:
	go install github.com/jkusniar/lara/cmd/lara
	go install github.com/jkusniar/lara/cmd/lara-ctl

test:
	go test -v ./...

clean:
	go clean
	rm -fr lara
	rm -fr lara-ctl