
mkdir dist

set GOARCH=wasm
set GOOS=js
go build -o dist/officestruggle.wasm main.go

FOR /F "tokens=*" %%i IN ('go env GOROOT') DO set VARIABLE=%%i
echo %VARIABLE%

copy "%VARIABLE%\misc\wasm\wasm_exec.js" .\dist\
copy scripts\index.html .\dist\

python -m http.server 8000 --directory dist\