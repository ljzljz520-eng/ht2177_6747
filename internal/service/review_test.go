package service

import "testing"

func TestWorkflowReviewArchive(t *testing.T) {
	svc := newTestService(t)
	if _, err := svc.ImportAndValidate(sampleDocument("B-REVIEW")); err != nil {
		t.Fatal(err)
	}
	record, event, err := svc.ReviewRecord("R1", "verified", true)
	if err != nil {
		t.Fatal(err)
	}
	if event.Decision != "approved" || record.Status != "reviewed" {
		t.Fatalf("review %#v %#v", record, event)
	}
	archived, err := svc.ArchiveRecord("R1")
	if err != nil {
		t.Fatal(err)
	}
	if archived.Status != "archived" {
		t.Fatal(archived.Status)
	}
}
