package parser

import (
	"errors"
	"strings"

	"storeledger/internal/domain"
)

func NormalizeDocument(doc BatchDocument) (BatchDocument, error) {
	if err := ValidateDocument(doc); err != nil {
		return BatchDocument{}, err
	}
	normalized := doc
	normalized.BatchID = strings.TrimSpace(doc.BatchID)
	normalized.Title = strings.TrimSpace(doc.Title)
	normalized.Rows = make([]Row, len(doc.Rows))
	for i, row := range doc.Rows {
		normalized.Rows[i] = normalizeRow(row)
	}
	return normalized, nil
}

func normalizeRow(row Row) Row {
	row.ID = strings.TrimSpace(row.ID)
	row.StoreID = strings.TrimSpace(row.StoreID)
	row.Inspector = strings.TrimSpace(row.Inspector)
	for i, finding := range row.Findings {
		row.Findings[i] = strings.TrimSpace(finding)
	}
	return row
}

func ValidateAttachments(doc BatchDocument, records []domain.InspectionRecord) error {
	known := make(map[string]bool, len(records))
	for _, record := range records {
		known[record.ID] = true
	}
	for _, item := range doc.Attachments {
		if item.ID == "" || item.Name == "" || item.Checksum == "" {
			return errors.New("attachment fields are required")
		}
		if !known[item.RecordID] {
			return errors.New("attachment references unknown record")
		}
	}
	return nil
}

func SplitDocuments(rows []Row, size int) [][]Row {
	if size < 1 {
		size = 1
	}
	result := make([][]Row, 0)
	for start := 0; start < len(rows); start += size {
		end := start + size
		if end > len(rows) {
			end = len(rows)
		}
		result = append(result, append([]Row(nil), rows[start:end]...))
	}
	return result
}

func MergeRows(groups [][]Row) []Row {
	total := 0
	for _, group := range groups {
		total += len(group)
	}
	result := make([]Row, 0, total)
	for _, group := range groups {
		result = append(result, group...)
	}
	return result
}
