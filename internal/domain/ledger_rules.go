package domain

import (
	"errors"
	"strings"
)

type BatchTransition struct {
	From   string
	To     string
	Actor  string
	Reason string
}

func AllowedBatchTransition(from, to string) bool {
	switch from {
	case BatchPending:
		return to == BatchConfirmed
	case BatchConfirmed:
		return to == BatchArchived
	case BatchArchived:
		return false
	default:
		return false
	}
}

func ApplyBatchTransition(batch InspectionBatch, transition BatchTransition) (InspectionBatch, error) {
	if transition.Actor == "" {
		return batch, errors.New("transition actor is required")
	}
	if transition.Reason == "" {
		return batch, errors.New("transition reason is required")
	}
	if transition.From != batch.Status || !AllowedBatchTransition(transition.From, transition.To) {
		return batch, errors.New("batch transition is not allowed")
	}
	batch.Status = transition.To
	return batch, nil
}

func FindingKey(finding string) string { return strings.ToLower(strings.TrimSpace(finding)) }

func CountFindings(records []InspectionRecord) map[string]int {
	result := make(map[string]int)
	for _, record := range records {
		for _, finding := range record.Findings {
			key := FindingKey(finding)
			if key != "" {
				result[key]++
			}
		}
	}
	return result
}

func ScoreBand(score int) string {
	if score >= 90 {
		return "excellent"
	}
	if score >= 70 {
		return "acceptable"
	}
	if score >= 50 {
		return "watch"
	}
	return "critical"
}

func IsActionRequired(record InspectionRecord) bool {
	if record.Score < 70 {
		return true
	}
	return len(record.Findings) > 0
}

func Reviewable(records []InspectionRecord) []InspectionRecord {
	result := make([]InspectionRecord, 0)
	for _, record := range records {
		if record.Status == RecordImported {
			result = append(result, record.Clone())
		}
	}
	return result
}
