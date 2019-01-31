
# GOFUNC calls Go function from BASIC code

## Usage:

BASIC code calling Go code:

    10 result = GOFUNC("func_name", arg1, arg2, ..., argN)
    20 print result

Go code called from BASIC:

    func func_name(arg1, arg2, ..., argN float64) float64
        return some_float64_value
    }

Append a '$' to func_name in order to call a string-returning function.

BASIC:

    10 result$ = GOFUNC("func_name$", arg1, arg2, ..., argN)

Go:

    func func_name(arg1, arg2, ..., argN float64) string
        return "some_string_value"
    }

## Example:

The Go functions:

    $ more gofunc.go
    package main
    
    func sum(a, b float64) float64 {
            return a + b
    }
    
    func concat(a, b string) string {
            return a + b
    }
    $

BASIC code using GOFUNC to call Go functions: 

    $ more gofunc.bas
    10 x=20:y=10:print gofunc("sum", sqr(4+5), x-y)
    20 x$="a":y$="b":print gofunc("concat$", x$, y$)
    $

Compile the BASIC code to Go (a.go) then build it along with the Go functions:

    $ basgo-build < gofunc.bas > a.go
    $ go run a.go gofunc.go
    13
    ab

