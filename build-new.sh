#!/bin/bash

go install golang.org/x/vuln/cmd/govulncheck@latest
go install golang.org/x/tools/cmd/deadcode@latest
go install github.com/mgechev/revive@latest
go install honnef.co/go/tools/cmd/staticcheck@latest
go install golang.org/x/tools/go/analysis/passes/modernize/cmd/modernize@latest
go install github.com/gordonklaus/ineffassign@latest
go install github.com/client9/misspell/cmd/misspell@latest
go install github.com/fzipp/gocyclo/cmd/gocyclo@latest

gofmt -s -w .

#revive ./...

staticcheck ./...

modernize -fix ./...

gocyclo -over 15 .

ineffassign ./...

#misspell .

go mod tidy

govulncheck ./...

#deadcode ./baslib-example

#go env -w CGO_ENABLED=0

#go get golang.org/x/tools/cmd/goyacc
go install modernc.org/goyacc@latest ;# supports %precedence
#goyacc -o ./basparser/parser.go -p Input ./basparser/parser.y
go generate ./basparser       ;# see ./basparser/generate.go

go test -failfast ./...

#go env -w CGO_ENABLED=0

go install ./...

go env -u CGO_ENABLED
