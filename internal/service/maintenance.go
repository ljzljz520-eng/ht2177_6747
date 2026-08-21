package service

import (
	"errors"
	"storeledger/internal/domain"
)

type BatchHealth struct {
	BatchID        string `json:"batch_id"`
	Status         string `json:"status"`
	Records        int    `json:"records"`
	ActionRequired int    `json:"action_required"`
	LedgerValid    bool   `json:"ledger_valid"`
}

func (s *Service) BatchHealth(batchID string) (BatchHealth, error) {
	if err := s.ensure(); err != nil {
		return BatchHealth{}, err
	}
	batch, err := s.GetBatch(batchID)
	if err != nil {
		return BatchHealth{}, err
	}
	records, err := s.ListBatchRecords(batchID)
	if err != nil {
		return BatchHealth{}, err
	}
	health := BatchHealth{BatchID: batchID, Status: batch.Status, Records: len(records)}
	for _, record := range records {
		if domain.IsActionRequired(record) {
			health.ActionRequired++
		}
	}
	health.LedgerValid = s.ValidateLedger(batchID) == nil
	return health, nil
}

func (s *Service) ReopenRecord(recordID string) (domain.InspectionRecord, error) {
	if err := s.ensure(); err != nil {
		return domain.InspectionRecord{}, err
	}
	record, err := s.GetRecord(recordID)
	if err != nil {
		return record, err
	}
	if record.Status != domain.RecordArchived {
		return record, errors.New("only archived records can be reopened")
	}
	record.Status = domain.RecordReviewed
	return record, s.Store.UpdateRecord(record)
}

func (s *Service) DeleteAttachment(id string) error {
	if id == "" {
		return errors.New("attachment id is required")
	}
	return nil
}

func (s *Service) ValidateAllBatches() (int, []string, error) {
	batches, err := s.Store.ListBatches()
	if err != nil {
		return 0, nil, err
	}
	invalid := make([]string, 0)
	for _, batch := range batches {
		if err := s.ValidateLedger(batch.BatchID); err != nil {
			invalid = append(invalid, batch.BatchID)
		}
	}
	return len(batches), invalid, nil
}
