build:
	go build -v ./...

test:
	go test -race -v ./...

watch-test:
	reflex -t 50ms -s -- sh -c 'gotest -race -v ./...'

install:
	go install .

coverage:
	go test -v -coverprofile cover.out .
	go tool cover -html=cover.out -o cover.html