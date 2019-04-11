package main

import (
	"testing"
	"os"
)

const (
	TestDatabase = "test.db"
	TestDatabaseErr = "/tmp"
)

func TestStartDatabase(t *testing.T) {
	err := StartDatabase(TestDatabase)
	if err != nil {
		t.Error(err)
	}
	ShutdownDatabase()
	os.Remove(TestDatabase)
}

func TestShutdownDatabase(t *testing.T) {
	err := StartDatabase(TestDatabase)
	if err != nil {
		t.Error(err)
	}
	ShutdownDatabase()
	os.Remove(TestDatabase)
}

func TestStartDatabase2(t *testing.T) {
	err := StartDatabase(TestDatabaseErr)
	if err == nil {
		t.Error()
	}
}
