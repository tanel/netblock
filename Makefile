
default: build

build:
	go build

cov:
	go test -coverprofile cover.out && go tool cover -html=cover.out -o cover.html && open cover.html

test:
	go test

install:
	go install