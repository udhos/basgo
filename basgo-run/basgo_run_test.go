package main

import (
	"bytes"
	"testing"

	"github.com/udhos/basgo/basgo"
)

func TestListRun(t *testing.T) {
	loadListRun(t, program1, list1, output1)
	loadListRun(t, program2, list2, output2)
}

func loadListRun(t *testing.T, source, expectedList, expectedOutput string) {
	b := basgo.New()

	b.ExecuteString(source) // Load

	// redirect stdout to buf
	bufList := make([]byte, 1000)
	outList := bytes.NewBuffer(bufList)
	b.Out = outList

	b.ExecuteCommandList() // List

	resultList := outList.String()
	if expectedList != resultList {
		t.Errorf("LIST MISMATCH:\nexpected: [%v]\nresult: [%v]", expectedList, resultList)
	}

	// redirect stdout to buf
	bufRun := make([]byte, 1000)
	outRun := bytes.NewBuffer(bufRun)
	b.Out = outRun

	b.ExecuteCommandRun() // Run

	resultRun := outRun.String()
	if expectedOutput != resultRun {
		t.Errorf("RUN MISMATCH:\nexpected: [%v]\nresult: [%v]", expectedOutput, resultRun)
	}
}

const program1 = `10 cls
20 print "hi"
30 a$="world"
40 print a$
50 end
60 print "world"
70  print "bad"
80  print  "spaces"
`

const list1 = `10 CLS
20 PRINT "hi"
30 A$="world"
40 PRINT A$
50 END
60 PRINT "world"
70  PRINT "bad"
80  PRINT  "spaces"
`

const output1 = `hi
world
`

const program2 = `1000 PRINT:PRINT : : PRINT
2000 : : : :   : :
3000 PRINT ":" :: PRINT ":::"
`

const list2 = `1000 PRINT:PRINT : : PRINT
2000 : : : :   : :
3000 PRINT ":" :: PRINT ":::"
`

const output2 = `


  :
  :::
`
