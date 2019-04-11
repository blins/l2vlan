package main

import (
	"fmt"
	"bytes"
	"testing"
)

type DebuggerTest struct {
	value string
}

func (dt *DebuggerTest) Println(v ...interface{}) {
	var buffer bytes.Buffer
	fmt.Fprintln(&buffer, v...)
	dt.value = string(buffer.Bytes())
}

func TestDebugger_SetLogger(t *testing.T) {
	var d Debugger
	d.SetLogger(&DebuggerTest{})
	if d.logger == nil {
		t.Error()
	}
}

func TestDebugger_Println(t *testing.T) {
	var d Debugger
	var dt DebuggerTest
	d.SetLogger(&dt)
	d.Println("test", 1000)
	if dt.value != "test 1000\n" {
		t.Error()
	}
}