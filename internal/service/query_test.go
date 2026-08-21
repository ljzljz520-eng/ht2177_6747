package service

import (
	"storeledger/internal/domain"
	"testing"
)

func TestWorkflowCollaborationQuery(t *testing.T) {
	svc := newTestService(t)
	if _, err := svc.ImportAndValidate(sampleDocument("B-QUERY")); err != nil {
		t.Fatal(err)
	}
	if _, err := svc.AddNote("R1", "collaborator", "check photo"); err != nil {
		t.Fatal(err)
	}
	page, err := svc.SearchText("clean", 1, 10)
	if err != nil {
		t.Fatal(err)
	}
	if page.Total != 1 {
		t.Fatalf("total %d", page.Total)
	}
	text, err := svc.Publish("R1")
	if err != nil {
		t.Fatal(err)
	}
	if text == "" {
		t.Fatal("empty publish")
	}
	_ = domain.Query{}
}
