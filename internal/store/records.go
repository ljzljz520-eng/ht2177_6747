package store

import (
	"errors"
	"go.etcd.io/bbolt"
	"storeledger/internal/domain"
)

func (s *Store) PutRecord(record domain.InspectionRecord) error {
	if err := record.Validate(); err != nil {
		return err
	}
	return s.transaction(true, func(tx *bbolt.Tx) error { return put(tx, "records", record.ID, record) })
}

func (s *Store) GetRecord(id string) (domain.InspectionRecord, error) {
	var record domain.InspectionRecord
	if err := validateEntityID(id); err != nil {
		return record, err
	}
	err := s.transaction(false, func(tx *bbolt.Tx) error { return get(tx, "records", id, &record) })
	if isMissing(err) {
		return record, errors.New("record not found")
	}
	return record, err
}

func (s *Store) ListRecords() ([]domain.InspectionRecord, error) {
	return list[domain.InspectionRecord](s.db, "records")
}

func (s *Store) ReplaceRecords(records []domain.InspectionRecord) error {
	return s.transaction(true, func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte("records"))
		for _, record := range records {
			if err := record.Validate(); err != nil {
				return err
			}
			data, err := encode(record)
			if err != nil {
				return err
			}
			if err := bucket.Put([]byte(record.ID), data); err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *Store) UpdateRecord(record domain.InspectionRecord) error { return s.PutRecord(record) }

func (s *Store) DeleteRecord(id string) error {
	if err := validateEntityID(id); err != nil {
		return err
	}
	return s.transaction(true, func(tx *bbolt.Tx) error { return remove(tx, "records", id) })
}

func (s *Store) RecordsForBatch(batchID string) ([]domain.InspectionRecord, error) {
	all, err := s.ListRecords()
	if err != nil {
		return nil, err
	}
	result := make([]domain.InspectionRecord, 0)
	for _, record := range all {
		if record.BatchID == batchID {
			result = append(result, record)
		}
	}
	return domain.SortRecords(result), nil
}
