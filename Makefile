all: deps test
	go build

deps:
	go get github.com/adarqui/snappass-core-go
	go get github.com/zenazn/goji
	go get github.com/zenazn/goji/web
	go get github.com/hypebeast/gojistatic

test:
	go test
