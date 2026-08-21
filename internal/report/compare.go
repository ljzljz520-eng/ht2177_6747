package report

import (
	"fmt"
	"storeledger/internal/domain"
)

type Comparison struct {
	Left         string         `json:"left"`
	Right        string         `json:"right"`
	ScoreDelta   int            `json:"score_delta"`
	RecordDelta  int            `json:"record_delta"`
	FindingDelta map[string]int `json:"finding_delta"`
}

func Compare(leftID, rightID string, left, right []domain.InspectionRecord) Comparison {
	result := Comparison{Left: leftID, Right: rightID, ScoreDelta: scoreTotal(right) - scoreTotal(left), RecordDelta: len(right) - len(left), FindingDelta: make(map[string]int)}
	for key, value := range domain.CountFindings(right) {
		result.FindingDelta[key] += value
	}
	for key, value := range domain.CountFindings(left) {
		result.FindingDelta[key] -= value
	}
	return result
}
func scoreTotal(records []domain.InspectionRecord) int {
	total := 0
	for _, record := range records {
		total += record.Score
	}
	return total
}
func (c Comparison) Text() string {
	return fmt.Sprintf("%s vs %s score_delta=%d record_delta=%d", c.Left, c.Right, c.ScoreDelta, c.RecordDelta)
}

func FilterActionRequired(records []domain.InspectionRecord) []domain.InspectionRecord {
	result := make([]domain.InspectionRecord, 0)
	for _, record := range records {
		if domain.IsActionRequired(record) {
			result = append(result, record.Clone())
		}
	}
	return result
}
