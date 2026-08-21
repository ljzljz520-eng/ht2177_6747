package report

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"storeledger/internal/domain"
)

type Summary struct {
	BatchID      string         `json:"batch_id"`
	Total        int            `json:"total"`
	Reviewed     int            `json:"reviewed"`
	Archived     int            `json:"archived"`
	AverageScore float64        `json:"average_score"`
	Findings     map[string]int `json:"findings"`
}

func Summarize(batchID string, records []domain.InspectionRecord) Summary {
	result := Summary{BatchID: batchID, Total: len(records), Findings: make(map[string]int)}
	var total int
	for _, record := range records {
		total += record.Score
		if record.Status == domain.RecordReviewed {
			result.Reviewed++
		}
		if record.Status == domain.RecordArchived {
			result.Archived++
		}
		for _, finding := range record.Findings {
			result.Findings[finding]++
		}
	}
	if len(records) > 0 {
		result.AverageScore = float64(total) / float64(len(records))
	}
	return result
}

func (s Summary) Marshal() ([]byte, error) { return json.MarshalIndent(s, "", "  ") }

func (s Summary) Text() string {
	keys := make([]string, 0, len(s.Findings))
	for key := range s.Findings {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	var out strings.Builder
	fmt.Fprintf(&out, "batch=%s total=%d reviewed=%d archived=%d average=%.2f", s.BatchID, s.Total, s.Reviewed, s.Archived, s.AverageScore)
	for _, key := range keys {
		fmt.Fprintf(&out, " finding[%s]=%d", key, s.Findings[key])
	}
	return out.String()
}

func CSV(records []domain.InspectionRecord) string {
	var out strings.Builder
	out.WriteString("id,store_id,batch_id,status,score,sequence\n")
	for _, record := range domain.SortRecords(records) {
		fmt.Fprintf(&out, "%s,%s,%s,%s,%d,%d\n", record.ID, record.StoreID, record.BatchID, record.Status, record.Score, record.Sequence)
	}
	return out.String()
}
