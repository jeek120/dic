#!/bin/bash

set CGO_ENABLED=0
set GOARCH=amd64

set GOOS=windows
go build  -o ./dic.exe ../.

set GOOS=darwin
go build -o ./dic_mac ../.

set GOOS=linux
go build -o ./dic ../.
