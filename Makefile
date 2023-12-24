ifndef $(GOPATH)
	GOPATH=$(shell go env GOPATH)
	export GOPATH
endif

build:
	go build -o bin/hpcadmin cmd/hpcadmin-cli/main.go

tidy:
	go mod tidy	

run: build
	./bin/hpcadmin

test:
	go test -v ./...
