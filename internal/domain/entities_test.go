package domain

import "testing"

func TestRecordValidationAndClone(t *testing.T) {
	record, err := NewRecord("R1", "S1", "B1", "A", 88, []string{"clean"}, 1)
	if err != nil {
		t.Fatal(err)
	}
	clone := record.Clone()
	clone.Findings[0] = "changed"
	if record.Findings[0] == clone.Findings[0] {
		t.Fatal("clone shares findings")
	}
	if err := record.Validate(); err != nil {
		t.Fatal(err)
	}
}

func TestBatchLifecycle(t *testing.T) {
	batch, err := NewBatch("B1", "title", "source")
	if err != nil {
		t.Fatal(err)
	}
	batch, err = batch.Confirm(2)
	if err != nil {
		t.Fatal(err)
	}
	batch, err = batch.Archive()
	if err != nil {
		t.Fatal(err)
	}
	if batch.Status != BatchArchived {
		t.Fatalf("status %s", batch.Status)
	}
}
