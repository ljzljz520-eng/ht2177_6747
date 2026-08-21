package report

import (
	"storeledger/internal/domain"
	"testing"
)

func TestSummaryAndCSV(t *testing.T) {
	records := []domain.InspectionRecord{{ID: "R1", BatchID: "B1", Score: 90, Status: domain.RecordReviewed, Findings: []string{"clean"}}}
	summary := Summarize("B1", records)
	if summary.AverageScore != 90 || !stringsContains(summary.Text(), "reviewed=1") {
		t.Fatal(summary)
	}
	if CSV(records) == "" {
		t.Fatal("missing csv")
	}
}
func stringsContains(value, part string) bool {
	for i := 0; i+len(part) <= len(value); i++ {
		if value[i:i+len(part)] == part {
			return true
		}
	}
	return false
}
