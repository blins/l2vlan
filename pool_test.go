package main

import (
	"testing"
)

func TestNewPoolv4(t *testing.T) {
	pool := NewPoolv4()
	if pool == nil {
		t.Error()
	}
	if pool.Data == nil {
		t.Error()
	}
}
