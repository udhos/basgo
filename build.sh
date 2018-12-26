#!/bin/bash

die() {
	echo 2>&1 $0: $#
	exit 
}

build() {
	local pkg="$1"

	gofmt -s -w "$pkg"
	go fix "$pkg"
	go vet "$pkg"

	#hash gosimple >/dev/null && gosimple "$pkg"
	hash golint >/dev/null && golint "$pkg"
	#hash staticcheck >/dev/null && staticcheck "$pkg"

	go test "$pkg"
	go install -v "$pkg"
}

go get golang.org/x/tools/cmd/goyacc
goyacc -o ./basparser/parser.go -p Input ./basparser/parser.y

build ./baslex
build ./baslex-run
build ./node
build ./basparser
build ./basparser-run
build ./basgo
build ./basgo-run
build ./basgo-build

