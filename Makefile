BINARY_NAME=officestruggle

all: build test

build:
	go build -o ${BINARY_NAME} main.go

test:
	go test -v -cover ./...

benchmark:
	go test -v -benchmem -bench=. ./...

run:
	go run -race main.go 2> err.log
 
clean:
	go clean
	rm ${BINARY_NAME}

deps:
	go install

profile: build
	${BINARY_NAME} -cpuprofile profile/cpu.prof -memprofile profile/mem.prof
	go tool pprof --pdf ${BINARY_NAME} profile/cpu.prof > profile/cpu.pdf
	go tool pprof --pdf ${BINARY_NAME} profile/mem.prof > profile/mem.pdf
