package baslib

import (
	"os"
	"strings"
)

func Environ(s string) {
	eq := strings.IndexByte(s, '=')
	if eq < 0 {
		alert("ENVIRON missing equal sign")
		return
	}
	k := s[:eq]
	if k == "" {
		alert("ENVIRON empty var name")
		return
	}
	v := s[eq+1:]
	if v == "" {
		err := os.Unsetenv(k)
		if err != nil {
			alert("ENVIRON deleting var: %v", err)
		}
		return
	}
	err := os.Setenv(k, v)
	if err != nil {
		alert("ENVIRON setting var: %v", err)
	}
}

func EnvironKey(key string) string {
	return os.Getenv(key)
}

func EnvironIndex(i int) string {
	if i < 1 {
		alert("ENVIRON$ index underflow: %d", i)
		return ""
	}
	list := os.Environ()
	if i > len(list) {
		return ""
	}
	return list[i-1]
}
