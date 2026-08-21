package parser

import (
	"strings"
	"testing"
)

func TestParseCSVAndBuildRecords(t *testing.T) {
	doc, err := ParseCSV("B1", "title", "csv", strings.NewReader("id,store,inspector,score,findings\nR1,S1,A,80,clean|safe\n"))
	if err != nil {
		t.Fatal(err)
	}
	records, err := BuildRecords(doc)
	if err != nil || len(records) != 1 || len(records[0].Findings) != 2 {
		t.Fatalf("records %#v %v", records, err)
	}
}
