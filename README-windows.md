# Recipe for Windows

# Requirements

## Install 7-zip

You will need 7-zip to extract mingw64 files in the next step.

If you have another archiving software that can extract .7z files, you can skip this step.

Install 7-zip from https://www.7-zip.org/

## Install mingw64

- Download x86_64-8.1.0-release-posix-sjlj-rt_v6-rev0.7z from:

https://sourceforge.net/projects/mingw-w64/files/Toolchains%20targetting%20Win64/Personal%20Builds/mingw-builds/8.1.0/threads-posix/sjlj/x86_64-8.1.0-release-posix-sjlj-rt_v6-rev0.7z

- Within the package, find the folder `mingw64` and extract it as c:\mingw64

- Add c:\mingw64\bin to the environment variable PATH

- If the installation is fine, you should be able to test it with `gcc --version`

```
C:\Users\evert>gcc --version
gcc (x86_64-posix-sjlj-rev0, Built by MinGW-W64 project) 8.1.0
Copyright (C) 2018 Free Software Foundation, Inc.
This is free software; see the source for copying conditions.  There is NO
warranty; not even for MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.

C:\Users\evert>
```

## Install Git Bash

If you have git in the system, skip this one.

Install Git Bash from https://gitforwindows.org/

Test it with `git version`:

```
C:\Users\evert>git version
git version 2.33.0.windows.2

C:\Users\evert>
```

## Install Golang

Install Golang from https://golang.org/dl/

Currently, pick this one: https://golang.org/dl/go1.17.windows-amd64.msi

Test it with `go version`

```
C:\Users\evert>go version
go version go1.17 windows/amd64

C:\Users\evert>
```

## Install basgo

- Download this file:

https://github.com/udhos/basgo/releases/download/v0.10.0/basgo_windows_amd64_v0.10.0.zip

- Extract these two files into a directory in your PATH:

```
basgo-build.exe
basc.exe
```

Test it with `basc`, like this:

```
C:\Users\evert>basc
2021/09/01 21:28:21 basc version 0.10.0 runtime go1.17 GOMAXPROC=12
basc: missing input files
usage: basc FILE [flags]
  -basgoBuild string
        basgo-build command (default "basgo-build")
  -baslibImport string
        baslib package (default "github.com/udhos/baslib/baslib")
  -baslibModule string
        baslib module (default "github.com/udhos/baslib@v0.11.0")
  -getFlags string
        go get flags
  -output string
        output file name

C:\Users\evert>
```

# Compilation

## Create a sample program

    echo 10 print "hello world!" > hello.bas

## Compile with basc

    basc hello.bas

That `basc` invokation will output the compiled code into a folder named `hello`.

## Run the compiled program

    .\hello\hello.exe

## Example

```
C:\Users\evert>echo 10 print "hello world!" > hello.bas

C:\Users\evert>basc hello.bas
2021/09/01 21:30:24 basc version 0.10.0 runtime go1.17 GOMAXPROC=12
2021/09/01 21:30:24 basc: basename: hello
2021/09/01 21:30:24 basc: cat input=hello.bas output=hello\hello.bas
2021/09/01 21:30:24 basc: basgo-build: command=basgo-build input=hello\hello.bas output=hello\hello.go baslibImport=github.com/udhos/baslib/baslib
2021/09/01 21:30:24 basgo-build version 0.10.0 runtime go1.17 GOMAXPROC=12
2021/09/01 21:30:24 basgo-build baslibImport=github.com/udhos/baslib/baslib
2021/09/01 21:30:24 basgo-build: compile: baslibImport: github.com/udhos/baslib/baslib
2021/09/01 21:30:24 basgo-build: reading BASIC code from stdin...
2021/09/01 21:30:24 basgo-build: DEBUG=[] debug=false level=0
2021/09/01 21:30:24 basgo-build: INPUTSZ=[] size=0
2021/09/01 21:30:24 basgo-build: input buffer size: 4096
2021/09/01 21:30:24 defineType: range a-z as FLOAT
2021/09/01 21:30:24 basgo-build: reading BASIC code from stdin...done
2021/09/01 21:30:24 basgo-build: FIXME WRITEME replace duplicate lines
2021/09/01 21:30:24 basgo-build: checking lines used/defined
2021/09/01 21:30:24 basgo-build: sorting lines
2021/09/01 21:30:24 basgo-build: scanning used vars
2021/09/01 21:30:24 basgo-build: issuing code
2021/09/01 21:30:24 basc: gofmt: hello\hello.go
2021/09/01 21:30:24 basc: build: dir=hello baslibModule=github.com/udhos/baslib@v0.11.0 output=
2021/09/01 21:30:24 basc: build: entering dir=hello
go: creating new go.mod: module hello
go: to add module requirements and sums:
        go mod tidy
go get: added github.com/udhos/baslib v0.11.0
2021/09/01 21:30:28 basc: output: hello\hello

C:\Users\evert>.\hello\hello.exe
2021/09/01 21:31:46 baslib: version 0.11.0 runtime go1.17 GOMAXPROC=12
2021/09/01 21:31:46 baslib: BASLIB_ALERT_OFF= showAlert=true
2021/09/01 21:31:46 baslib: env var BASLIB_ALERT_OFF is empty, set it to non-empty to disable alerts
2021/09/01 21:31:46 loading codepage 437
2021/09/01 21:31:46 loading codepage 437: found 256 symbols
2021/09/01 21:31:46 BASLIB ALERT: newInkey(): will consume os.Stdin
hello world!

C:\Users\evert>
```
