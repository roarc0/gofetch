package memory

import (
	"io"
	"log"

	bolt "go.etcd.io/bbolt"
)

type Memory interface {
	Put(key string, value string) error
	Get(key string) (*string, error)
	Has(key string) bool
	Del(key string) error
	io.Closer
}

type memory struct {
	db     *bolt.DB
	bucket []byte
}

type Config struct {
	FilePath string
}

func NewMemory(filePath string, bucket string) (Memory, error) {
	db, err := bolt.Open(filePath, 0600, nil)
	if err != nil {
		return nil, err
	}

	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucket))
		return err
	})

	if err != nil {
		return nil, err
	}

	return &memory{
		db:     db,
		bucket: []byte(bucket),
	}, nil
}

func (m *memory) Put(key string, value string) error {
	return m.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(m.bucket)
		err := b.Put([]byte(key), []byte(value))
		return err
	})
}

func (m *memory) Get(key string) (*string, error) {
	var value string
	err := m.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(m.bucket)
		v := b.Get([]byte(key))
		value = string(v)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return &value, nil
}

func (m *memory) Has(key string) bool {
	var exists bool
	err := m.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(m.bucket)
		v := b.Get([]byte(key))
		exists = v != nil
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	return exists
}

func (m *memory) Del(key string) error {
	return m.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(m.bucket)
		err := b.Delete([]byte(key))
		return err
	})
}

func (m *memory) Close() error {
	return m.db.Close()
}
