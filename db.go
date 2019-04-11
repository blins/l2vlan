package main

import (
	"go.etcd.io/bbolt"
)

func StartDatabase(dbname string) error {
	var err error
	db, err = bbolt.Open(dbname, 0666, nil)
	return err
}

func ShutdownDatabase() {
	db.Close()
}

var (
	db *bbolt.DB
)

