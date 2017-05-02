.PHONY: default
default: all ;

# default target
all: clean lint vet test build

# install build dependencies, REQUIRED
build-deps:
	go get -u -v github.com/golang/lint/golint

# install development dependencies, OPTIONAL
dev-deps:
	go get -u -v github.com/kardianos/govendor
	go get -u -v golang.org/x/tools/cmd/stringer

# call golint on sources
lint:
	for p in $$(go list ./... | grep -v /vendor/); do golint $$p; done

# call go vet on sources
vet:
	for p in $$(go list ./... | grep -v /vendor/); do go vet $$p; done

# build executables on default arch
build:
	go build -v github.com/jkusniar/lara/cmd/lara
	go build -v github.com/jkusniar/lara/cmd/lara-ctl

# build & install executables on default arch
install:
	go install github.com/jkusniar/lara/cmd/lara
	go install github.com/jkusniar/lara/cmd/lara-ctl

# run all tests
test:
	go test -v ./...

# clean build artifacts
clean:
	go clean
	rm -fr lara
	rm -fr lara-ctl