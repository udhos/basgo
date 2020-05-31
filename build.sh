#!/bin/bash

msg() {
	echo 2>&1 "$0": $@
}

build() {
	local pkg="$1"

	gofmt -s -w "$pkg"
	go fix "$pkg"
	go vet "$pkg"

	#hash gosimple >/dev/null && gosimple "$pkg"
	#hash golint >/dev/null && golint "$pkg"
	#hash staticcheck >/dev/null && staticcheck "$pkg"

	go test "$pkg"
	go install -v "$pkg"
}

#go get golang.org/x/tools/cmd/goyacc
go get modernc.org/goyacc          ;# supports %precedence
#goyacc -o ./basparser/parser.go -p Input ./basparser/parser.y
go generate ./basparser ;# see ./basparser/generate.go

build ./basc
build ./baslex
build ./baslex-run
build ./node
build ./basparser
build ./basparser-run
build ./basgo
build ./basgo-run
build ./baslib/codepage
build ./baslib/file
build ./baslib

msg "PLEASE STAND BY: tests for basgo-build may take some minutes"
msg "                 for full testing output, stop testing and run the command below:"
msg "                 go test -test.v ./basgo-build"
build ./basgo-build

