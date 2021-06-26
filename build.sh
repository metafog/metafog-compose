#!/bin/zsh

env CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 go build -ldflags "-s -w" -a -o ./bin/macos/planetr-compose cmd/task/task.go
upx ./bin/macos/planetr-compose

env CGO_ENABLED=1 GOOS=windows GOARCH=amd64 CC="x86_64-w64-mingw32-gcc" go build -a -ldflags "-linkmode external -extldflags '-static' -s -w" -o ./bin/win64/planetr-compose.exe cmd/task/task.go
upx ./bin/win64/planetr-compose.exe

env CGO_ENABLED=1 GOOS=linux GOARCH=amd64 CC="x86_64-linux-musl-gcc" go build -a -ldflags "-linkmode external -extldflags '-static' -s -w"  -o ./bin/linux-amd64/planetr-compose cmd/task/task.go
upx ./bin/linux-amd64/planetr-compose
