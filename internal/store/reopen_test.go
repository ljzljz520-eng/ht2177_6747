package store

import (
	"path/filepath"
	"testing"

	"storeledger/internal/domain"
)

func TestPersistenceSurvivesReopen(t *testing.T) {
	path := filepath.Join(t.TempDir(), "ledger.db")
	st, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	batch, _ := domain.NewBatch("B-REOPEN", "reopen", "fixture")
	record, _ := domain.NewRecord("R-REOPEN", "S1", "B-REOPEN", "A", 81, []string{"light"}, 0)
	if err := st.PutBatch(batch); err != nil {
		t.Fatal(err)
	}
	if err := st.PutRecord(record); err != nil {
		t.Fatal(err)
	}
	if err := st.PutAudit(domain.AuditEvent{ID: "A-REOPEN", RecordID: record.ID, Reviewer: "R", Decision: "approved"}); err != nil {
		t.Fatal(err)
	}
	if err := st.PutAttachment(domain.Attachment{ID: "F-REOPEN", RecordID: record.ID, Name: "photo", Checksum: "abc"}); err != nil {
		t.Fatal(err)
	}
	if err := st.PutNote(domain.CollaborationNote{ID: "N-REOPEN", RecordID: record.ID, Author: "A", Body: "ok"}); err != nil {
		t.Fatal(err)
	}
	if err := st.Close(); err != nil {
		t.Fatal(err)
	}
	st, err = Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	if _, err := st.GetBatch(batch.BatchID); err != nil {
		t.Fatal(err)
	}
	if _, err := st.GetRecord(record.ID); err != nil {
		t.Fatal(err)
	}
	audits, _ := st.ListAudits()
	attachments, _ := st.ListAttachments()
	notes, _ := st.ListNotes()
	if len(audits) != 1 || len(attachments) != 1 || len(notes) != 1 {
		t.Fatalf("entities did not reopen: %d %d %d", len(audits), len(attachments), len(notes))
	}
}
