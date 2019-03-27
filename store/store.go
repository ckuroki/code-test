package store

import (
	"errors"
	"github.com/boltdb/bolt"
)

var (
  bucket = []byte("urls")    // Bucket name
	errorNoData = errors.New("no data found")
)

type DB struct {
	db *bolt.DB
}

// Open key value store
func Open(dbfile string) (*DB, error) {
	if db, err := bolt.Open(dbfile, 0600, nil); err != nil {
		return nil, err
	} else {
		err := db.Update(func(tx *bolt.Tx) error {
			_, err := tx.CreateBucketIfNotExists(bucket)
			return err
		})
		if err != nil {
			return nil, err
		} else {
			return &DB{db: db}, nil
		}
	}
}

// Put inserts data into the store.
func (s *DB) Put(key string, value string) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket(bucket).Put([]byte(key), []byte(value))
	})
}

// Get retrieves data from the store.
func (s *DB) Get(key string) (val string,err error) {
	err = s.db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket(bucket).Cursor()
		if k, v := c.Seek([]byte(key)); k == nil || string(k) != key {
			return errorNoData
		} else {
      val = string(v)
			return nil
		}
	})
return
}

// Delete removes an entry from store
func (s *DB) Delete(key string) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		c := tx.Bucket(bucket).Cursor()
		if k, _ := c.Seek([]byte(key)); k == nil || string(k) != key {
			return errorNoData
		} else {
			return c.Delete()
		}
	})
}

// Close 
func (s *DB) Close() error {
	return s.db.Close()
}
