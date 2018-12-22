%{

package main

import (
	"bufio"
	"fmt"
	"os"
	"unicode"
)

%}

// fields inside this union end up as the fields in a structure known
// as ${PREFIX}SymType, of which a reference is passed to the lexer.
%union{
	val string
}

// any non-terminal which returns a value needs a type, which is
// really a field name in the above union struct
%type <val> input

// same for terminals
%token <val> CHARACTER

%%

in : /* empty */
  | in input '\n'
     { fmt.Printf("Read character: %s\n", $2) }
  ;

input : CHARACTER
  | input CHARACTER
      { $$ = $1 + $2 }
  ;
%%

type InputLex struct {
        // contains one complete input string (with the trailing \n)
        s string
        // used to keep track of parser position along the above imput string
        pos int
}

func (l *InputLex) Lex(lval *InputSymType) int {
	var c rune = ' '

        // skip through all the spaces, both at the ends and in between
	for c == ' ' {
		if l.pos == len(l.s) {
			return 0
		}
		c = rune(l.s[l.pos])
		l.pos += 1
	}

        // only look for input characters that are either digits or lower case
	if unicode.IsDigit(c) || unicode.IsLower(c) {
	    lval.val = string(c)
	    return CHARACTER
	}

        // do not return any token in case of unrecognized grammer
        // this results in syntax error
	return int(c)
}

func (l *InputLex) Error(s string) {
	fmt.Printf("syntax error: %s\n", s)
}

func main() {
	fi := bufio.NewReader(os.NewFile(0, "stdin"))

	for {
		var eqn string
		var ok bool

		fmt.Printf("input: ")
		if eqn, ok = readline(fi); ok {
			InputParse(&InputLex{s: eqn})
		} else {
			break
		}
	}
}

func readline(fi *bufio.Reader) (string, bool) {
	s, err := fi.ReadString('\n')
	if err != nil {
		return "", false
	}
	return s, true
}
