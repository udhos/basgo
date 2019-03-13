package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
)

type buildTest struct {
	source     string
	input      string
	output     string
	buildError int
}

const sourceGoto = `
10 goto 900
700 print "3"
710 end 
800 print "2"
810 goto 700
900 print "1"
910 goto 800
`

const outputGoto = `1
2
3
`

const sourceOnGosub = `
10 for i=1 to 3:on i gosub 100,200,300:next
20 print "end";
30 end
100 print "1";:gosub 1000:return
200 print "2";:return
300 print "3";:return
1000 print "push";:return
`

const sourceOnGosub2 = `
5 rem output: 'a 1 2 3b 4'
10 for i=1 to 4:on i gosub 100,,,400:print str$(i);:next:end
100 print "a";:return
400 print "b";:return
`

const sourceOnGoto = `
5 rem output: 'x'
10 on 2 goto 100,,300:print "x";:end
100 print "a";:end
300 print "b";:end
`

const sourceGoto2 = `
10 goto 30
20 def fna(x)=x*x:print fna(2)
30 print 3
`

const sourceGofunc = `
10 _goimport("math")
20 _godecl("func degToRad(d float64) float64 {")
30 _godecl("    return d*math.Pi/180")
40 _godecl("}")
50 input d
60 r = _gofunc("degToRad", d)
70 print d;"degrees in radians is";r;
`

const sourceReturnLine = `
10 print "1";:gosub 100
20 print "2";
30 print "3";
40 end
100 print "4";
110 return 30 
`

const sourceRestore = `
10 data 1,2
20 data 3,4:data 5,6
30 data 7,8
40 read a:print a;
50 restore 30:read a:print a;
60 restore 20:read a:print a;
`

const sourceDefint = `
10 a=1.1:print a;
20 defstr a:a="b":print a;
30 defint a:a=2.2:print a;
`

const (
	OK      = iota
	WRONG   = iota
	RUNTIME = iota
)

var testTable = []buildTest{
	{"", "", "", OK},                  // empty program
	{"ugh", "", "", WRONG},            // invalid program
	{`10 print "ab"`, "", "ab\n", OK}, // minimum program

	{`10 a$="/tmp/x":open a$ for output as 1:print#1,"xyz";:close:open a$ for append as 1:print#1,"abc";:close:open a$ for input as 1:print input$(6,#1);`, "", "xyzabc", OK},

	{`10 print len(a%);`, "", " 8 ", OK},
	{`10 print len(1.1);`, "", " 8 ", OK},
	{`10 print len("123456");`, "", " 6 ", OK},

	{`10 print atn(0);`, "", " 0 ", OK},
	{`10 print atn(1);`, "", " 0.7853981633974483 ", OK},
	{`10 print tan(atn(1));`, "", " 1 ", OK},

	{`10 a$="/tmp/x":b$="/tmp/y":open a$ for output as 1:print#1,"xyz":close:kill b$:name a$ as b$:open b$ for input as 1:print input$(3,#1);`, "", "xyz", OK},

	{`10 print oct$(7);oct$(8);`, "", "710", OK},
	{`10 print hex$(9);hex$(15);hex$(&haa);`, "", "9FAA", OK},

	{`10 chdir "/etc":files "passwd"`, "", "passwd\n", OK},

	{`10 a$="/tmp/x":open a$ for output as 1:print#1,"xyz":close:open a$ for input as 1:print lof(1);`, "", " 4 ", OK},

	{`10 a$="/tmp/x":open a$ for output as 1:print#1,"xyz":close:kill a$:open a$ for input as 1:print input$(2,#1);`, "", "", OK},

	{`10 a$="/tmp/x":open a$ for output as 1:print#1,"xyz":close:open a$ for input as 1:print input$(2,#1);`, "", "xy", OK},

	{`10 a$="/tmp/x":open a$ for output as 1:print#1,"a":print#1,"b":close:open a$ for input as 1:input#1,x$,y$:print x$;"-";y$;`, "", "a-b", OK},

	{`10 print environ$("basgo-ttt");`, "", "", OK},
	{`10 environ "basgo-ttt=zz":print environ$("basgo-ttt");`, "", "zz", OK},
	{`10 environ "basgo-ttt=zz":environ "basgo-ttt=22":print environ$("basgo-ttt");`, "", "22", OK},
	{`10 environ "basgo-ttt=zz":environ "basgo-ttt=":print environ$("basgo-ttt");`, "", "", OK},

	{`10 print eof(1);`, "", " -1 ", OK},
	{`10 open "/etc/passwd" for input as 1:print eof(1);`, "", " 0 ", OK},

	{`10 print int(1);`, "", " 1 ", OK},
	{`10 print fix(1);`, "", " 1 ", OK},
	{`10 a%=1:print int(a%);`, "", " 1 ", OK},
	{`10 a%=1:print fix(a%);`, "", " 1 ", OK},
	{`10 print int(-8.4);`, "", " -9 ", OK},
	{`10 print fix(-8.4);`, "", " -8 ", OK},
	{`10 print int(8.9);`, "", " 8 ", OK},
	{`10 print fix(8.9);`, "", " 8 ", OK},

	{`10 erase a$`, "", "", WRONG},
	{`10 dim a$(20):a$(1)="a":print a$(1);:print a$(1);`, "", "aa", OK},
	{`10 dim a$(20):a$(1)="a":print a$(1);:erase a$:print a$(1);`, "", "a", OK},

	{`10 print pos(0);`, "", " 1 ", OK},
	{`10 print pos(0);pos(0);`, "", " 1  4 ", OK},
	{`10 print "abcd";pos(0);`, "", "abcd 5 ", OK},
	{`10 print "abcd";chr$(13);pos(0);`, "", "abcd\n 1 ", OK},
	{`10 print "abcd";:print:print pos(0);`, "", "abcd\n 1 ", OK},

	{sourceDefint, "", " 1.1 b 2 ", OK},

	{`10 a$=input$(3):print a$;`, "abcde\n", "abc", OK},

	{`10 input a:print a`, "2\n", "?  2 \n", OK},

	{sourceOnGoto, "", "x", OK},
	{sourceOnGosub2, "", "a 1 2 3b 4", OK},

	{"10 data\n20 read a$:print a$;", "", "", OK},
	{"10 data ,\n20 read a$:print a$;", "", "", OK},
	{`10 data:read a$:print a$;`, "", "", OK},
	{`10 data ,:read a$:print a$;`, "", "", OK},
	{`10 data a,b:read a$,b$:print a$,b$;`, "", "ab", OK},
	{`10 data ,b:read a$,b$:print a$,b$;`, "", "b", OK},
	{`10 data a,,b:read a$,b$:print a$,b$;`, "", "a", OK},
	{`10 data a,&hb:read a$,b:print a$,b;`, "", "a 11 ", OK},
	{`10 data a, 3:read a$,b:print a$,b;`, "", "a 3 ", OK},
	{`10 data a, .1:read a$,b:print a$,b;`, "", "a 0.1 ", OK},

	{sourceRestore, "", " 1  7  3 ", OK},

	{`10 rem a`, "", "", OK},
	{`10 ' a`, "", "", OK},
	{`10 print "1";:rem a`, "", "1", OK},
	{`10 print "1";:' a`, "", "1", OK},
	{`10 print "1";rem a`, "", "", WRONG},
	{`10 print "1";' a`, "", "1", OK},

	{`10 print ""`, "", "\n", OK},
	{`10 print "`, "", "\n", OK},
	{`10 print "hello`, "", "hello\n", OK},
	{`10 print "hello  `, "", "hello  \n", OK},

	{`10 _goimport("fmt"):a$=_gofunc("fmt.Sprintf$","gofunc-good"): print a$;`, "", "gofunc-good", OK},
	{`10 _goimport("fmt"):print "goproc-";:_goproc("fmt.Print","good")`, "", "goproc-good", OK},
	{sourceGofunc, "180", "?  180 degrees in radians is 3.141592653589793 ", OK},

	{`10 print instr("abcabcabc","");`, "", " 1 ", OK},
	{`10 print instr("abcabcabc","bc");`, "", " 2 ", OK},
	{`10 print instr(1,"abcabcabc","bc");`, "", " 2 ", OK},
	{`10 print instr(2,"abcabcabc","bc");`, "", " 2 ", OK},
	{`10 print instr(3,"abcabcabc","bc");`, "", " 5 ", OK},
	{`10 print instr(3,"abcabcabc","");`, "", " 3 ", OK},
	{`10 print instr(3,"abcabcabc","xy");`, "", " 0 ", OK},

	{sourceReturnLine, "", "143", OK},

	{`10 print 2!`, "", " 2 \n", OK},
	{`10 print 2!;`, "", " 2 ", OK},
	{`10 a=1:print a=1!`, "", " -1 \n", OK},
	{`10 a=1:print a=1!;`, "", " -1 ", OK},

	{`10 print &`, "", " 0 \n", WRONG},
	{`10 print &h`, "", " 0 \n", WRONG},
	{`10 print &hg`, "", " 0 \n", WRONG},
	{`10 print &hb`, "", " 11 \n", OK},
	{"10 print &hb", "", " 11 \n", OK},
	{`10 print &hb;`, "", " 11 ", OK},

	{"10 print 0 => 0", "", " -1 \n", OK},
	{"10 print 1 => 0", "", " -1 \n", OK},
	{"10 print 0 => 1", "", " 0 \n", OK},
	{"10 print 0 =< 0", "", " -1 \n", OK},
	{"10 print 1 =< 0", "", " 0 \n", OK},
	{"10 print 0 =< 1", "", " -1 \n", OK},
	{"10 print 0 >< 0", "", " 0 \n", OK},
	{"10 print 0 >< 1", "", " -1 \n", OK},

	{`10 print sin(0);`, "", " 0 ", OK},
	{`10 print cos(0);`, "", " 1 ", OK},
	{`10 print tan(0);`, "", " 0 ", OK},
	{`10 print sqr(9);`, "", " 3 ", OK},
	{`10 print sqr(-1);`, "", " NaN ", OK},

	{`10 while a$<>"b":a$=inkey$:print a$;:wend`, "abcd", "ab", OK},

	{sourceOnGosub, "", "1push23end", OK},

	{sourceGoto2, "", " 3 \n", OK},

	{`10 def fa() = 1:print fa()`, "", "", WRONG},
	{`10 def fn() = 1:print fn()`, "", "", WRONG},
	{`10 def fn2() = 1:print fn2()`, "", "", WRONG},
	{`10 print fa(1)`, "", " 0 \n", OK},
	{`10 print fn(1)`, "", " 0 \n", OK},
	{`10 print fn2(1)`, "", " 0 \n", OK},
	{`10 print fna(1)`, "", "", WRONG},
	{`10 def fna() = 3:print fna();`, "", " 3 ", OK},
	{`10 def fna2() = 3:print fna2();`, "", " 3 ", OK},
	{`10 def fna$(a$) = a$+a$:print fna$("1");`, "", "11", OK},
	{`10 def fna(a) = a+a:print fna(1);`, "", " 2 ", OK},
	{`10 b=2:def fna() = b+b:print fna();`, "", " 4 ", OK},
	{`10 b=2:def fna(a) = a+b:print fna(3);`, "", " 5 ", OK},
	{`10 b=2:def fna(a,b) = a+b:print fna(5,5);:print b;`, "", " 10  2 ", OK},
	{`10 b=2:c=3:def fna(a,b) = a+b+c:print fna(5,5);:print b;`, "", " 13  2 ", OK},
	{`10 b=2:def fna() = int(b): print fna();`, "", " 2 ", OK},
	{`10 b=2:def fna%() = b: print fna%();`, "", " 2 ", OK},
	{`10 def fna(b,b) = b:print fna(4,5)`, "", "", WRONG},

	{`10 print len(time$);`, "", " 8 ", OK},
	{`10 print len(date$);`, "", " 10 ", OK},
	{`10 t=timer:print (t>=0) and (t<86400);`, "", " -1 ", OK},

	{`10 print sgn`, "", "", WRONG},
	{`10 print sgn()`, "", "", WRONG},
	{`10 print sgn("")`, "", "", WRONG},
	{`10 print sgn(2);`, "", " 1 ", OK},
	{`10 print sgn(0);`, "", " 0 ", OK},
	{`10 print sgn(-3);`, "", " -1 ", OK},
	{`10 print sgn(1.2);`, "", " 1 ", OK},
	{`10 print sgn(-1.2);`, "", " -1 ", OK},

	{`10 print abs`, "", "", WRONG},
	{`10 print abs()`, "", "", WRONG},
	{`10 print abs("")`, "", "", WRONG},
	{`10 print abs(1);`, "", " 1 ", OK},
	{`10 print abs(0);`, "", " 0 ", OK},
	{`10 print abs(-1);`, "", " 1 ", OK},
	{`10 print abs(1.2);`, "", " 1.2 ", OK},
	{`10 print abs(-1.2);`, "", " 1.2 ", OK},

	{`10 print 0e`, "", " 0 \n", OK},
	{`10 print 1e`, "", " 1 \n", OK},
	{`10 print 1e+`, "", " 1 \n", WRONG},
	{`10 print 1e-`, "", " 1 \n", WRONG},
	{`10 print .1e`, "", " 0.1 \n", OK},
	{`10 print 1e2`, "", " 100 \n", OK},
	{`10 print .12345e+5`, "", " 12345 \n", OK},
	{`10 print 12.34e56`, "", " 1.234e+57 \n", OK},
	{`10 print 12.34e+56`, "", " 1.234e+57 \n", OK},
	{`10 print 12.34e-56`, "", " 1.234e-55 \n", OK},

	{`10 print 0e;`, "", " 0 ", OK},
	{`10 print 1e;`, "", " 1 ", OK},
	{`10 print 1e+;`, "", " 1 ", WRONG},
	{`10 print 1e-;`, "", " 1 ", WRONG},
	{`10 print .1e;`, "", " 0.1 ", OK},
	{`10 print 1e2;`, "", " 100 ", OK},
	{`10 print .12345e+5;`, "", " 12345 ", OK},
	{`10 print 12.34e56;`, "", " 1.234e+57 ", OK},
	{`10 print 12.34e+56;`, "", " 1.234e+57 ", OK},
	{`10 print 12.34e-56;`, "", " 1.234e-55 ", OK},

	{`10 a=12.34e56:print a;`, "", " 1.234e+57 ", OK},
	{`10 a=12.34ee56:print a;`, "", "", WRONG},
	{`10 a=12.34e+-56:print a;`, "", "", WRONG},
	{`10 a=12.34e++56:print a;`, "", "", WRONG},
	{`10 a=12.34e--56:print a;`, "", "", WRONG},

	{`10 print int(98.89);`, "", " 98 ", OK},
	{`10 print int(-12.11);`, "", " -13 ", OK},
	{`10 print a-int(1.1);`, "", " -1 ", OK},
	{`10 print a+int(1.1);`, "", " 1 ", OK},
	{`10 print a*int(1.1);`, "", " 0 ", OK},
	{`10 print a/int(1.1);`, "", " 0 ", OK},
	{`10 print int(1.1)-a;`, "", " 1 ", OK},
	{`10 print int(1.1)+a;`, "", " 1 ", OK},
	{`10 print int(1.1)*a;`, "", " 0 ", OK},
	{`10 a=1:print int(1.1)/a;`, "", " 1 ", OK},
	{`10 print int(5/2)*2+1;`, "", " 5 ", OK},

	{"10 print 1+2", "", " 3 \n", OK},
	{"10 print 1.1+2", "", " 3.1 \n", OK},
	{"10 print 1.1+2.2", "", " 3.3 \n", OK},
	{`10 print "a"+"b"`, "", "ab\n", OK},
	{`10 print 1+"b"`, "", "", WRONG},

	{"10 print 1-2", "", " -1 \n", OK},
	{"10 print 1.2-2", "", " -0.8 \n", OK},
	{"10 print 1.1-2.2", "", " -1.1 \n", OK},
	{`10 print "a"-"b"`, "", "", WRONG},

	{"10 print 5 MOD 3", "", " 2 \n", OK},
	{"10 print 5.5 MOD 3.3", "", " 0 \n", OK},
	{`10 print "a" MOD "b"`, "", "", WRONG},

	{`10 print 5 \ 3`, "", " 1 \n", OK},
	{`10 print 5 \ 2`, "", " 2 \n", OK},
	{`10 print 5 \ 2.5`, "", " 1 \n", OK},
	{`10 print 6.6 \ 3.3`, "", " 2 \n", OK},
	{`10 print "a" \ "b"`, "", "", WRONG},

	{`10 print 5 * 3`, "", " 15 \n", OK},
	{`10 print 1.1 * 2`, "", " 2.2 \n", OK},
	{`10 print 2 * 2.5`, "", " 5 \n", OK},
	{`10 print "a" * "b"`, "", "", WRONG},

	{`10 print 5 / 5`, "", " 1 \n", OK},
	{`10 print 5 / 4`, "", " 1.25 \n", OK},
	{`10 print 5 / 2`, "", " 2.5 \n", OK},
	{`10 print 5 / 1`, "", " 5 \n", OK},
	{`10 print 5 / 2.5`, "", " 2 \n", OK},
	{`10 print 6.6 / 3.3`, "", " 2 \n", OK},
	{`10 print "a" / "b"`, "", "", WRONG},

	{`10 print 2 ^ 3`, "", " 8 \n", OK},
	{`10 print 16 ^ .5`, "", " 4 \n", OK},
	{`10 print "a" ^ "b"`, "", "", WRONG},

	{`10 print +10`, "", " 10 \n", OK},
	{`10 print +(2.1-1)`, "", " 1.1 \n", OK},
	{`10 print +"a"`, "", "", WRONG},

	{`10 print -10`, "", " -10 \n", OK},
	{`10 print -(2.1-1)`, "", " -1.1 \n", OK},
	{`10 print -"a"`, "", "", WRONG},

	{`10 print (22)`, "", " 22 \n", OK},
	{`10 print ("a"+"b")`, "", "ab\n", OK},
	{`10 print 2*(3+4)`, "", " 14 \n", OK},

	{`10 print LEN("hello")`, "", " 5 \n", OK},
	{`10 print LEN(2)`, "", " 8 \n", OK},
	{`10 print LEN(3.3)`, "", " 8 \n", OK},

	{`10 print not 0`, "", " -1 \n", OK},
	{`10 print not -1`, "", " 0 \n", OK},
	{`10 print not 1.1`, "", " -2 \n", OK},
	{`10 print not ""`, "", "", WRONG},

	{"10 print -1 and -1", "", " -1 \n", OK},
	{"10 print 0 and 0", "", " 0 \n", OK},
	{"10 print 1 and 0", "", " 0 \n", OK},
	{"10 print 1 and 3", "", " 1 \n", OK},
	{`10 print "" and ""`, "", "", WRONG},

	{"10 print -1 or -1", "", " -1 \n", OK},
	{"10 print 0 or 0", "", " 0 \n", OK},
	{"10 print 1 or 0", "", " 1 \n", OK},
	{"10 print 1 or 3", "", " 3 \n", OK},
	{`10 print "" or ""`, "", "", WRONG},

	{"10 print -1 xor -1", "", " 0 \n", OK},
	{"10 print 0 xor -1", "", " -1 \n", OK},
	{"10 print 0 xor 0", "", " 0 \n", OK},
	{"10 print 1 xor 0", "", " 1 \n", OK},
	{"10 print 1 xor 3", "", " 2 \n", OK},
	{`10 print "" xor ""`, "", "", WRONG},

	{"10 print -1 eqv -1", "", " -1 \n", OK},
	{"10 print 0 eqv -1", "", " 0 \n", OK},
	{"10 print 0 eqv 0", "", " -1 \n", OK},
	{`10 print "" eqv ""`, "", "", WRONG},

	{"10 print -1 imp -1", "", " -1 \n", OK},
	{"10 print 0 imp -1", "", " -1 \n", OK},
	{"10 print 0 imp 0", "", " -1 \n", OK},
	{"10 print -1 imp 0", "", " 0 \n", OK},
	{`10 print "" imp ""`, "", "", WRONG},

	{"10 print 0 = 0", "", " -1 \n", OK},
	{"10 print 1.1 = 1", "", " 0 \n", OK},
	{"10 print 2.2 = 2.2", "", " -1 \n", OK},
	{`10 print "" = ""`, "", " -1 \n", OK},
	{"10 print 0 = 1", "", " 0 \n", OK},
	{"10 print 1.1 = 2", "", " 0 \n", OK},
	{"10 print 2.2 = 3.3", "", " 0 \n", OK},
	{`10 print "a" = ""`, "", " 0 \n", OK},
	{`10 print 0 = ""`, "", "", WRONG},

	{"10 print 0 <> 0", "", " 0 \n", OK},
	{"10 print 1.1 <> 1", "", " -1 \n", OK},
	{"10 print 2.2 <> 2.2", "", " 0 \n", OK},
	{`10 print "" <> ""`, "", " 0 \n", OK},
	{"10 print 0 <> 1", "", " -1 \n", OK},
	{"10 print 1.1 <> 2", "", " -1 \n", OK},
	{"10 print 2.2 <> 3.3", "", " -1 \n", OK},
	{`10 print "a" <> ""`, "", " -1 \n", OK},
	{`10 print 0 <> ""`, "", "", WRONG},

	{"10 print 0 > 0", "", " 0 \n", OK},
	{"10 print 1.1 > 1", "", " -1 \n", OK},
	{"10 print 2.2 > 2.2", "", " 0 \n", OK},
	{`10 print "" > ""`, "", " 0 \n", OK},
	{`10 print "a" > "b"`, "", " 0 \n", OK},
	{"10 print 0 > 1", "", " 0 \n", OK},
	{"10 print 1 > 0", "", " -1 \n", OK},
	{"10 print 1.1 > 2", "", " 0 \n", OK},
	{"10 print 2.2 > 3.3", "", " 0 \n", OK},
	{"10 print 3.3 > 2.2", "", " -1 \n", OK},
	{`10 print "a" > ""`, "", " -1 \n", OK},
	{`10 print 0 > ""`, "", "", WRONG},

	{"10 print 0 < 0", "", " 0 \n", OK},
	{"10 print 1.1 < 1", "", " 0 \n", OK},
	{"10 print 2.2 < 2.2", "", " 0 \n", OK},
	{`10 print "" < ""`, "", " 0 \n", OK},
	{`10 print "a" < "b"`, "", " -1 \n", OK},
	{"10 print 0 < 1", "", " -1 \n", OK},
	{"10 print 1 < 0", "", " 0 \n", OK},
	{"10 print 1.1 < 2", "", " -1 \n", OK},
	{"10 print 2.2 < 3.3", "", " -1 \n", OK},
	{"10 print 3.3 < 2.2", "", " 0 \n", OK},
	{`10 print "a" < ""`, "", " 0 \n", OK},
	{`10 print 0 < ""`, "", "", WRONG},

	{"10 print 0 >= 0", "", " -1 \n", OK},
	{"10 print 1.1 >= 1", "", " -1 \n", OK},
	{"10 print 2.2 >= 2.2", "", " -1 \n", OK},
	{`10 print "" >= ""`, "", " -1 \n", OK},
	{`10 print "a" >= "b"`, "", " 0 \n", OK},
	{"10 print 0 >= 1", "", " 0 \n", OK},
	{"10 print 1 >= 0", "", " -1 \n", OK},
	{"10 print 1.1 >= 2", "", " 0 \n", OK},
	{"10 print 2.2 >= 3.3", "", " 0 \n", OK},
	{"10 print 3.3 >= 2.2", "", " -1 \n", OK},
	{`10 print "a" >= ""`, "", " -1 \n", OK},
	{`10 print 0 >= ""`, "", "", WRONG},

	{"10 print 0 <= 0", "", " -1 \n", OK},
	{"10 print 1.1 <= 1", "", " 0 \n", OK},
	{"10 print 2.2 <= 2.2", "", " -1 \n", OK},
	{`10 print "" <= ""`, "", " -1 \n", OK},
	{`10 print "a" <= "b"`, "", " -1 \n", OK},
	{"10 print 0 <= 1", "", " -1 \n", OK},
	{"10 print 1 <= 0", "", " 0 \n", OK},
	{"10 print 1.1 <= 2", "", " -1 \n", OK},
	{"10 print 2.2 <= 3.3", "", " -1 \n", OK},
	{"10 print 3.3 <= 2.2", "", " 0 \n", OK},
	{`10 print "a" <= ""`, "", " 0 \n", OK},
	{`10 print 0 <= ""`, "", "", WRONG},

	{`10 print a!<a%;`, "", " 0 ", OK},
	{`10 print a!<=a%;`, "", " -1 ", OK},
	{`10 print a!=a%;`, "", " -1 ", OK},
	{`10 print a!>=a%;`, "", " -1 ", OK},
	{`10 print a!<>a%;`, "", " 0 ", OK},

	{`10 goto 20`, "", "", WRONG},
	{sourceGoto, "", outputGoto, OK},

	{`10 print "hi"`, "", "hi\n", OK},
	{`10 print "hi";`, "", "hi", OK},

	{`10 if "" then print 1`, "", "", WRONG},
	{`10 if 0 then print 1`, "", "", OK},
	{`10 if -1 then print 1`, "", " 1 \n", OK},
	{`10 if 1.1 then print 1`, "", " 1 \n", OK},
	{`10 if 0 then print 1 else print 2`, "", " 2 \n", OK},
	{`10 if -1 then print 1 else print 2`, "", " 1 \n", OK},
	{`10 if 1.1 then print 1 else print 2`, "", " 1 \n", OK},

	{"10 if 0 then 20\n20 end\n30 print 30", "", "", OK},
	{"10 if 0 goto 20\n20 end\n30 print 30", "", "", OK},
	{"10 if 0 then 20 else 30\n20 end\n30 print 30", "", " 30 \n", OK},
	{"10 if 0 goto 20 else 30\n20 end\n30 print 30", "", " 30 \n", OK},

	{`10 if 1 then print 1:print 2:print 3 else print 4:print 5:print 6`, "", " 1 \n 2 \n 3 \n", OK},
	{`10 if 0 then print 1:print 2:print 3 else print 4:print 5:print 6`, "", " 4 \n 5 \n 6 \n", OK},

	{`10 print "abc" : if 1 then print 1:print 2:print 3 else print 4:print 5:print 6`, "", "abc\n 1 \n 2 \n 3 \n", OK},
	{`10 print "abc" : if 0 then print 1:print 2:print 3 else print 4:print 5:print 6`, "", "abc\n 4 \n 5 \n 6 \n", OK},

	{`10 for a=1 to 2:if a>1 then next`, "", "", OK},

	{`10 line input a:print a`, "", "", WRONG},
	{`10 line input a$:print a$`, "2\n", "2\n", OK},
	{`10 line input "choice? ",a$:print a$`, "2\n", "choice? 2\n", OK},
	{`10 line input "choice? ";a$:print a$`, "2\n", "choice? 2\n", OK},

	{`10 input a : print a`, "2\n", "?  2 \n", OK},
	{`10 input a(1) : print a(1)`, "2\n", "?  2 \n", OK},
	{`10 input "",a : print a`, "2\n", " 2 \n", OK},
	{`10 input a! : print a!`, "2.1\n", "?  2.1 \n", OK},
	{`10 input a# : print a#`, "2.1\n", "?  2.1 \n", OK},
	{`10 input a% : print a%`, "2\n", "?  2 \n", OK},
	{`10 input a$ : print a$`, "abc\n", "? abc\n", OK},
	{`10 input a:input b:print a" "b;`, "2\n3\n", "? ?  2   3 ", OK},
	{`10 input a,b:print a" "b;`, "2,3\n", "?  2   3 ", OK},

	{`10 a="":print a`, "", "", WRONG},
	{`10 a%="":print a%`, "", "", WRONG},
	{`10 a!="":print a!`, "", "", WRONG},
	{`10 a#="":print a#`, "", "", WRONG},
	{`10 a$="":print a$`, "", "\n", OK},
	{`10 a$=1:print a$`, "", "", WRONG},

	{`10 a=print:print a`, "", "", WRONG},
	{`10 a=let:print a`, "", "", WRONG},
	{`10 print=1`, "", "", WRONG},
	{`10 let=1`, "", "", WRONG},
	{`10 rnd=1`, "", "", WRONG},

	{`10 a=a%:print a`, "", " 0 \n", OK},
	{`10 a=a#:print a`, "", " 0 \n", OK},
	{`10 a%=a:print a%`, "", " 0 \n", OK},
	{`10 a%=a#:print a%`, "", " 0 \n", OK},
	{`10 a#=a:print a#`, "", " 0 \n", OK},
	{`10 a#=a%:print a#`, "", " 0 \n", OK},

	{`10 a%=5.1:print a%`, "", " 5 \n", OK},
	{`10 a%=5.9:print a%`, "", " 6 \n", OK},

	{`10 a=45:print A`, "", " 45 \n", OK},
	{`10 A=45:print a`, "", " 45 \n", OK},
	{`10 input a:print A;`, "23\n", "?  23 ", OK},
	{`10 input A:print a;`, "23\n", "?  23 ", OK},

	{`10 print int(2)`, "", " 2 \n", OK},
	{`10 print int(2.1)`, "", " 2 \n", OK},
	{`10 print int(2.9)`, "", " 2 \n", OK},
	{`10 print int("")`, "", "", WRONG},

	{`10 print left$("abc",-1)`, "", "\n", OK},
	{`10 print left$("abc",0)`, "", "\n", OK},
	{`10 print left$("abc",1)`, "", "a\n", OK},
	{`10 print left$("abc",2)`, "", "ab\n", OK},
	{`10 print left$("abc",3)`, "", "abc\n", OK},
	{`10 print left$("abc",4)`, "", "abc\n", OK},
	{`10 print left$("abc",1.1)`, "", "a\n", OK},
	{`10 print left$("abc",.6+.5)`, "", "a\n", OK},
	{`10 print left$("abc",2.9)`, "", "abc\n", OK},
	{`10 print left$("abc","")`, "", "", WRONG},
	{`10 print left$("abc",a$)`, "", "", WRONG},
	{`10 print left$(1,1)`, "", "", WRONG},

	{`10 print 1:end:print 2`, "", " 1 \n", OK},
	{`10 print 1:stop:print 2`, "", " 1 \n", OK},

	{"10 on 1 goto 20\n20 rem", "", "", OK},
	{"10 on 2 goto 10,20\n20 rem", "", "", OK},
	{`10 on 1 goto x`, "", "", WRONG},
	{`10 on 1 goto 20`, "", "", WRONG},
	{`10 on "" goto 10`, "", "", WRONG},
	{`10 on a$ goto 10`, "", "", WRONG},
	{"10 on 1 goto 20,30,40\n20 print 20:end\n30 print 30:end\n40 print 40", "", " 20 \n", OK},
	{"10 on 2 goto 20,30,40\n20 print 20:end\n30 print 30:end\n40 print 40", "", " 30 \n", OK},
	{"10 on 3 goto 20,30,40\n20 print 20:end\n30 print 30:end\n40 print 40", "", " 40 \n", OK},
	{"10 on 4 goto 20,30,40\n20 print 20:end\n30 print 30:end\n40 print 40", "", " 20 \n", OK},
	{"10 a=1:on a goto 30\n20 print 20:end\n30 print 30", "", " 30 \n", OK},
	{"10 a=1.1:on a goto 30\n20 print 20:end\n30 print 30", "", " 30 \n", OK},

	{"10 for a=1 to 3:print a;:next", "", " 1  2  3 ", OK},
	{"10 for a=1 to 3:print a;:next a", "", " 1  2  3 ", OK},
	{"10 for a=1 to 3:print a;:next b", "", "", WRONG},
	{"10 for a=4 to 0 step -2:print a;:next a", "", " 4  2  0 ", OK},
	{"10 for a=1 to 3 step 2:print a;:next a", "", " 1  3 ", OK},
	{"10 for a%=1 to 3 step 2:print a%;:next", "", " 1  3 ", OK},
	{"10 b%=1:for a=b% to 3 step 2:print a;:next", "", " 1  3 ", OK},
	{"10 b%=1:c%=3:for a=b% to c% step 2:print a;:next", "", " 1  3 ", OK},
	{"10 b%=1:c%=3:d%=2:for a=b% to c% step d%:print a;:next", "", " 1  3 ", OK},
	{"10 b=1:c=3:d=2:for a%=b to c step d:print a%;:next", "", " 1  3 ", OK},
	{"10 for a=1 to 3:for b=4 to 5:print str$(a),str$(b);:next:next", "", " 1 4 1 5 2 4 2 5 3 4 3 5", OK},
	{"10 for a=1 to 3:for b=4 to 5:print str$(a),str$(b);:next b,a", "", " 1 4 1 5 2 4 2 5 3 4 3 5", OK},
	{"10 for a=1 to 3:for b=4 to 5:print a,b;:next a,b", "", "", WRONG},
	{"10 for a=1 to 3:for b=4 to 5:print a,b;:next", "", "", WRONG},
	{"10 for a=1 to 3 step -1:print a;:next", "", "", OK},
	{"10 for a=3 to 1 step 1:print a;:next", "", "", OK},

	{`10 data`, "", "", OK},
	{`10 read`, "", "", WRONG},
	{`10 data a:print a;`, "", " 0 ", OK},
	{`10 read a:print a;`, "", "", RUNTIME},
	{`10 data 3:read 3`, "", "", WRONG},
	{`10 data a:read a:print a;`, "", "", RUNTIME},
	{`10 data 3:read a:print a;`, "", " 3 ", OK},
	{`10 data 3:read a,b:print a,b;`, "", "", RUNTIME},
	{`10 data "3":read a$:print a$;`, "", "3", OK},
	{`10 data 3:read a$:print a$;`, "", "", RUNTIME},
	{`10 data 3,"4":read a,b$:print a,b$;`, "", " 3 4", OK},
	{`10 data 3,"4":read a:print a;`, "", " 3 ", OK},
	{`10 data 3,"4",5,"6":read a,b$,c,d$:print a,b$,c,d$;`, "", " 3 4 5 6", OK},
	{`10 data 1,-2,3,-4:read a,b,c,d:print a,b,c,d;`, "", " 1  -2  3  -4 ", OK},
	{`10 data +1.1,-2.1,+3.1,-4.1:read a,b,c,d:print a,b,c,d;`, "", " 1.1  -2.1  3.1  -4.1 ", OK},

	{`10 print a()`, "", "", WRONG},
	{`10 print a("")`, "", "", WRONG},
	{`10 print a(1)`, "", " 0 \n", OK},
	{`10 print a(1,2)`, "", " 0 \n", OK},
	{`10 print a$(1,2)`, "", "\n", OK},
	{`10 i=1.1:print a(i,2)`, "", " 0 \n", OK},
	{`10 print a(11,2)`, "", "", RUNTIME},

	{`10 a()=2:print a()`, "", "", WRONG},
	{`10 a("")=2:print a("")`, "", "", WRONG},
	{`10 a(1)=2:print a(1)`, "", " 2 \n", OK},
	{`10 a(1,2)=3:print a(1,2)`, "", " 3 \n", OK},
	{`10 x=1.2:i=1.1:a(i,2)=x:print a(i,2)`, "", " 1.2 \n", OK},
	{`10 a(11)=2:print a(11)`, "", "", RUNTIME},
	{`10 a(1)=2:print a(1,1)`, "", "", WRONG},

	{`10 data 2:read a(1):print a(1)`, "", " 2 \n", OK},
	{`10 data "3":read a$(1,2):print a$(1,2)`, "", "3\n", OK},
	{`10 data 2,3:read a(1),b:print a(1),b`, "", " 2  3 \n", OK},

	{`10 for a(1)=1 to 3:print str$(a(1));:next a(1)`, "", " 1 2 3", OK},
	{`10 for a(1)=1 to 3:print str$(a(1));:next`, "", " 1 2 3", OK},
	{`10 for a(1,2)=1 to 3:print str$(a(1,2));:next a(1,2)`, "", " 1 2 3", OK},
	{`10 for a(1,2)=1 to 3:print str$(a(1,2));:next`, "", " 1 2 3", OK},
	{`10 for a(1)=1 to 3:print a(1);:next a(2)`, "", "", WRONG},
	{`10 for a(1,2)=1 to 3:print a(1,2);:next a(1,1)`, "", "", WRONG},

	{`10 dim a():print a()`, "", "", WRONG},
	{`10 dim a(""):print a("")`, "", "", WRONG},
	{`10 dim a(b):a(b)=2:print a(b)`, "", "", WRONG},
	{`10 dim a(-1):a(-1)=2:print a(-1)`, "", "", WRONG},
	{`10 dim a(0):a(0)=2:print a(0)`, "", " 2 \n", OK},
	{`10 dim a(1):a(0)=2:print a(0)`, "", " 2 \n", OK},
	{`10 dim a(1):a(1)=2:print a(1)`, "", " 2 \n", OK},
	{`10 dim a(1),b(2):a(1)=2:print a(1),b(2)`, "", " 2  0 \n", OK},
	{`10 dim a(20):a(20)=2:print a(20)`, "", " 2 \n", OK},
	{`10 dim a(20):a(21)=2:print a(21)`, "", "", RUNTIME},
	{`10 dim a(20,30):a(20,30)=2:print a(20,30)`, "", " 2 \n", OK},
	{`10 dim a(1):a(1,1)=2:print a(1):a(1,1)`, "", "", WRONG},
	{`10 a(1)=2:dim a(1):print a(1)`, "", " 0 \n", OK},
	{`10 a(1)=2:dim a(2):print a(1)`, "", " 0 \n", OK},
	{`10 a(1)=2:dim a(1,1):print a(1),a(1,1)`, "", "", WRONG},
	{`10 dim a(1):dim a(1):a(1)=2:print a(1)`, "", " 2 \n", OK},
	{`10 dim a(1):dim a(2):a(1)=2:print a(1)`, "", "", WRONG},
	{`10 dim a(1,1):dim a(1,1):a(1,1)=2:print a(1,1)`, "", " 2 \n", OK},
	{`10 a(1,1)=2:dim a(1,1):print a(1,1)`, "", " 0 \n", OK},

	{`10 restore`, "", "", RUNTIME},
	{`10 data 2,3:read a:print a:read a:print a`, "", " 2 \n 3 \n", OK},
	{`10 data 2,3:read a:print a:restore:read a:print a`, "", " 2 \n 2 \n", OK},

	{`10 print mid$(1,1,1)`, "", "", WRONG},
	{`10 print mid$("abc","",1)`, "", "", WRONG},
	{`10 print mid$("abc",1,"")`, "", "", WRONG},
	{`10 print mid$("abc",1,1);`, "", "a", OK},
	{`10 print mid$("abc",1,2);`, "", "ab", OK},
	{`10 print mid$("abc",1,3);`, "", "abc", OK},
	{`10 print mid$("abc",2,1);`, "", "b", OK},
	{`10 print mid$("abc",2,2);`, "", "bc", OK},
	{`10 print mid$("abc",3,1);`, "", "c", OK},
	{`10 print mid$("abc",1,0);`, "", "", OK},
	{`10 print mid$("abc",2,0);`, "", "", OK},
	{`10 print mid$("abc",3,0);`, "", "", OK},
	{`10 print mid$("abc",4,1);`, "", "", OK},
	{`10 print mid$("abc",0,1);`, "", "a", OK},

	{`10 print mid$("abc",0);`, "", "abc", OK},
	{`10 print mid$("abc",1);`, "", "abc", OK},
	{`10 print mid$("abc",2);`, "", "bc", OK},
	{`10 print mid$("abc",3);`, "", "c", OK},
	{`10 print mid$("abc",4);`, "", "", OK},

	{`10 gosub 20`, "", "", WRONG},
	{`10 return`, "", "", RUNTIME},
	{"10 gosub 20\n20 end", "", "", OK},
	{"10 gosub 20\n20 print 3", "", " 3 \n", OK},
	{"10 print 1;:gosub 20:print 2;:end\n20 print 3;:return", "", " 1  3  2 ", OK},
	{sourceGosub, "", "1234567", OK},

	{`10 print "1";:while a<3:print str$(a);:a=a+1:wend`, "", "1 0 1 2", OK},
	{`10 print "1";:while 0:print a;:a=a+1:wend`, "", "1", OK},
	{`10 print "1";:while "":print a;:a=a+1:wend`, "", "", WRONG},
	{`10 print "1";:while a<3:print a;:a=a+1`, "", "", WRONG},
	{`10 while 0`, "", "", WRONG},
	{`10 wend`, "", "", WRONG},

	{`10 swap`, "", "", WRONG},
	{`10 swap a:print a`, "", "", WRONG},
	{`10 swap a,b$:print a,b$`, "", "", WRONG},
	{`10 swap a,b%:print a,b%`, "", "", WRONG},
	{`10 a=1:b=2:swap a,b:print a,b;`, "", " 2  1 ", OK},
	{`10 a=1.1:b=2.2:swap a,b:print a,b;`, "", " 2.2  1.1 ", OK},
	{`10 a$="a":b$="b":swap a$,b$:print a$,b$;`, "", "ba", OK},

	{`10 print str$();`, "", "", WRONG},
	{`10 print str$("");`, "", "", WRONG},
	{`10 print str$(1);`, "", " 1", OK},
	{`10 print str$(1.1);`, "", " 1.1", OK},
	{`10 print str$(1+1);`, "", " 2", OK},
	{`10 print str$(1+.1);`, "", " 1.1", OK},
	{`10 print "1"+str$(1+.1);`, "", "1 1.1", OK},

	{`10 print val();`, "", "", WRONG},
	{`10 print val(0);`, "", "", WRONG},
	{`10 print val("");`, "", " 0 ", OK},
	{`10 print val("1");`, "", " 1 ", OK},
	{`10 print val("1.1");`, "", " 1.1 ", OK},
	{`10 print val("1"+"1");`, "", " 11 ", OK},
	{`10 print val("1"+".1");`, "", " 1.1 ", OK},
	{`10 print 1+val("1"+".1");`, "", " 2.1 ", OK},

	{`10 print right$("abc",-1)`, "", "\n", OK},
	{`10 print right$("abc",0)`, "", "\n", OK},
	{`10 print right$("abc",1)`, "", "c\n", OK},
	{`10 print right$("abc",2)`, "", "bc\n", OK},
	{`10 print right$("abc",3)`, "", "abc\n", OK},
	{`10 print right$("abc",4)`, "", "abc\n", OK},
	{`10 print right$("abc",1.1)`, "", "c\n", OK},
	{`10 print right$("abc",.6+.5)`, "", "c\n", OK},
	{`10 print right$("abc",2.9)`, "", "abc\n", OK},
	{`10 print right$("abc","")`, "", "", WRONG},
	{`10 print right$("abc",a$)`, "", "", WRONG},
	{`10 print right$(1,1)`, "", "", WRONG},

	{`10 print "a";tab(0);"b";`, "", "a\nb", OK},
	{`10 print "a";tab(1);"b";`, "", "a\nb", OK},
	{`10 print "a";tab(2);"b";`, "", "ab", OK},
	{`10 print "a";tab(3);"b";`, "", "a b", OK},
	{`10 print "a";tab(4);"b";`, "", "a  b", OK},

	{`10 print "a"+spc(0)+"b";`, "", "ab", OK},
	{`10 print "a"+spc(1)+"b";`, "", "a b", OK},
	{`10 print "a"+spc(2)+"b";`, "", "a  b", OK},

	{`10 print "a"+space$(0)+"b";`, "", "ab", OK},
	{`10 print "a"+space$(1)+"b";`, "", "a b", OK},
	{`10 print "a"+space$(2)+"b";`, "", "a  b", OK},

	{`10 print "a"+string$(0,"1")+"b";`, "", "ab", OK},
	{`10 print "a"+string$(1,"")+"b";`, "", "ab", OK},
	{`10 print "a"+string$(1,"1")+"b";`, "", "a1b", OK},
	{`10 print "a"+string$(2,"1")+"b";`, "", "a11b", OK},
	{`10 print "a"+string$(2,"21")+"b";`, "", "a22b", OK},
	{`10 print "a"+string$(2,32)+"b";`, "", "a  b", OK},

	{`10 print chr$(32);`, "", " ", OK},
	{`10 print chr$(0);`, "", string(0), OK},

	{`10 print asc(" ");`, "", " 32 ", OK},
	{`10 print asc(" a");`, "", " 32 ", OK},
	{`10 print asc("");`, "", " 0 ", OK},

	{`10 print chr$(asc("a"));`, "", "a", OK},
	{`10 print asc(chr$(32));`, "", " 32 ", OK},
	{`10 print chr$(asc("a"))="a";`, "", " -1 ", OK},
	{`10 print asc(chr$(32))=32;`, "", " -1 ", OK},
}

const sourceGosub = `
100 print "1";
110 gosub 200
120 print "7";
130 end

200 print "2";
210 gosub 300
220 print "6";
230 return

300 print "3";
310 gosub 400
320 print "5";
330 return

400 print "4";
410 return
`

func TestBuild(t *testing.T) {

	for _, data := range testTable {

		t.Logf("source: %q\n", data.source)

		tmp := "/tmp/basgo-test-build.go"
		w, errCreate := os.Create(tmp)
		if errCreate != nil {
			t.Errorf("create tmp file %s: %v", tmp, errCreate)
			return
		}

		printf := func(format string, v ...interface{}) (int, error) {
			s := fmt.Sprintf(format, v...)
			_, err := w.Write([]byte(s))
			if err != nil {
				msg := fmt.Errorf("TestBuild printf: %v", err)
				t.Errorf(msg.Error())
				return 0, msg
			}
			return 0, nil
		}

		r := strings.NewReader(data.source)

		status, errors := compile(r, printf)

		w.Close()

		t.Logf("status=%d errors=%d\n", status, errors)

		if data.buildError == WRONG {
			// build error expected
			if status == 0 && errors == 0 {
				t.Errorf("unexpected build success")
				return
			}
			continue
		} else {
			// build error NOT expected
			if status != 0 {
				t.Errorf("unexpected build status=%d", status)
				return
			}
			if errors != 0 {
				t.Errorf("unexpected build errors=%d", errors)
				return
			}
		}

		cmd := exec.Command("go", "run", tmp)
		cmd.Stdin = strings.NewReader(data.input)
		output := bytes.Buffer{}
		cmd.Stdout = &output
		errExec := cmd.Run()

		if data.buildError == RUNTIME {
			// RUNTIME error expected
			if errExec == nil {
				t.Errorf("unexpected RUNTIME success")
				return
			}
		} else {
			// RUNTIME error NOT expected
			if errExec != nil {
				t.Errorf("unexpected RUNTIME error: go run %s: %v", tmp, errExec)
				return
			}
		}

		result := output.String()
		t.Logf("output: %q\n", result)
		if result != data.output {
			t.Errorf("unexpected output: got=%q expected=%q", result, data.output)
			return
		}
	}

}
