package service

import (
	"errors"
	"storeledger/internal/domain"
)

func (s *Service) BuildBatchLedger(batchID string) (domain.RecordLedger, error) {
	if err := s.ensure(); err != nil {
		return domain.RecordLedger{}, err
	}
	batch, err := s.GetBatch(batchID)
	if err != nil {
		return domain.RecordLedger{}, err
	}
	records, err := s.Store.RecordsForBatch(batchID)
	if err != nil {
		return domain.RecordLedger{}, err
	}
	if batch.RecordCount != len(records) {
		return domain.RecordLedger{}, errors.New("batch record count is stale")
	}
	return domain.BuildLedger(batchID, records)
}

func (s *Service) AddToLedger(batchID, recordID string) (domain.RecordLedger, error) {
	ledger, err := s.BuildBatchLedger(batchID)
	if err != nil {
		return ledger, err
	}
	updated, err := ledger.Add(recordID)
	if err != nil {
		return ledger, err
	}
	return updated, updated.Validate()
}

func (s *Service) InsertLedgerRecord(batchID, before, recordID string) (domain.RecordLedger, error) {
	ledger, err := s.BuildBatchLedger(batchID)
	if err != nil {
		return ledger, err
	}
	updated, err := ledger.InsertBefore(before, recordID)
	if err != nil {
		return ledger, err
	}
	return updated, updated.Validate()
}

func (s *Service) RemoveLedgerRecord(batchID, recordID string) (domain.RecordLedger, error) {
	ledger, err := s.BuildBatchLedger(batchID)
	if err != nil {
		return ledger, err
	}
	updated, err := ledger.Remove(recordID)
	if err != nil {
		return ledger, err
	}
	return updated, updated.Validate()
}

func (s *Service) Neighbors(batchID, recordID string) (string, string, error) {
	ledger, err := s.BuildBatchLedger(batchID)
	if err != nil {
		return "", "", err
	}
	previous, _ := ledger.Previous(recordID)
	next, _ := ledger.Next(recordID)
	if !ledger.Contains(recordID) {
		return "", "", errors.New("record not found in ledger")
	}
	return previous, next, nil
}

func (s *Service) ValidateLedger(batchID string) error {
	ledger, err := s.BuildBatchLedger(batchID)
	if err != nil {
		return err
	}
	records, err := s.ListBatchRecords(batchID)
	if err != nil {
		return err
	}
	return domain.ReconcileLedger(ledger, records)
}
