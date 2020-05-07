package main

import (
	"errors"
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/sirupsen/logrus"
	"github.com/snapcore/snapd/osutil"
	"time"
)

var db *bolt.DB

const (
	dbPath   = "data.db"
	dbBucket = "records"
)

var (
	errBucketEmpty = errors.New("[error] bucket are  empty")
	errValEmpty    = errors.New("[error] value are empty")
)

func init() {
	if !osutil.FileExists(dbPath) {
		if err := startUp(false); err != nil {
			logrus.Fatal(err)
		}

		defer func() {
			if err := db.Close(); err != nil {
				logrus.Fatal(err)
			}
		}()
	}
}

func startUp(readOnly bool) error {
	var err error
	db, err = bolt.Open(dbPath, 0666,
		&bolt.Options{
			Timeout:  1 * time.Second,
			ReadOnly: readOnly,
		},
	)
	if err != nil {
		return fmt.Errorf("[error opening  the database] %s", err.Error())
	}

	return nil
}

func dbWrite(key, val []byte) error {
	if err := startUp(false); err != nil {
		return err
	}

	defer func() {
		if err := db.Close(); err != nil {
			logrus.Fatalf("[error closing the database] %s", err.Error())
		}
	}()

	return db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(dbBucket))
		if err != nil {
			return err
		}

		if err = b.Put(key, val); err != nil {
			return err
		}
		return nil
	})
}

func dbRead(key string) ([]byte, error) {
	err := startUp(true)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := db.Close(); err != nil {
			logrus.Fatalf("[error closing the database] %s", err.Error())
		}
	}()

	var val []byte
	if err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(dbBucket))
		if b == nil {
			return errBucketEmpty
		}

		result := b.Get([]byte(key))
		val = make([]byte, len(result))
		copy(val, result)

		if len(val) == 0 {
			return errValEmpty
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return val, nil
}

func dbDelete(key string) error {
	err := startUp(false)
	if err != nil {
		return err
	}

	defer func() {
		if err := db.Close(); err != nil {
			logrus.Fatalf("Error close db: %s", err)
		}
	}()

	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(dbBucket))
		if b == nil {
			return errBucketEmpty
		}
		return b.Delete([]byte(key))
	})
}
