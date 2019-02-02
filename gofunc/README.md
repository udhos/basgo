
# _GOFUNC calls Go function from BASIC code

_GOFUNC() is a BASIC keyword introduced by the 'basgo' compiler in order to call a Go function from BASIC code.

## Usage:

BASIC code calling Go code:

    10 result = _GOFUNC("func_name", arg1, arg2, ..., argN)
    20 print result

Go code called from BASIC:

    func func_name(arg1, arg2, ..., argN float64) float64
        return some_float64_value
    }

## Hints

- Append a '$' to _GOFUNC func_name in order to call a string-returning function.

- Use _GOPROC() to call a Go function with no return value.

BASIC:

    10 result$ = _GOFUNC("func_name$", arg1, arg2, ..., argN)

Go:

    func func_name(arg1, arg2, ..., argN float64) string
        return "some_string_value"
    }

## Example:

The Go functions:

    $ more gofunc.go
    package main

    import (
            "fmt"
    )

    func sum(a, b float64) float64 {
            return a + b
    }

    func concat(a, b string) string {
            return a + b
    }

    func printString(s string) {
            fmt.Print(s)
    }
    $

BASIC code using _GOFUNC and _GOPROC to call Go functions: 

    $ more gofunc.bas
    10 print _gofunc("sum", sqr(4+5), 20-10)
    20 print _gofunc("concat$", "a", "b")
    30 _goproc("printString", "c"): print
    $

Compile the BASIC code to Go (a.go) then build it along with the Go functions:

    $ basgo-build < gofunc.bas > a.go
    $ go run a.go gofunc.go
    13
    ab
    c
    $

