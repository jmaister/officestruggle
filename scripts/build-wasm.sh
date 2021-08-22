
mkdir -p dist

GOARCH=wasm GOOS=js go build -o dist/officestruggle.wasm main.go

cp $(go env GOROOT)/misc/wasm/wasm_exec.js ./dist/
cp scripts/index.html ./dist/

python3 -m http.server 8000 --directory dist/
