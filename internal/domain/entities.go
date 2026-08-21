package domain

import (
	"errors"
	"fmt"
	"sort"
	"strings"
)

const (
	BatchPending   = "pending"
	BatchConfirmed = "confirmed"
	BatchArchived  = "archived"
	RecordImported = "imported"
	RecordReviewed = "reviewed"
	RecordArchived = "archived"
)

type InspectionRecord struct {
	ID         string   `json:"id"`
	StoreID    string   `json:"store_id"`
	BatchID    string   `json:"batch_id"`
	Inspector  string   `json:"inspector"`
	Status     string   `json:"status"`
	Score      int      `json:"score"`
	Findings   []string `json:"findings"`
	ImportedAt string   `json:"imported_at"`
	Sequence   int      `json:"sequence"`
}

type AuditEvent struct {
	ID        string `json:"id"`
	RecordID  string `json:"record_id"`
	Reviewer  string `json:"reviewer"`
	Decision  string `json:"decision"`
	Note      string `json:"note"`
	CreatedAt string `json:"created_at"`
}

type InspectionBatch struct {
	ID             string `json:"id"`
	BatchID        string `json:"batch_id"`
	Title          string `json:"title"`
	Status         string `json:"status"`
	Source         string `json:"source"`
	RecordCount    int    `json:"record_count"`
	ConfirmedCount int    `json:"confirmed_count"`
}

type Attachment struct {
	ID       string `json:"id"`
	RecordID string `json:"record_id"`
	Name     string `json:"name"`
	Checksum string `json:"checksum"`
	Kind     string `json:"kind"`
}

type CollaborationNote struct {
	ID        string `json:"id"`
	RecordID  string `json:"record_id"`
	Author    string `json:"author"`
	Body      string `json:"body"`
	CreatedAt string `json:"created_at"`
}

type Query struct {
	BatchID  string
	StoreID  string
	Status   string
	MinScore int
	Text     string
	Page     int
	PageSize int
}

type Page struct {
	Items    []InspectionRecord `json:"items"`
	Page     int                `json:"page"`
	PageSize int                `json:"page_size"`
	Total    int                `json:"total"`
	Pages    int                `json:"pages"`
}

func NewBatch(batchID, title, source string) (InspectionBatch, error) {
	batchID = strings.TrimSpace(batchID)
	if batchID == "" {
		return InspectionBatch{}, errors.New("batch id is required")
	}
	if title == "" {
		title = "Untitled inspection batch"
	}
	return InspectionBatch{ID: batchID, BatchID: batchID, Title: title, Status: BatchPending, Source: source}, nil
}

func NewRecord(id, storeID, batchID, inspector string, score int, findings []string, sequence int) (InspectionRecord, error) {
	if strings.TrimSpace(id) == "" || strings.TrimSpace(storeID) == "" || strings.TrimSpace(batchID) == "" {
		return InspectionRecord{}, errors.New("record identity is required")
	}
	if score < 0 || score > 100 {
		return InspectionRecord{}, fmt.Errorf("score %d outside 0..100", score)
	}
	clean := make([]string, 0, len(findings))
	for _, finding := range findings {
		finding = strings.TrimSpace(finding)
		if finding != "" {
			clean = append(clean, finding)
		}
	}
	return InspectionRecord{ID: id, StoreID: storeID, BatchID: batchID, Inspector: inspector, Status: RecordImported, Score: score, Findings: clean, Sequence: sequence}, nil
}

func (r InspectionRecord) Validate() error {
	if r.ID == "" || r.StoreID == "" || r.BatchID == "" {
		return errors.New("record identity is incomplete")
	}
	if r.Score < 0 || r.Score > 100 {
		return errors.New("record score is invalid")
	}
	if r.Status != RecordImported && r.Status != RecordReviewed && r.Status != RecordArchived {
		return fmt.Errorf("unknown record status %q", r.Status)
	}
	return nil
}

func (r InspectionRecord) SearchText() string {
	return strings.ToLower(strings.Join(append([]string{r.ID, r.StoreID, r.BatchID, r.Inspector}, r.Findings...), " "))
}

func (r InspectionRecord) Clone() InspectionRecord {
	copyRecord := r
	copyRecord.Findings = append([]string(nil), r.Findings...)
	return copyRecord
}

func (b InspectionBatch) Confirm(count int) (InspectionBatch, error) {
	if count < 0 {
		return b, errors.New("negative record count")
	}
	if b.Status == BatchArchived {
		return b, errors.New("archived batch cannot be confirmed")
	}
	b.Status = BatchConfirmed
	b.RecordCount = count
	b.ConfirmedCount = count
	return b, nil
}

func (b InspectionBatch) Archive() (InspectionBatch, error) {
	if b.Status != BatchConfirmed {
		return b, errors.New("only confirmed batch can be archived")
	}
	b.Status = BatchArchived
	return b, nil
}

func SortRecords(records []InspectionRecord) []InspectionRecord {
	out := make([]InspectionRecord, len(records))
	for i, record := range records {
		out[i] = record.Clone()
	}
	sort.SliceStable(out, func(i, j int) bool {
		if out[i].Score != out[j].Score {
			return out[i].Score > out[j].Score
		}
		return out[i].Sequence < out[j].Sequence
	})
	return out
}

func StatusAllowed(status string) bool {
	switch status {
	case "", RecordImported, RecordReviewed, RecordArchived:
		return true
	default:
		return false
	}
}

func NormalizeQuery(q Query) Query {
	if q.Page < 1 {
		q.Page = 1
	}
	if q.PageSize < 1 {
		q.PageSize = 20
	}
	if q.PageSize > 100 {
		q.PageSize = 100
	}
	q.BatchID = strings.TrimSpace(q.BatchID)
	q.StoreID = strings.TrimSpace(q.StoreID)
	q.Status = strings.TrimSpace(strings.ToLower(q.Status))
	q.Text = strings.TrimSpace(strings.ToLower(q.Text))
	return q
}

func ApplyQuery(records []InspectionRecord, q Query) []InspectionRecord {
	q = NormalizeQuery(q)
	filtered := make([]InspectionRecord, 0, len(records))
	for _, record := range records {
		if q.BatchID != "" && record.BatchID != q.BatchID {
			continue
		}
		if q.StoreID != "" && record.StoreID != q.StoreID {
			continue
		}
		if q.Status != "" && record.Status != q.Status {
			continue
		}
		if record.Score < q.MinScore {
			continue
		}
		if q.Text != "" && !strings.Contains(record.SearchText(), q.Text) {
			continue
		}
		filtered = append(filtered, record.Clone())
	}
	return SortRecords(filtered)
}

func MakePage(records []InspectionRecord, q Query) Page {
	q = NormalizeQuery(q)
	total := len(records)
	pages := 0
	if total > 0 {
		pages = (total + q.PageSize - 1) / q.PageSize
	}
	start := (q.Page - 1) * q.PageSize
	if start > total {
		start = total
	}
	end := start + q.PageSize
	if end > total {
		end = total
	}
	items := make([]InspectionRecord, end-start)
	copy(items, records[start:end])
	return Page{Items: items, Page: q.Page, PageSize: q.PageSize, Total: total, Pages: pages}
}

func ValidateBatchConsistency(batch InspectionBatch, records []InspectionRecord) error {
	if batch.BatchID == "" {
		return errors.New("batch id is required")
	}
	if batch.RecordCount != len(records) {
		return fmt.Errorf("batch count %d does not match records %d", batch.RecordCount, len(records))
	}
	for _, record := range records {
		if record.BatchID != batch.BatchID {
			return errors.New("record belongs to another batch")
		}
		if err := record.Validate(); err != nil {
			return err
		}
	}
	return nil
}
