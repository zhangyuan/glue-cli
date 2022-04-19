build: clean
	go build -o build/glue cmd/glue.go

.PHONY: clean

clean:
	rm -rf build

glue-linux_amd64:
	env GOOS=linux GOARCH=amd64 go build -o build/glue-linux_amd64 cmd/glue.go
darwin_amd64:
	env GOOS=darwin GOARCH=amd64 go build -o build/glue-darwin_amd64 cmd/glue.go
windows_amd64:
	env GOOS=windows GOARCH=amd64 go build -o build/glue-windows_amd64 cmd/glue.go

release: clean glue-linux_amd64 darwin_amd64 windows_amd64
