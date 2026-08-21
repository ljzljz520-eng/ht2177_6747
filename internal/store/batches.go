package store

import (
	"errors"
	"go.etcd.io/bbolt"
	"storeledger/internal/domain"
)

func (s *Store) PutBatch(batch domain.InspectionBatch) error {
	if err := validateEntityID(batch.BatchID); err != nil {
		return err
	}
	return s.transaction(true, func(tx *bbolt.Tx) error { return put(tx, "batches", batch.BatchID, batch) })
}

func (s *Store) GetBatch(id string) (domain.InspectionBatch, error) {
	var batch domain.InspectionBatch
	if err := validateEntityID(id); err != nil {
		return batch, err
	}
	err := s.transaction(false, func(tx *bbolt.Tx) error { return get(tx, "batches", id, &batch) })
	if isMissing(err) {
		return batch, errors.New("batch not found")
	}
	return batch, err
}

func (s *Store) ListBatches() ([]domain.InspectionBatch, error) {
	return list[domain.InspectionBatch](s.db, "batches")
}

func (s *Store) ConfirmBatch(id string, count int) (domain.InspectionBatch, error) {
	var result domain.InspectionBatch
	err := s.transaction(true, func(tx *bbolt.Tx) error {
		var current domain.InspectionBatch
		if err := get(tx, "batches", id, &current); err != nil {
			return err
		}
		updated, err := current.Confirm(count)
		if err != nil {
			return err
		}
		result = updated
		return put(tx, "batches", id, updated)
	})
	return result, err
}

func (s *Store) ArchiveBatch(id string) (domain.InspectionBatch, error) {
	var result domain.InspectionBatch
	err := s.transaction(true, func(tx *bbolt.Tx) error {
		var current domain.InspectionBatch
		if err := get(tx, "batches", id, &current); err != nil {
			return err
		}
		updated, err := current.Archive()
		if err != nil {
			return err
		}
		result = updated
		return put(tx, "batches", id, updated)
	})
	return result, err
}
