package parser

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"storeledger/internal/domain"
)

type Row struct {
	ID        string
	StoreID   string
	Inspector string
	Score     int
	Findings  []string
}
type BatchDocument struct {
	BatchID     string
	Title       string
	Source      string
	Rows        []Row
	Attachments []domain.Attachment
}

func ParseCSV(batchID, title, source string, input io.Reader) (BatchDocument, error) {
	if strings.TrimSpace(batchID) == "" {
		return BatchDocument{}, errors.New("batch id is required")
	}
	reader := csv.NewReader(bufio.NewReader(input))
	reader.TrimLeadingSpace = true
	doc := BatchDocument{BatchID: batchID, Title: title, Source: source}
	line := 0
	for {
		fields, err := reader.Read()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return BatchDocument{}, fmt.Errorf("csv line %d: %w", line+1, err)
		}
		line++
		if len(fields) == 0 || strings.TrimSpace(strings.Join(fields, "")) == "" {
			continue
		}
		if line == 1 && strings.EqualFold(strings.TrimSpace(fields[0]), "id") {
			continue
		}
		row, err := parseRow(fields, line)
		if err != nil {
			return BatchDocument{}, err
		}
		doc.Rows = append(doc.Rows, row)
	}
	if len(doc.Rows) == 0 {
		return BatchDocument{}, errors.New("document contains no rows")
	}
	return doc, nil
}

func parseRow(fields []string, line int) (Row, error) {
	if len(fields) < 4 {
		return Row{}, fmt.Errorf("csv line %d needs id, store, inspector, score", line)
	}
	score, err := strconv.Atoi(strings.TrimSpace(fields[3]))
	if err != nil {
		return Row{}, fmt.Errorf("csv line %d score: %w", line, err)
	}
	findings := make([]string, 0)
	if len(fields) > 4 {
		findings = splitFindings(fields[4])
	}
	return Row{ID: strings.TrimSpace(fields[0]), StoreID: strings.TrimSpace(fields[1]), Inspector: strings.TrimSpace(fields[2]), Score: score, Findings: findings}, nil
}

func splitFindings(value string) []string {
	parts := strings.Split(value, "|")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		if part = strings.TrimSpace(part); part != "" {
			result = append(result, part)
		}
	}
	return result
}

func ValidateDocument(doc BatchDocument) error {
	if doc.BatchID == "" {
		return errors.New("batch id is required")
	}
	seen := make(map[string]bool)
	for _, row := range doc.Rows {
		if row.ID == "" || row.StoreID == "" {
			return errors.New("row identity is required")
		}
		if seen[row.ID] {
			return fmt.Errorf("duplicate row %s", row.ID)
		}
		seen[row.ID] = true
		if row.Score < 0 || row.Score > 100 {
			return fmt.Errorf("row %s score out of range", row.ID)
		}
	}
	return nil
}

func BuildRecords(doc BatchDocument) ([]domain.InspectionRecord, error) {
	if err := ValidateDocument(doc); err != nil {
		return nil, err
	}
	records := make([]domain.InspectionRecord, 0, len(doc.Rows))
	for sequence, row := range doc.Rows {
		record, err := domain.NewRecord(row.ID, row.StoreID, doc.BatchID, row.Inspector, row.Score, row.Findings, sequence)
		if err != nil {
			return nil, err
		}
		records = append(records, record)
	}
	return records, nil
}

func BuildAttachments(doc BatchDocument, records []domain.InspectionRecord) []domain.Attachment {
	result := make([]domain.Attachment, 0, len(doc.Attachments))
	for _, item := range doc.Attachments {
		for _, record := range records {
			if item.RecordID == record.ID {
				result = append(result, item)
			}
		}
	}
	return result
}

func CanonicalRows(records []domain.InspectionRecord) []Row {
	result := make([]Row, 0, len(records))
	for _, record := range records {
		result = append(result, Row{ID: record.ID, StoreID: record.StoreID, Inspector: record.Inspector, Score: record.Score, Findings: append([]string(nil), record.Findings...)})
	}
	return result
}
