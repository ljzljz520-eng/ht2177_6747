package store

import (
	"path/filepath"
	"testing"

	"storeledger/internal/domain"
)

func TestStoreRoundTrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), "ledger.db")
	st, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	batch, _ := domain.NewBatch("B1", "title", "source")
	record, _ := domain.NewRecord("R1", "S1", "B1", "A", 90, nil, 0)
	if err := st.PutBatch(batch); err != nil {
		t.Fatal(err)
	}
	if err := st.PutRecord(record); err != nil {
		t.Fatal(err)
	}
	loaded, err := st.GetRecord("R1")
	if err != nil || loaded.ID != "R1" {
		t.Fatalf("loaded %#v %v", loaded, err)
	}
}

func TestStoreRejectsClosed(t *testing.T) {
	st, err := Open(filepath.Join(t.TempDir(), "ledger.db"))
	if err != nil {
		t.Fatal(err)
	}
	if err := st.Close(); err != nil {
		t.Fatal(err)
	}
	if _, err := st.GetRecord("R1"); err == nil {
		t.Fatal("expected closed error")
	}
}
