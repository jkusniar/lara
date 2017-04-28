.PHONY: default
default: test ;

build:
	go build github.com/jkusniar/lara/cmd/lara
	go build github.com/jkusniar/lara/cmd/lara-ctl

install:
	go install github.com/jkusniar/lara/cmd/lara
	go install github.com/jkusniar/lara/cmd/lara-ctl

# TODO postgres_test, or use "go test -v ./..." for all test at once
test:
	go test -v ./http_test
	go test -v ./crypto_test
	go test -v ./postgres

clean:
	go clean
	rm -fr lara
	rm -fr lara-ctl