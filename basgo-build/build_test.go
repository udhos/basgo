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

var testTable = []buildTest{
	{"", "", "", false},   // empty program
	{"ugh", "", "", true}, // invalid program
	{"10 print 1+2", "", "3\n", false},
	{"10 print 1.1+2", "", "3.1\n", false},
	{"10 print 1.1+2.2", "", "3.3\n", false},
	{`10 print "a"+"b"`, "", "ab\n", false},
	{`10 print 1+"b"`, "", "", true},
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
			}
			continue
		} else {
			// build error NOT expected
			if status != 0 {
				t.Errorf("unexpected build status=%d", status)
			}
			if errors != 0 {
				t.Errorf("unexpected build errors=%d", errors)
			}
		}

		w.Close()

		if status != 0 || errors != 0 {
			continue
		}

		cmd := exec.Command("go", "run", tmp)
		cmd.Stdin = strings.NewReader(data.input)
		output := bytes.Buffer{}
		cmd.Stdout = &output
		errExec := cmd.Run()
		if errExec != nil {
			t.Errorf("go run %s: %v", tmp, errExec)
			continue
		}
		result := output.String()
		t.Logf("output: %q\n", result)
		if result != data.output {
			t.Errorf("unexpected output: got=[%s] expected=[%s]", result, data.output)
		}
	}

}
