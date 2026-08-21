package service

import (
	"errors"
	"storeledger/internal/domain"
	"storeledger/internal/parser"
)

type BulkResult struct {
	Imported int
	Failed   int
	Errors   []string
}

func (s *Service) ImportDocuments(documents []parser.BatchDocument) BulkResult {
	result := BulkResult{Errors: make([]string, 0)}
	for _, doc := range documents {
		if _, err := s.ImportAndValidate(doc); err != nil {
			result.Failed++
			result.Errors = append(result.Errors, err.Error())
			continue
		}
		result.Imported++
	}
	return result
}

func (s *Service) ReviewBatch(batchID string, approve bool, note string) (int, int, error) {
	records, err := s.ListBatchRecords(batchID)
	if err != nil {
		return 0, 0, err
	}
	reviewed, failed := 0, 0
	for _, record := range records {
		if _, _, err := s.ReviewRecord(record.ID, note, approve); err != nil {
			failed++
			continue
		}
		reviewed++
	}
	return reviewed, failed, nil
}

func (s *Service) ArchiveReviewed(batchID string) (int, error) {
	records, err := s.ListBatchRecords(batchID)
	if err != nil {
		return 0, err
	}
	archived := 0
	for _, record := range records {
		if record.Status != domain.RecordReviewed {
			continue
		}
		if _, err := s.ArchiveRecord(record.ID); err != nil {
			return archived, err
		}
		archived++
	}
	return archived, nil
}

func (s *Service) RequireConfirmed(batchID string) error {
	batch, err := s.GetBatch(batchID)
	if err != nil {
		return err
	}
	if batch.Status != domain.BatchConfirmed {
		return errors.New("batch must be confirmed")
	}
	return nil
}

func (s *Service) RecordCount(batchID string) (int, error) {
	records, err := s.ListBatchRecords(batchID)
	return len(records), err
}
