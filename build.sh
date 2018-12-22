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

#build_lex() {
#	local pkg="$1"
#
#	go test "$pkg"
#	go install -v "$pkg"
#}
#
#go get github.com/blynn/nex
#hash nex || die missing nex
#pushd lex
#nex -s -o generated-lex.go lex.nex
#popd
#build_lex ./lex
#build ./lex-run

build ./baslex
build ./baslex-run
build ./basgo
build ./basgo-run

go get golang.org/x/tools/cmd/goyacc

#go tool yacc -o ./basparser/parser.go -p Input ./basparser/parser.y
goyacc -o ./basparser/parser.go -p Input ./basparser/parser.y

build ./basparser

