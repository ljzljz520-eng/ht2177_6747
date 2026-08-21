package service

import (
	"errors"

	"storeledger/internal/domain"
	"storeledger/internal/parser"
)

type ImportResult struct {
	Batch       domain.InspectionBatch
	Records     []domain.InspectionRecord
	Attachments []domain.Attachment
	Summary     string
}

func (s *Service) ImportAndValidate(doc parser.BatchDocument) (ImportResult, error) {
	if err := s.ensure(); err != nil {
		return ImportResult{}, err
	}
	if err := parser.ValidateDocument(doc); err != nil {
		return ImportResult{}, err
	}
	batch, err := s.RegisterBatch(doc.BatchID, doc.Title, doc.Source)
	if err != nil {
		return ImportResult{}, err
	}
	records, err := parser.BuildRecords(doc)
	if err != nil {
		return ImportResult{}, err
	}
	for _, record := range records {
		if err := s.ValidateRecord(record); err != nil {
			return ImportResult{}, err
		}
	}
	ordered := domain.SortRecords(records)
	if err := s.Store.ReplaceRecords(ordered); err != nil {
		return ImportResult{}, err
	}
	for _, attachment := range parser.BuildAttachments(doc, ordered) {
		if err := s.Store.PutAttachment(attachment); err != nil {
			return ImportResult{}, err
		}
	}
	confirmed, err := s.Store.ConfirmBatch(batch.BatchID, len(ordered))
	if err != nil {
		return ImportResult{}, err
	}
	result, err := s.Report(batch.BatchID)
	if err != nil {
		return ImportResult{}, err
	}
	return ImportResult{Batch: confirmed, Records: cloneRecords(ordered), Attachments: parser.BuildAttachments(doc, ordered), Summary: result.Text()}, nil
}

func cloneRecords(records []domain.InspectionRecord) []domain.InspectionRecord {
	out := make([]domain.InspectionRecord, len(records))
	for i, record := range records {
		out[i] = record.Clone()
	}
	return out
}

func (s *Service) ReimportAndValidate(doc parser.BatchDocument) (ImportResult, error) {
	if doc.BatchID == "" {
		return ImportResult{}, errors.New("batch id is required")
	}
	return s.ImportAndValidate(doc)
}

func (s *Service) ConfirmIndependent(batchID string) error {
	batch, err := s.GetBatch(batchID)
	if err != nil {
		return err
	}
	if batch.Status != domain.BatchConfirmed {
		return errors.New("batch is not confirmed")
	}
	return nil
}
