

go build
go test -v -benchmem -cpuprofile profile/cpu-astar.prof -memprofile profile/mem-astar.prof -bench=. ./astar

go tool pprof -png astar.test.exe profile/cpu-astar.prof > profile/cpu-astar.png
go tool pprof -png astar.test.exe profile/mem-astar.prof > profile/mem-astar.png