package service

import (
	"storeledger/internal/domain"
)

type Metrics struct {
	Records  int
	Batches  int
	Reviewed int
	Archived int
	Average  float64
	Risks    map[string]int
}

func (s *Service) Metrics() (Metrics, error) {
	records, err := s.Store.ListRecords()
	if err != nil {
		return Metrics{}, err
	}
	batches, err := s.Store.ListBatches()
	if err != nil {
		return Metrics{}, err
	}
	result := Metrics{Records: len(records), Batches: len(batches), Risks: domain.CountFindings(records)}
	for _, record := range records {
		if record.Status == domain.RecordReviewed {
			result.Reviewed++
		}
		if record.Status == domain.RecordArchived {
			result.Archived++
		}
		result.Average += float64(record.Score)
	}
	if result.Records > 0 {
		result.Average /= float64(result.Records)
	}
	return result, nil
}

func (m Metrics) Complete() bool {
	return m.Records >= 0 && m.Batches >= 0 && m.Reviewed >= 0 && m.Archived >= 0
}

func (m Metrics) RiskCount() int {
	total := 0
	for _, count := range m.Risks {
		total += count
	}
	return total
}
