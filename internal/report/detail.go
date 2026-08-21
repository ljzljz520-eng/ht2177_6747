package report

import (
	"fmt"
	"sort"
	"strings"

	"storeledger/internal/domain"
)

type RecordDetail struct {
	Record      domain.InspectionRecord
	Audits      []domain.AuditEvent
	Notes       []domain.CollaborationNote
	Attachments []domain.Attachment
}

func Detail(record domain.InspectionRecord, audits []domain.AuditEvent, notes []domain.CollaborationNote, attachments []domain.Attachment) RecordDetail {
	return RecordDetail{Record: record.Clone(), Audits: append([]domain.AuditEvent(nil), audits...), Notes: append([]domain.CollaborationNote(nil), notes...), Attachments: append([]domain.Attachment(nil), attachments...)}
}

func (d RecordDetail) Text() string {
	var out strings.Builder
	fmt.Fprintf(&out, "%s store=%s batch=%s status=%s score=%d\n", d.Record.ID, d.Record.StoreID, d.Record.BatchID, d.Record.Status, d.Record.Score)
	if len(d.Record.Findings) > 0 {
		out.WriteString("findings: ")
		out.WriteString(strings.Join(d.Record.Findings, ", "))
		out.WriteByte('\n')
	}
	for _, audit := range d.Audits {
		fmt.Fprintf(&out, "audit %s %s %s\n", audit.Reviewer, audit.Decision, audit.Note)
	}
	for _, note := range d.Notes {
		fmt.Fprintf(&out, "note %s: %s\n", note.Author, note.Body)
	}
	return out.String()
}

func TopStores(records []domain.InspectionRecord, limit int) []string {
	counts := make(map[string]int)
	for _, record := range records {
		counts[record.StoreID]++
	}
	keys := make([]string, 0, len(counts))
	for key := range counts {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool {
		if counts[keys[i]] != counts[keys[j]] {
			return counts[keys[i]] > counts[keys[j]]
		}
		return keys[i] < keys[j]
	})
	if limit < 1 || limit > len(keys) {
		limit = len(keys)
	}
	return keys[:limit]
}

func RiskDistribution(records []domain.InspectionRecord) map[string]int {
	result := make(map[string]int)
	for _, record := range records {
		result[domain.ScoreBand(record.Score)]++
	}
	return result
}

func RenderTable(records []domain.InspectionRecord) string {
	var out strings.Builder
	out.WriteString("ID | Store | Status | Score\n")
	out.WriteString("---|---|---|---\n")
	for _, record := range domain.SortRecords(records) {
		fmt.Fprintf(&out, "%s | %s | %s | %d\n", record.ID, record.StoreID, record.Status, record.Score)
	}
	return out.String()
}
