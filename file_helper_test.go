package ap_tools

import (
	"fmt"
	"testing"
)

func TestReadLineByLineAndProcess(t *testing.T) {
	err := ReadLineByLineAndProcess("test.txt", func(lineNumber int64, lineString string) {
		fmt.Println(lineNumber, lineString)
	})
	if err != nil {
		t.Fatalf(err.Error())
	}
}

type people struct {
	Name string `index:"1"`
	Age  int64  `index:"2"`
	Addr string `index:"3"`
}

func TestLineString2Struct(t *testing.T) {
	err := ReadLineByLineAndProcess("test.txt", func(lineNumber int64, lineString string) {
		p := people{}
		e := LineString2Struct(&p, lineString, ",")
		if e != nil {
			t.Fatalf(e.Error())
		}
		fmt.Println(lineNumber, p.Name, p.Age, p.Addr)
	})
	if err != nil {
		t.Fatalf(err.Error())
	}
}
