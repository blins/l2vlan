package main

import "testing"

const (
	avIntValue = 1000
	avStringValue = "1000"
)

type avStringer struct{}

func (s avStringer) String() string {
	return avStringValue
}

func TestAnyVal_Int(t *testing.T) {
	var i int = avIntValue
	av := AnyVal{i}
	if av.Int() != avIntValue {
		t.Error()
	}
}

func TestAnyVal_Int2(t *testing.T) {
	var i string = avStringValue
	av := AnyVal{i}
	if av.Int() != avIntValue {
		t.Error()
	}
}

func TestAnyVal_Int3(t *testing.T) {
	var i avStringer
	av := AnyVal{i}
	if av.Int() != avIntValue {
		t.Error()
	}
}

func TestAnyVal_Int4(t *testing.T) {
	var i []byte
	av := AnyVal{i}
	if av.Int() != 0 {
		t.Error()
	}
}

func TestAnyVal_String(t *testing.T) {
	var i int = avIntValue
	av := AnyVal{i}
	if av.String() != avStringValue {
		t.Error()
	}
}

func TestAnyVal_String2(t *testing.T) {
	var i string = avStringValue
	av := AnyVal{i}
	if av.String() != avStringValue {
		t.Error()
	}
}

func TestAnyVal_String3(t *testing.T) {
	var i avStringer
	av := AnyVal{i}
	if av.String() != avStringValue {
		t.Error()
	}
}

func TestAnyVal_String4(t *testing.T) {
	var i []byte
	av := AnyVal{i}
	if av.String() != "" {
		t.Error()
	}
}