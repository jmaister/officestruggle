
go build
go test -v -benchmem -cpuprofile profile/cpu-%1.prof -memprofile profile/mem-%1.prof -bench=. ./%1

go tool pprof -png %1.test.exe profile/cpu-%1.prof > profile/cpu-%1.png
go tool pprof -png %1.test.exe profile/mem-%1.prof > profile/mem-%1.png