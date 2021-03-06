#!/bin/bash

die() {
	echo 2>&1 "$0": $#
	exit 
}

lint() {
	local pkg="$1"

	echo working: "$pkg"

	go vet -vettool=$(which shadow) "$pkg"
	gosimple "$pkg"
	golint "$pkg"
	staticcheck "$pkg"
}

go install golang.org/x/tools/go/analysis/passes/shadow/cmd/shadow

lint ./baslex
lint ./baslex-run
lint ./node
lint ./basparser
lint ./basparser-run
lint ./basgo
lint ./basgo-run
lint ./basgo-build

