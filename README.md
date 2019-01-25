[![license](http://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/udhos/basgo/blob/master/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/udhos/basgo)](https://goreportcard.com/report/github.com/udhos/basgo)

# basgo compiles BASIC-lang to Golang

* [Requirements](#requirements)
* [Install](#install)
* [Run the Compiler](#run-the-compiler)
  * [Example](#example)
  * [Sample \- Hello World](#sample---hello-world)
* [Run the Interpreter](#run-the-interpreter)
* [References](#references)

Created by [gh-md-toc](https://github.com/ekalinin/github-markdown-toc.go)

# Requirements

In order to build the 'basgo-build' compiler, a recent version of Go is required.

If your system lacks Go, this recipe will install a current release of Go:

    git clone https://github.com/udhos/update-golang
    cd update-golang
    sudo ./update-golang.sh

# Install

The recipe below will install 'basgo-build' under "~/go/bin".

    git clone https://github.com/udhos/basgo
    cd basgo
    ./build.sh

# Run the Compiler

Status: the compiler currently can handle very simple programs.

    basgo-build < program.bas > program.go
    go run program.go                      ;# builds and runs program.go

## Example

    basgo-build < examples/game.bas > game.go
    go run game.go                            ;# builds and runs game.go

## Sample - Hello World

    $ echo '10 print "hello world!"' | basgo-build > a.go
    $
    $ go run a.go
    hello world!
    $

# Run the Interpreter

Status: the interpreter currently can only parse simple programs, but is unable to execute anything.

    # interpreter interactively reads from stdin
    basgo-run

# References

http://www.classicbasicgames.org/ - Classic BASIC Games

http://www.vintage-basic.net/games.html - BASIC Computer Games

https://hwiegman.home.xs4all.nl/gw-man/ - GW-BASIC User's Guide

http://www.antonis.de/qbebooks/gwbasman/ - GW-BASIC User's Guide

http://www.antonis.de/qbebooks/gwbasman/chapter%206.html - Operators

https://robhagemans.github.io/pcbasic/doc/1.2/#guide - Language Guide

https://github.com/robhagemans/pcbasic - GW-BASIC emulator

https://godoc.org/modernc.org/golex - lex/flex-like utility

https://github.com/skx/gobasic/ - BASIC interpreter in Golang

WHILE

END
