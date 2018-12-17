#!/bin/sh

build() {
	local pkg="$1"

	gofmt -s -w "$pkg"
	go tool fix "$pkg"
	go tool vet "$pkg"

	#hash gosimple >/dev/null && gosimple "$pkg"
	hash golint >/dev/null && golint "$pkg"
	#hash staticcheck >/dev/null && staticcheck "$pkg"

	go test "$pkg"
	go install -v "$pkg"
}

build ./basgo-run
