package service

import (
	"path/filepath"
	"testing"

	"storeledger/internal/parser"
	"storeledger/internal/store"
)

func newTestService(t *testing.T) *Service {
	t.Helper()
	st, err := store.Open(filepath.Join(t.TempDir(), "ledger.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { st.Close() })
	svc, err := New(st, DeterministicClock("2024-01-01T00:00:00Z"), "reviewer")
	if err != nil {
		t.Fatal(err)
	}
	return svc
}
func sampleDocument(batch string) parser.BatchDocument {
	return parser.BatchDocument{BatchID: batch, Title: "sample", Source: "test", Rows: []parser.Row{{ID: "R1", StoreID: "S1", Inspector: "A", Score: 90, Findings: []string{"clean"}}, {ID: "R2", StoreID: "S1", Inspector: "B", Score: 80}}}
}

func TestWorkflowImportValidateReport(t *testing.T) {
	svc := newTestService(t)
	result, err := svc.ImportAndValidate(sampleDocument("B-IMPORT"))
	if err != nil {
		t.Fatal(err)
	}
	if result.Batch.Status != "confirmed" || len(result.Records) != 2 {
		t.Fatalf("unexpected import %#v", result)
	}
	if result.Summary == "" {
		t.Fatal("missing summary")
	}
}
