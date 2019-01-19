# basgo

# Install

    git clone https://github.com/udhos/basgo
    cd basgo
    ./build.sh

# Run the Compiler

    basgo-build < program.bas > program.go
    go run program.go

## Example

    basgo-build < examples/game.bas > game.go
    go run game.go

# Run the Interpreter

    # interpreter interactively reads from stdin
    basgo-run

# References

http://www.antonis.de/qbebooks/gwbasman/ - GW-BASIC User's Guide

http://www.antonis.de/qbebooks/gwbasman/chapter%206.html - Operators

https://robhagemans.github.io/pcbasic/doc/1.2/#guide - Language Guide

https://github.com/robhagemans/pcbasic - GW-BASIC emulator

https://godoc.org/modernc.org/golex - lex/flex-like utility
