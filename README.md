[![license](http://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/udhos/basgo/blob/master/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/udhos/basgo)](https://goreportcard.com/report/github.com/udhos/basgo)

# basgo compiles BASIC-lang to Golang

* [Requirements](#requirements)
* [Install](#install)
* [Run the Compiler](#run-the-compiler)
  * [Status and Limitations](#status-and-limitations)
  * [Example](#example)
  * [Sample \- Hello World](#sample---hello-world)
  * [Use \_GOFUNC to call Go function from BASIC code](#use-_gofunc-to-call-go-function-from-basic-code)
* [Run the Interpreter](#run-the-interpreter)
* [BASIC References](#basic-references)
  * [BASIC programs and games](#basic-programs-and-games)
  * [BASIC documentation](#basic-documentation)
  * [BASIC interpreters and compilers](#basic-interpreters-and-compilers)

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

    basgo-build < program.bas > program.go
    go run program.go                      ;# builds and runs program.go

## Status and Limitations

The compiler currently can handle many simple programs.

Limitations include lack of support for sound, graphics and hardware-specific instructions (POKE, PEEK, etc).

See also known issues: https://github.com/udhos/basgo/issues

## Example

    basgo-build < examples/game.bas > game.go
    go run game.go                            ;# builds and runs game.go

## Sample - Hello World

    $ echo '10 print "hello world!"' | basgo-build > a.go
    $
    $ go run a.go
    hello world!
    $

## Use \_GOFUNC to call Go function from BASIC code

\_GOFUNC() is a BASIC keyword introduced by the 'basgo' compiler in order to call a Go function from BASIC code.

    10 result = _GOFUNC("func_name", arg1, arg2, ..., argN)
    20 print result

See [gofunc](gofunc)

# Run the Interpreter

Status: the interpreter currently can only parse simple programs, but is unable to execute anything.

    # interpreter interactively reads from stdin
    basgo-run

# BASIC References

## BASIC programs and games

https://www.completelyfreesoftware.com/old_games.html - A Collection Of 1980s Games

http://www.dunnington.info/public/startrek/index.html - Star Trek

https://sparcie.wordpress.com/tag/gwbasic/ - Few GW-BASIC games

http://www.eddiesegoura.com/Games/ - BASIC Games

http://peyre.x10.mx/GWBASIC/ - A page about GWBASIC Games & Other Programs

http://www.ifarchive.org/indexes/if-archive/games/source/basic/

http://www.moorecad.com/classicbasic/index.html - Classic Basic Games Page

http://www.classicbasicgames.org/ - Classic BASIC Games

http://www.vintage-basic.net/games.html - BASIC Computer Games

## BASIC documentation

https://hwiegman.home.xs4all.nl/gw-man/ - GW-BASIC User's Guide

http://www.antonis.de/qbebooks/gwbasman/ - GW-BASIC User's Guide

http://www.antonis.de/qbebooks/gwbasman/chapter%206.html - Operators

https://robhagemans.github.io/pcbasic/doc/1.2/#guide - Language Guide

http://www.worldofspectrum.org/ZXBasicManual/ - SINCLAIR ZX SPECTRUM - BASIC Programming

## BASIC interpreters and compilers

https://github.com/robhagemans/pcbasic - GW-BASIC emulator

https://github.com/skx/gobasic/ - BASIC interpreter in Golang

