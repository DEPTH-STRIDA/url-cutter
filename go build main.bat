@echo off
setlocal
set GOOS=linux
set GOARCH=amd64
go build main.go
endlocal
