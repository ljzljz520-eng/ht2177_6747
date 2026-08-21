package store

import (
	"errors"
	"time"

	"go.etcd.io/bbolt"
)

func (s *Store) PurgeBucket(bucket string) error {
	allowed := map[string]bool{"audits": true, "attachments": true, "notes": true}
	if !allowed[bucket] {
		return errors.New("bucket cannot be purged")
	}
	return s.transaction(true, func(tx *bbolt.Tx) error {
		cursor := tx.Bucket([]byte(bucket)).Cursor()
		keys := make([][]byte, 0)
		for key, _ := cursor.First(); key != nil; key, _ = cursor.Next() {
			keys = append(keys, append([]byte(nil), key...))
		}
		for _, key := range keys {
			if err := tx.Bucket([]byte(bucket)).Delete(key); err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *Store) CompactHint() string {
	if s == nil || !s.Healthy() {
		return "closed"
	}
	return time.Unix(0, 0).UTC().Format(time.RFC3339)
}

func (s *Store) BucketSizes() (map[string]int, error) {
	result := make(map[string]int)
	for _, bucket := range []string{"batches", "records", "audits", "attachments", "notes"} {
		count, err := s.Count(bucket)
		if err != nil {
			return nil, err
		}
		result[bucket] = count
	}
	return result, nil
}
