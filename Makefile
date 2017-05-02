.PHONY: default
default: all ;

all: clean lint vet test build

lint:
	for p in $$(go list ./... | grep -v /vendor/); do golint $$p; done

vet:
	for p in $$(go list ./... | grep -v /vendor/); do go vet $$p; done

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