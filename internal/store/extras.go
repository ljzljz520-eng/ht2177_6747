package store

import (
	"go.etcd.io/bbolt"
	"storeledger/internal/domain"
)

func (s *Store) PutAttachment(item domain.Attachment) error {
	if err := validateEntityID(item.ID); err != nil {
		return err
	}
	return s.transaction(true, func(tx *bbolt.Tx) error { return put(tx, "attachments", item.ID, item) })
}

func (s *Store) ListAttachments() ([]domain.Attachment, error) {
	return list[domain.Attachment](s.db, "attachments")
}

func (s *Store) AttachmentsForRecord(recordID string) ([]domain.Attachment, error) {
	items, err := s.ListAttachments()
	if err != nil {
		return nil, err
	}
	result := make([]domain.Attachment, 0)
	for _, item := range items {
		if item.RecordID == recordID {
			result = append(result, item)
		}
	}
	return result, nil
}

func (s *Store) PutNote(note domain.CollaborationNote) error {
	if err := validateEntityID(note.ID); err != nil {
		return err
	}
	return s.transaction(true, func(tx *bbolt.Tx) error { return put(tx, "notes", note.ID, note) })
}

func (s *Store) ListNotes() ([]domain.CollaborationNote, error) {
	return list[domain.CollaborationNote](s.db, "notes")
}

func (s *Store) NotesForRecord(recordID string) ([]domain.CollaborationNote, error) {
	items, err := s.ListNotes()
	if err != nil {
		return nil, err
	}
	result := make([]domain.CollaborationNote, 0)
	for _, item := range items {
		if item.RecordID == recordID {
			result = append(result, item)
		}
	}
	return result, nil
}
