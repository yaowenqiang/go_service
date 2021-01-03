SHELL := /bin/bash

run:
	go run app/sales-api/main.go
runa:
	go run app/admin/main.go
test:
	gotest -v ./... -count=1
	staticcheck ./...
tidy:
	go mod tidy
	go mod vendor
