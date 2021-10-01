.PHONY: build build-osx build-linux

build: build-osx build-linux

build-osx:
	GOOS=darwin go build -o dist/osx/reconstore cmd/reconstore/reconstore.go

build-linux:
	GOOS=linux go build -o dist/linux/reconstore cmd/reconstore/reconstore.go
