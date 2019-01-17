#!/bin/bash

die() {
	echo 2>&1 $0: $#
	exit 
}

build() {
	local pkg="$1"

	gofmt -s -w "$pkg"
	go fix "$pkg"
	go vet -vettool=$(which shadow) "$pkg"

	#hash gosimple >/dev/null && gosimple "$pkg"
	#hash golint >/dev/null && golint "$pkg"
	#hash staticcheck >/dev/null && staticcheck "$pkg"

	go test "$pkg"
	go install -v "$pkg"
}

go install golang.org/x/tools/go/analysis/passes/shadow/cmd/shadow

#go get golang.org/x/tools/cmd/goyacc
go get modernc.org/goyacc          ;# supports %precedence
goyacc -o ./basparser/parser.go -p Input ./basparser/parser.y

build ./baslex
build ./baslex-run
build ./node
build ./basparser
build ./basparser-run
build ./basgo
build ./basgo-run
build ./basgo-build

