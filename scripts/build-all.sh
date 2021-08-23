

echo "Compiling..."
echo "win-x64"
GOOS=windows GOARCH=amd64 go build -o bin/main-win-x64.exe
echo "win-x32"
GOOS=windows GOARCH=386 go build -o bin/main-win-x32.exe

echo "linux-x64"
GOOS=linux GOARCH=amd64 go build -o bin/main-linux-x64.bin
echo "linux-x32"
GOOS=linux GOARCH=386 go build -o bin/main-linux-x32.bin

echo "linux-arm64"
GOOS=linux GOARCH=arm64 go build -o bin/main-linux-arm64.bin
echo "linux-arm32"
GOOS=linux GOARCH=arm go build -o bin/main-linux-arm32.bin
echo "linux-ppc64"
GOOS=linux GOARCH=ppc64 go build -o bin/main-linux-ppc64.bin

echo "mac-arm64"
GOOS=darwin GOARCH=arm64 go build -o bin/main-mac-arm64.bin
echo "mac-x64"
GOOS=darwin GOARCH=amd64 go build -o bin/main-mac-x64.bin

echo "webasm"
GOOS=js GOARCH=wasm go build -o wasm/main.wasm
