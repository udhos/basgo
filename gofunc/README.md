
# _GOFUNC calls Go function from BASIC code

_GOFUNC() is a BASIC keyword introduced by the 'basgo' compiler in order to call a Go function from BASIC code.

# Usage

BASIC code calling Go code:

    10 result = _GOFUNC("func_name", arg1, arg2, ..., argN)
    20 print result

Go code called from BASIC:

    func func_name(arg1, arg2, ..., argN float64) float64
        return some_float64_value
    }

# Hints

- Append a '$' to _GOFUNC func_name in order to call a string-returning function.

BASIC:

    10 result$ = _GOFUNC("func_name$", arg1, arg2, ..., argN)

Go:

    func func_name(arg1, arg2, ..., argN float64) string
        return "some_string_value"
    }

- Use _GOPROC() to call a Go function with no return value.

# Example:

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

# Defining Go code within BASIC code with _GOIMPORT and _GODECL

_GOFUNC and _GOPROC can call Go function from BASIC.
However, how can that Go function be defined?
Naturally, the Go function can be defined using a full-Go source file, as depicted in the gofunc.go example above.
Another option is to embed the Go code within BASIC code, using _GOIMPORT and _GODECL.
In the example rad.bas below, _GOIMPORT includes the Go "math" package, and _GODECL is used to define the Go "degToRad" function. Then, _GOFUNC can be used to call that Go function from BASIC.
The _GOIMPORT/_GODECL option is cumbersome for large Go code, but it can be useful for defining small Go functions directly within BASIC code.

    $ more gofunc/rad.bas 
    110 rem Using _GOIMPORT and _GODECL to embed Go code within BASIC code
    120 rem
    130 _goimport("math")
    140 _godecl("func degToRad(d float64) float64 {")
    150 _godecl("    return d*math.Pi/180")
    160 _godecl("}")
    170 rem
    180 rem Now using _GOFUNC to call that Go function from BASIC code
    190 rem
    200 d = 180
    210 r = _gofunc("degToRad", d)
    220 print d;"degrees in radians is";r

Running the example:

    $ basgo-build < gofunc/rad.bas > a.go && go run a.go
     180 degrees in radians is 3.141592653589793 
    $

