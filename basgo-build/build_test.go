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
	buildError bool
}

const sourceGoto = `
10 goto 900
700 print 3
710 end 
800 print 2
810 goto 700
900 print 1
910 goto 800
`

const outputGoto = `1
2
3
`

const (
	OK    = false
	WRONG = true
)

var testTable = []buildTest{
	{"", "", "", OK},       // empty program
	{"ugh", "", "", WRONG}, // invalid program

	{"10 print 1+2", "", "3\n", OK},
	{"10 print 1.1+2", "", "3.1\n", OK},
	{"10 print 1.1+2.2", "", "3.3\n", OK},
	{`10 print "a"+"b"`, "", "ab\n", OK},
	{`10 print 1+"b"`, "", "", WRONG},

	{"10 print 1-2", "", "-1\n", OK},
	{"10 print 1.1-2", "", "-0.9\n", OK},
	{"10 print 1.1-2.2", "", "-1.1\n", OK},
	{`10 print "a"-"b"`, "", "", WRONG},

	{"10 print 5 MOD 3", "", "2\n", OK},
	{"10 print 5.5 MOD 3.3", "", "0\n", OK},
	{`10 print "a" MOD "b"`, "", "", WRONG},

	{`10 print 5 \ 3`, "", "1\n", OK},
	{`10 print 5 \ 2`, "", "2\n", OK},
	{`10 print 5 \ 2.5`, "", "1\n", OK},
	{`10 print 6.6 \ 3.3`, "", "2\n", OK},
	{`10 print "a" \ "b"`, "", "", WRONG},

	{`10 print 5 * 3`, "", "15\n", OK},
	{`10 print 1.1 * 2`, "", "2.2\n", OK},
	{`10 print 2 * 2.5`, "", "5\n", OK},
	{`10 print "a" * "b"`, "", "", WRONG},

	{`10 print 5 / 5`, "", "1\n", OK},
	{`10 print 5 / 4`, "", "1.25\n", OK},
	{`10 print 5 / 2`, "", "2.5\n", OK},
	{`10 print 5 / 1`, "", "5\n", OK},
	{`10 print 5 / 2.5`, "", "2\n", OK},
	{`10 print 6.6 / 3.3`, "", "2\n", OK},
	{`10 print "a" / "b"`, "", "", WRONG},

	{`10 print 2 ^ 3`, "", "8\n", OK},
	{`10 print 16 ^ .5`, "", "4\n", OK},
	{`10 print "a" ^ "b"`, "", "", WRONG},

	{`10 print +10`, "", "10\n", OK},
	{`10 print +(2.2-1)`, "", "1.2\n", OK},
	{`10 print +"a"`, "", "", WRONG},

	{`10 print -10`, "", "-10\n", OK},
	{`10 print -(2.2-1)`, "", "-1.2\n", OK},
	{`10 print -"a"`, "", "", WRONG},

	{`10 print (22)`, "", "22\n", OK},
	{`10 print ("a"+"b")`, "", "ab\n", OK},
	{`10 print 2*(3+4)`, "", "14\n", OK},

	{`10 print LEN "hello"`, "", "5\n", OK},
	{`10 print LEN 2`, "", "8\n", OK},
	{`10 print LEN 3.3`, "", "8\n", OK},

	{`10 print not 0`, "", "-1\n", OK},
	{`10 print not -1`, "", "0\n", OK},
	{`10 print not 1.1`, "", "-2\n", OK},
	{`10 print not ""`, "", "", WRONG},

	{"10 print -1 and -1", "", "-1\n", OK},
	{"10 print 0 and 0", "", "0\n", OK},
	{"10 print 1 and 0", "", "0\n", OK},
	{"10 print 1 and 3", "", "1\n", OK},
	{`10 print "" and ""`, "", "", WRONG},

	{"10 print -1 or -1", "", "-1\n", OK},
	{"10 print 0 or 0", "", "0\n", OK},
	{"10 print 1 or 0", "", "1\n", OK},
	{"10 print 1 or 3", "", "3\n", OK},
	{`10 print "" or ""`, "", "", WRONG},

	{"10 print -1 xor -1", "", "0\n", OK},
	{"10 print 0 xor -1", "", "-1\n", OK},
	{"10 print 0 xor 0", "", "0\n", OK},
	{"10 print 1 xor 0", "", "1\n", OK},
	{"10 print 1 xor 3", "", "2\n", OK},
	{`10 print "" xor ""`, "", "", WRONG},

	{"10 print -1 eqv -1", "", "-1\n", OK},
	{"10 print 0 eqv -1", "", "0\n", OK},
	{"10 print 0 eqv 0", "", "-1\n", OK},
	{`10 print "" eqv ""`, "", "", WRONG},

	{"10 print -1 imp -1", "", "-1\n", OK},
	{"10 print 0 imp -1", "", "-1\n", OK},
	{"10 print 0 imp 0", "", "-1\n", OK},
	{"10 print -1 imp 0", "", "0\n", OK},
	{`10 print "" imp ""`, "", "", WRONG},

	{"10 print 0 = 0", "", "-1\n", OK},
	{"10 print 1.1 = 1", "", "0\n", OK},
	{"10 print 2.2 = 2.2", "", "-1\n", OK},
	{`10 print "" = ""`, "", "-1\n", OK},
	{"10 print 0 = 1", "", "0\n", OK},
	{"10 print 1.1 = 2", "", "0\n", OK},
	{"10 print 2.2 = 3.3", "", "0\n", OK},
	{`10 print "a" = ""`, "", "0\n", OK},
	{`10 print 0 = ""`, "", "", WRONG},

	{"10 print 0 <> 0", "", "0\n", OK},
	{"10 print 1.1 <> 1", "", "-1\n", OK},
	{"10 print 2.2 <> 2.2", "", "0\n", OK},
	{`10 print "" <> ""`, "", "0\n", OK},
	{"10 print 0 <> 1", "", "-1\n", OK},
	{"10 print 1.1 <> 2", "", "-1\n", OK},
	{"10 print 2.2 <> 3.3", "", "-1\n", OK},
	{`10 print "a" <> ""`, "", "-1\n", OK},
	{`10 print 0 <> ""`, "", "", WRONG},

	{"10 print 0 > 0", "", "0\n", OK},
	{"10 print 1.1 > 1", "", "-1\n", OK},
	{"10 print 2.2 > 2.2", "", "0\n", OK},
	{`10 print "" > ""`, "", "0\n", OK},
	{`10 print "a" > "b"`, "", "0\n", OK},
	{"10 print 0 > 1", "", "0\n", OK},
	{"10 print 1 > 0", "", "-1\n", OK},
	{"10 print 1.1 > 2", "", "0\n", OK},
	{"10 print 2.2 > 3.3", "", "0\n", OK},
	{"10 print 3.3 > 2.2", "", "-1\n", OK},
	{`10 print "a" > ""`, "", "-1\n", OK},
	{`10 print 0 > ""`, "", "", WRONG},

	{"10 print 0 < 0", "", "0\n", OK},
	{"10 print 1.1 < 1", "", "0\n", OK},
	{"10 print 2.2 < 2.2", "", "0\n", OK},
	{`10 print "" < ""`, "", "0\n", OK},
	{`10 print "a" < "b"`, "", "-1\n", OK},
	{"10 print 0 < 1", "", "-1\n", OK},
	{"10 print 1 < 0", "", "0\n", OK},
	{"10 print 1.1 < 2", "", "-1\n", OK},
	{"10 print 2.2 < 3.3", "", "-1\n", OK},
	{"10 print 3.3 < 2.2", "", "0\n", OK},
	{`10 print "a" < ""`, "", "0\n", OK},
	{`10 print 0 < ""`, "", "", WRONG},

	{"10 print 0 >= 0", "", "-1\n", OK},
	{"10 print 1.1 >= 1", "", "-1\n", OK},
	{"10 print 2.2 >= 2.2", "", "-1\n", OK},
	{`10 print "" >= ""`, "", "-1\n", OK},
	{`10 print "a" >= "b"`, "", "0\n", OK},
	{"10 print 0 >= 1", "", "0\n", OK},
	{"10 print 1 >= 0", "", "-1\n", OK},
	{"10 print 1.1 >= 2", "", "0\n", OK},
	{"10 print 2.2 >= 3.3", "", "0\n", OK},
	{"10 print 3.3 >= 2.2", "", "-1\n", OK},
	{`10 print "a" >= ""`, "", "-1\n", OK},
	{`10 print 0 >= ""`, "", "", WRONG},

	{"10 print 0 <= 0", "", "-1\n", OK},
	{"10 print 1.1 <= 1", "", "0\n", OK},
	{"10 print 2.2 <= 2.2", "", "-1\n", OK},
	{`10 print "" <= ""`, "", "-1\n", OK},
	{`10 print "a" <= "b"`, "", "-1\n", OK},
	{"10 print 0 <= 1", "", "-1\n", OK},
	{"10 print 1 <= 0", "", "0\n", OK},
	{"10 print 1.1 <= 2", "", "-1\n", OK},
	{"10 print 2.2 <= 3.3", "", "-1\n", OK},
	{"10 print 3.3 <= 2.2", "", "0\n", OK},
	{`10 print "a" <= ""`, "", "0\n", OK},
	{`10 print 0 <= ""`, "", "", WRONG},

	{`10 goto 20`, "", "", WRONG},
	{sourceGoto, "", outputGoto, OK},

	{`10 print "hi"`, "", "hi\n", OK},
	{`10 print "hi";`, "", "hi", OK},

	{`10 if "" then print 1`, "", "", WRONG},
	{`10 if 0 then print 1`, "", "", OK},
	{`10 if -1 then print 1`, "", "1\n", OK},
	{`10 if 1.1 then print 1`, "", "1\n", OK},
	{`10 if 0 then print 1 else print 2`, "", "2\n", OK},
	{`10 if -1 then print 1 else print 2`, "", "1\n", OK},
	{`10 if 1.1 then print 1 else print 2`, "", "1\n", OK},

	{"10 if 0 then 20\n20 end\n30 print 30", "", "", OK},
	{"10 if 0 goto 20\n20 end\n30 print 30", "", "", OK},
	{"10 if 0 then 20 else 30\n20 end\n30 print 30", "", "30\n", OK},
	{"10 if 0 goto 20 else 30\n20 end\n30 print 30", "", "30\n", OK},

	{`10 if 1 then print 1:print 2:print 3 else print 4:print 5:print 6`, "", "1\n2\n3\n", OK},
	{`10 if 0 then print 1:print 2:print 3 else print 4:print 5:print 6`, "", "4\n5\n6\n", OK},

	{`10 print "abc" : if 1 then print 1:print 2:print 3 else print 4:print 5:print 6`, "", "abc\n1\n2\n3\n", OK},
	{`10 print "abc" : if 0 then print 1:print 2:print 3 else print 4:print 5:print 6`, "", "abc\n4\n5\n6\n", OK},

	{`10 input a : print a`, "2\n", "2\n", OK},
	{`10 input a! : print a!`, "2.1\n", "2.1\n", OK},
	{`10 input a# : print a#`, "2.1\n", "2.1\n", OK},
	{`10 input a% : print a%`, "2\n", "2\n", OK},
	{`10 input a$ : print a$`, "abc\n", "abc\n", OK},
	{`10 input a:input b:print a" "b;`, "2\n3\n", "2 3", OK},

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

	{`10 a=a%:print a`, "", "0\n", OK},
	{`10 a=a#:print a`, "", "0\n", OK},
	{`10 a%=a:print a%`, "", "0\n", OK},
	{`10 a%=a#:print a%`, "", "0\n", OK},
	{`10 a#=a:print a#`, "", "0\n", OK},
	{`10 a#=a%:print a#`, "", "0\n", OK},

	{`10 a%=5.1:print a%`, "", "5\n", OK},
	{`10 a%=5.9:print a%`, "", "6\n", OK},

	{`10 a=45:print A`, "", "45\n", OK},
	{`10 A=45:print a`, "", "45\n", OK},
	{`10 input a:print A;`, "23\n", "23", OK},
	{`10 input A:print a;`, "23\n", "23", OK},

	{`10 print int 2`, "", "2\n", OK},
	{`10 print int 2.1`, "", "2\n", OK},
	{`10 print int 2.9`, "", "2\n", OK},
	{`10 print int ""`, "", "", WRONG},

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

	{`10 print 1:end:print 2`, "", "1\n", OK},
	{`10 print 1:stop:print 2`, "", "1\n", OK},

	{"10 on 1 goto 20\n20 rem", "", "", OK},
	{"10 on 2 goto 10,20\n20 rem", "", "", OK},
	{`10 on 1 goto x`, "", "", WRONG},
	{`10 on 1 goto 20`, "", "", WRONG},
	{`10 on "" goto 10`, "", "", WRONG},
	{`10 on a$ goto 10`, "", "", WRONG},
	{"10 on 1 goto 20,30,40\n20 print 20:end\n30 print 30:end\n40 print 40", "", "20\n", OK},
	{"10 on 2 goto 20,30,40\n20 print 20:end\n30 print 30:end\n40 print 40", "", "30\n", OK},
	{"10 on 3 goto 20,30,40\n20 print 20:end\n30 print 30:end\n40 print 40", "", "40\n", OK},
	{"10 on 4 goto 20,30,40\n20 print 20:end\n30 print 30:end\n40 print 40", "", "20\n", OK},
	{"10 a=1:on a goto 30\n20 print 20:end\n30 print 30", "", "30\n", OK},
	{"10 a=1.1:on a goto 30\n20 print 20:end\n30 print 30", "", "30\n", OK},
}

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

		t.Logf("status=%d errors=%d\n", status, errors)

		if data.buildError {
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

		w.Close()

		cmd := exec.Command("go", "run", tmp)
		cmd.Stdin = strings.NewReader(data.input)
		output := bytes.Buffer{}
		cmd.Stdout = &output
		errExec := cmd.Run()
		if errExec != nil {
			t.Errorf("go run %s: %v", tmp, errExec)
			return
		}
		result := output.String()
		t.Logf("output: %q\n", result)
		if result != data.output {
			t.Errorf("unexpected output: got=%q expected=%q", result, data.output)
			return
		}
	}

}
