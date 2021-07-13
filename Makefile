BINARY_NAME=officestruggle

all: build test

build:
	go build -o ${BINARY_NAME} main.go

test:
	go test -v ./...

run:
	go run -race main.go 2> err.log
 
clean:
	go clean
	rm ${BINARY_NAME}

deps:
	go install
