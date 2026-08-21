package domain

import "testing"

func TestQueryFiltersAndPagination(t *testing.T) {
	records := []InspectionRecord{{ID: "1", StoreID: "S", BatchID: "B", Status: RecordImported, Score: 80, Sequence: 1}, {ID: "2", StoreID: "S", BatchID: "B", Status: RecordReviewed, Score: 90, Sequence: 2}, {ID: "3", StoreID: "X", BatchID: "B", Status: RecordImported, Score: 95, Sequence: 3}}
	filtered := ApplyQuery(records, Query{StoreID: "S", MinScore: 80})
	page := MakePage(filtered, Query{Page: 1, PageSize: 1})
	if page.Total != 2 || len(page.Items) != 1 || page.Items[0].ID != "2" {
		t.Fatalf("unexpected page %#v", page)
	}
}
