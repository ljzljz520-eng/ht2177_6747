package service

import (
	"errors"
	"storeledger/internal/domain"
)

func (s *Service) ReviewRecord(recordID, note string, approve bool) (domain.InspectionRecord, domain.AuditEvent, error) {
	if err := s.ensure(); err != nil {
		return domain.InspectionRecord{}, domain.AuditEvent{}, err
	}
	record, err := s.Store.GetRecord(recordID)
	if err != nil {
		return record, domain.AuditEvent{}, err
	}
	if approve {
		if err := s.Policy.CanApprove(record, note); err != nil {
			return record, domain.AuditEvent{}, err
		}
		record, event, err := domain.Approve(record, s.Reviewer, note, s.Clock.Now())
		if err != nil {
			return record, event, err
		}
		if err := s.Store.UpdateRecord(record); err != nil {
			return record, event, err
		}
		if err := s.Store.PutAudit(event); err != nil {
			return record, event, err
		}
		return record, event, nil
	}
	if note == "" {
		return record, domain.AuditEvent{}, errors.New("rejection note is required")
	}
	record, event, err := domain.Reject(record, s.Reviewer, note, s.Clock.Now())
	if err != nil {
		return record, event, err
	}
	if err := s.Store.UpdateRecord(record); err != nil {
		return record, event, err
	}
	if err := s.Store.PutAudit(event); err != nil {
		return record, event, err
	}
	return record, event, nil
}

func (s *Service) ArchiveRecord(recordID string) (domain.InspectionRecord, error) {
	if err := s.ensure(); err != nil {
		return domain.InspectionRecord{}, err
	}
	record, err := s.Store.GetRecord(recordID)
	if err != nil {
		return record, err
	}
	archived, err := domain.ArchiveRecord(record)
	if err != nil {
		return record, err
	}
	return archived, s.Store.UpdateRecord(archived)
}

func (s *Service) ArchiveBatch(batchID string) (domain.InspectionBatch, error) {
	if err := s.ensure(); err != nil {
		return domain.InspectionBatch{}, err
	}
	records, err := s.Store.RecordsForBatch(batchID)
	if err != nil {
		return domain.InspectionBatch{}, err
	}
	if len(records) == 0 {
		return domain.InspectionBatch{}, errors.New("batch has no records")
	}
	for _, record := range records {
		if record.Status != domain.RecordArchived {
			return domain.InspectionBatch{}, errors.New("all records must be archived")
		}
	}
	return s.Store.ArchiveBatch(batchID)
}

func (s *Service) Audit(recordID string) ([]domain.AuditEvent, error) {
	if err := s.ensure(); err != nil {
		return nil, err
	}
	return s.Store.AuditsForRecord(recordID)
}
