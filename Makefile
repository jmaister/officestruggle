BINARY_NAME=officestruggle

all: build test

build:
	go build -o ${BINARY_NAME} main.go

test:
	go test -v main.go

run:
	go run -race main.go
 
clean:
	go clean
	rm ${BINARY_NAME}

deps:
	go install