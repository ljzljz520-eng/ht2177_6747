package store

import (
	"go.etcd.io/bbolt"
	"storeledger/internal/domain"
)

func (s *Store) PutAudit(event domain.AuditEvent) error {
	if err := validateEntityID(event.ID); err != nil {
		return err
	}
	return s.transaction(true, func(tx *bbolt.Tx) error { return put(tx, "audits", event.ID, event) })
}

func (s *Store) ListAudits() ([]domain.AuditEvent, error) {
	return list[domain.AuditEvent](s.db, "audits")
}

func (s *Store) AuditsForRecord(recordID string) ([]domain.AuditEvent, error) {
	all, err := s.ListAudits()
	if err != nil {
		return nil, err
	}
	result := make([]domain.AuditEvent, 0)
	for _, event := range all {
		if event.RecordID == recordID {
			result = append(result, event)
		}
	}
	return result, nil
}

func (s *Store) LastDecision(recordID string) (domain.AuditEvent, bool, error) {
	events, err := s.AuditsForRecord(recordID)
	if err != nil {
		return domain.AuditEvent{}, false, err
	}
	if len(events) == 0 {
		return domain.AuditEvent{}, false, nil
	}
	return events[len(events)-1], true, nil
}
