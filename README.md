## INSTALL ALL DEPENDENCY
```bash
go mod tidy
```

## RUN LOCAL COMMAND
```bash
go run main.go
```

## BUILD COMMAND
```bash
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/api cmd/main.go 
```