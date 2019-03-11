package codepage

import (
	"testing"
)

type codepageTest struct {
	b int
	u int
}

var testTable = []codepageTest{
	{0x00, 0x0000},
	{0x01, 0x263a},
	{0x02, 0x263b},
	{0xfe, 0x25a0},
	{0xff, 0x00a0},
}

func TestCodepage(t *testing.T) {
	for _, c := range testTable {
		u := ByteToUnicode(c.b)
		if u != c.u {
			t.Errorf("wrong unicode for byte=%x: expected=%x found=%x", c.b, c.u, u)
		}
		b := UnicodeToByte(c.u)
		if b != c.b {
			t.Errorf("wrong byte for unicode=%x: expected=%x found=%x", c.u, c.b, b)
		}
	}
}
