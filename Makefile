linux: main.go
	go build -o bin/timeline main.go

windows: main.go
	GOOS=windows GOARCH=386 go build -o bin/timeline.exe main.go
