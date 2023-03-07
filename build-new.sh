#!/bin/bash

gofmt -s -w .

#revive ./...

go mod tidy

#go get golang.org/x/tools/cmd/goyacc
go install modernc.org/goyacc@latest ;# supports %precedence
#goyacc -o ./basparser/parser.go -p Input ./basparser/parser.y
go generate ./basparser       ;# see ./basparser/generate.go

go test -failfast ./...

go install ./...
