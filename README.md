[![license](http://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/udhos/basgo/blob/master/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/udhos/basgo)](https://goreportcard.com/report/github.com/udhos/basgo)

# basgo compiles BASIC-lang to Golang.

**basgo** compiles BASIC-lang to Golang. Then `go build` can translate the code to native executable binary.

* [basgo compiles BASIC\-lang to Golang\.](#basgo-compiles-basic-lang-to-golang)
* [Requirements](#requirements)
  * [Version 0\.5 requires GCC](#version-05-requires-gcc)
  * [Install mingw64 to provide GCC for Windows](#install-mingw64-to-provide-gcc-for-windows)
* [Install](#install)
  * [Install only the 'basgo\-build' compiler](#install-only-the-basgo-build-compiler)
  * [Full install for development](#full-install-for-development)
* [Run the Compiler](#run-the-compiler)
  * [Run the script 'basc'](#run-the-script-basc)
  * [Run the compiler manually](#run-the-compiler-manually)
  * [Status and Limitations](#status-and-limitations)
  * [Example](#example)
  * [Sample \- Hello World](#sample---hello-world)
  * [Use \_GOFUNC to call Go function from BASIC code](#use-_gofunc-to-call-go-function-from-basic-code)
* [Run the Interpreter](#run-the-interpreter)
* [BASIC References](#basic-references)
  * [BASIC programs and games](#basic-programs-and-games)
  * [BASIC documentation](#basic-documentation)
  * [BASIC interpreters and compilers](#basic-interpreters-and-compilers)
* [2D Graphics Packages](#2d-graphics-packages)

Created by [gh-md-toc](https://github.com/ekalinin/github-markdown-toc.go)

# Requirements

In order to build the 'basgo-build' compiler, a recent version of Go is required.

If your system lacks Go, this recipe will install a current release of Go:

    git clone https://github.com/udhos/update-golang
    cd update-golang
    sudo ./update-golang.sh

For Windows systems, get the Go installer here: https://golang.org/dl/

## Version 0.5 requires GCC

Versions up to 0.4 of 'basgo-build' compiler did not require GCC.

In version 0.5 the experimental support for graphics introduced GCC as requirement.

## Install mingw64 to provide GCC for Windows

This is a quick recipe on how to install mingw64 on Windows.

- Download x86_64-8.1.0-release-posix-sjlj-rt_v6-rev0.7z from:

https://sourceforge.net/projects/mingw-w64/files/Toolchains%20targetting%20Win64/Personal%20Builds/mingw-builds/8.1.0/threads-posix/sjlj/x86_64-8.1.0-release-posix-sjlj-rt_v6-rev0.7z

- Extract the folder 'mingw64' as c:\mingw64

- Add c:\mingw64\bin to %PATH%

- Test GCC:

Open CMD.exe and run 'gcc --version':

    C:\Users\evert>gcc --version
    gcc (x86_64-posix-sjlj-rev0, Built by MinGW-W64 project) 8.1.0
    Copyright (C) 2018 Free Software Foundation, Inc.
    This is free software; see the source for copying conditions.  There is NO
    warranty; not even for MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
    
    
    C:\Users\evert>

# Install

If you don't want to hack the compiler, you can download a precompiled binary release here:

https://github.com/udhos/basgo/releases

## Install only the 'basgo-build' compiler

The recipe below will install 'basgo-build' under "~/go/bin".

    git clone https://github.com/udhos/basgo
    cd basgo
    go get modernc.org/goyacc
    go generate ./basparser
    go install ./basgo-build

## Full install for development

If you want to hack the compiler, perform a full build (including tests):

    git clone https://github.com/udhos/basgo
    cd basgo
    ./build.sh

# Run the Compiler

## Run the script 'basc'

The utility 'basc' performs the full compilation steps automatically:

    echo '10 print "hello world"' > hello.bas
    basc hello.bas                            ;# compile hello.bas to ./hello/hello
    ./hello/hello                             ;# execute the resulting binary

## Run the compiler manually

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

# 2D Graphics Packages

https://github.com/fyne-io/fyne - UI toolkit

https://github.com/faiface/pixel - 2D game library

https://github.com/fogleman/gg - 2D rendering only, does not send to screen

