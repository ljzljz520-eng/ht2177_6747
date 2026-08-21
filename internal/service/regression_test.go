package service

import "testing"

func TestBusiness07Regression(t *testing.T) {
	svc := newTestService(t)
	doc := sampleDocument("RB2177-07")
	doc.Rows[0].Score = 80
	doc.Rows[1].Score = 80
	result, err := svc.ImportAndValidate(doc)
	if err != nil {
		t.Fatal(err)
	}
	if result.Records[0].Sequence != 0 {
		t.Fatalf("same-score records must retain import sequence, got %d", result.Records[0].Sequence)
	}
	second, err := svc.ReimportAndValidate(doc)
	if err != nil {
		t.Fatal(err)
	}
	if second.Records[0].Sequence != 0 {
		t.Fatalf("reimport changed same-score order: %d", second.Records[0].Sequence)
	}
}
