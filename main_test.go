package main

import (
	"testing"
)

func TestNewNetwork(t *testing.T) {
	ss := NetToCIDR(*nets)
	n := NewNetwork()
	if n.IP.IsUnspecified() {
		t.Error()
	}
	if s := NetToCIDR(n); s == ss {
		t.Error("Invalid", s)
	}
}


