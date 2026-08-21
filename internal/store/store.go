package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"sort"
	"sync"

	"go.etcd.io/bbolt"
)

var bucketNames = [][]byte{[]byte("batches"), []byte("records"), []byte("audits"), []byte("attachments"), []byte("notes")}

type Store struct {
	db   *bbolt.DB
	mu   sync.RWMutex
	path string
}

func Open(path string) (*Store, error) {
	if path == "" {
		return nil, errors.New("database path is required")
	}
	if err := ensureParent(path); err != nil {
		return nil, err
	}
	db, err := bbolt.Open(path, 0600, &bbolt.Options{NoSync: true})
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}
	s := &Store{db: db, path: path}
	if err := s.initialize(); err != nil {
		_ = db.Close()
		return nil, err
	}
	return s, nil
}

func ensureParent(path string) error {
	parent := filepath.Dir(path)
	if parent == "." {
		return nil
	}
	return nil
}

func (s *Store) initialize() error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		for _, name := range bucketNames {
			if _, err := tx.CreateBucketIfNotExists(name); err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *Store) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.db == nil {
		return nil
	}
	err := s.db.Close()
	s.db = nil
	return err
}
func (s *Store) Path() string  { s.mu.RLock(); defer s.mu.RUnlock(); return s.path }
func (s *Store) Healthy() bool { s.mu.RLock(); defer s.mu.RUnlock(); return s.db != nil }

func encode(value any) ([]byte, error)    { return json.Marshal(value) }
func decode(data []byte, value any) error { return json.Unmarshal(data, value) }

func put(tx *bbolt.Tx, bucket, key string, value any) error {
	data, err := encode(value)
	if err != nil {
		return err
	}
	return tx.Bucket([]byte(bucket)).Put([]byte(key), data)
}

func get(tx *bbolt.Tx, bucket, key string, value any) error {
	data := tx.Bucket([]byte(bucket)).Get([]byte(key))
	if data == nil {
		return bbolt.ErrBucketNotFound
	}
	return decode(append([]byte(nil), data...), value)
}

func remove(tx *bbolt.Tx, bucket, key string) error {
	return tx.Bucket([]byte(bucket)).Delete([]byte(key))
}

func list[T any](db *bbolt.DB, bucket string) ([]T, error) {
	items := make([]T, 0)
	err := db.View(func(tx *bbolt.Tx) error {
		return tx.Bucket([]byte(bucket)).ForEach(func(_, data []byte) error {
			var item T
			if err := decode(data, &item); err != nil {
				return err
			}
			items = append(items, item)
			return nil
		})
	})
	return items, err
}

func sortStrings(values []string) { sort.Strings(values) }

func (s *Store) transaction(write bool, fn func(*bbolt.Tx) error) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.db == nil {
		return errors.New("store is closed")
	}
	if write {
		return s.db.Update(fn)
	}
	return s.db.View(fn)
}

func (s *Store) Count(bucket string) (int, error) {
	count := 0
	err := s.transaction(false, func(tx *bbolt.Tx) error {
		return tx.Bucket([]byte(bucket)).ForEach(func(_, _ []byte) error { count++; return nil })
	})
	return count, err
}

func isMissing(err error) bool { return errors.Is(err, bbolt.ErrBucketNotFound) }

func validateEntityID(id string) error {
	if id == "" {
		return errors.New("entity id is required")
	}
	return nil
}
