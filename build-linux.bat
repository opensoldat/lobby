@ECHO OFF

SETLOCAL

SET "GOOS=linux"
SET "GOARCH=amd64"

go build -o bin/

ENDLOCAL
