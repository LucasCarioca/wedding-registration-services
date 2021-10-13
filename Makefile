.PHONY := all

download:
	go mod download

start:
	go run .

test:
	go test ./pkg/... -v -covermode=count -coverprofile=coverage.out

lint:
	golint ./...

vet:
	go vet ./...

fmt:
	go fmt ./...

fmt-check:
	test -z $$(gofmt -l .)

fix:
	go fix ./...

chk:
	staticcheck ./...