BINARY_NAME=officestruggle
TS := $(shell /bin/date "+%Y%m%d-%H%M%S")


all: build test

build:
	go build -o ${BINARY_NAME} main.go

test:
	go test -v -cover ./...

testprofile: build
	go test -v -cpuprofile profile/cpu${TS}.prof -memprofile profile/mem${TS}.prof ./systems
	go tool pprof -png ${BINARY_NAME} profile/cpu${TS}.prof > profile/cpu${TS}.png
	go tool pprof -png ${BINARY_NAME} profile/mem${TS}.prof > profile/mem${TS}.png

benchmark:
	go test -v -benchmem -bench=. ./...

benchmarkprofile:
	go test -v -benchmem -cpuprofile profile/cpu${TS}.prof -memprofile profile/mem${TS}.prof -bench=. ./systems
	go tool pprof -png ${BINARY_NAME} profile/cpu${TS}.prof > profile/cpu${TS}.png
	go tool pprof -png ${BINARY_NAME} profile/mem${TS}.prof > profile/mem${TS}.png

run:
	go run -race main.go 2> err.log
 
clean:
	go clean
	rm ${BINARY_NAME}

deps:
	sudo apt install graphviz
	go install

profile: build
	mkdir -p profile
	./${BINARY_NAME} -cpuprofile profile/cpu${TS}.prof -memprofile profile/mem${TS}.prof
	go tool pprof -png ${BINARY_NAME} profile/cpu${TS}.prof > profile/cpu${TS}.png
	go tool pprof -png ${BINARY_NAME} profile/mem${TS}.prof > profile/mem${TS}.png
