
# GOFUNC calls Go function from BASIC code

## Usage:

    10 result = GOFUNC("func_name", arg1, arg2, ..., argN)

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

