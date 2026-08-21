package store

import (
	"encoding/json"
	"errors"
	"sort"

	"go.etcd.io/bbolt"
	"storeledger/internal/domain"
)

type Snapshot struct {
	Batches     []domain.InspectionBatch   `json:"batches"`
	Records     []domain.InspectionRecord  `json:"records"`
	Audits      []domain.AuditEvent        `json:"audits"`
	Attachments []domain.Attachment        `json:"attachments"`
	Notes       []domain.CollaborationNote `json:"notes"`
}

func (s *Store) Snapshot() (Snapshot, error) {
	if s == nil {
		return Snapshot{}, errors.New("store is nil")
	}
	batches, err := s.ListBatches()
	if err != nil {
		return Snapshot{}, err
	}
	records, err := s.ListRecords()
	if err != nil {
		return Snapshot{}, err
	}
	audits, err := s.ListAudits()
	if err != nil {
		return Snapshot{}, err
	}
	attachments, err := s.ListAttachments()
	if err != nil {
		return Snapshot{}, err
	}
	notes, err := s.ListNotes()
	if err != nil {
		return Snapshot{}, err
	}
	sort.Slice(batches, func(i, j int) bool { return batches[i].BatchID < batches[j].BatchID })
	sort.Slice(records, func(i, j int) bool { return records[i].ID < records[j].ID })
	return Snapshot{Batches: batches, Records: records, Audits: audits, Attachments: attachments, Notes: notes}, nil
}

func (s *Store) ExportJSON() ([]byte, error) {
	snapshot, err := s.Snapshot()
	if err != nil {
		return nil, err
	}
	return json.MarshalIndent(snapshot, "", "  ")
}

func (s *Store) ImportSnapshot(snapshot Snapshot) error {
	return s.transaction(true, func(tx *bbolt.Tx) error {
		for _, batch := range snapshot.Batches {
			if err := put(tx, "batches", batch.BatchID, batch); err != nil {
				return err
			}
		}
		for _, record := range snapshot.Records {
			if err := put(tx, "records", record.ID, record); err != nil {
				return err
			}
		}
		for _, audit := range snapshot.Audits {
			if err := put(tx, "audits", audit.ID, audit); err != nil {
				return err
			}
		}
		for _, attachment := range snapshot.Attachments {
			if err := put(tx, "attachments", attachment.ID, attachment); err != nil {
				return err
			}
		}
		for _, note := range snapshot.Notes {
			if err := put(tx, "notes", note.ID, note); err != nil {
				return err
			}
		}
		return nil
	})
}

func DecodeSnapshot(data []byte) (Snapshot, error) {
	var snapshot Snapshot
	if len(data) == 0 {
		return snapshot, errors.New("snapshot is empty")
	}
	return snapshot, json.Unmarshal(data, &snapshot)
}
