package service

import (
	"errors"
	"fmt"
	"strings"

	"storeledger/internal/domain"
	"storeledger/internal/report"
	"storeledger/internal/store"
)

type Clock interface{ Now() string }
type StaticClock struct{ Value string }

func (c StaticClock) Now() string { return c.Value }

type Service struct {
	Store    *store.Store
	Clock    Clock
	Reviewer string
	Policy   domain.ReviewPolicy
}

func New(st *store.Store, clock Clock, reviewer string) (*Service, error) {
	if st == nil {
		return nil, errors.New("store is required")
	}
	if clock == nil {
		return nil, errors.New("clock is required")
	}
	if reviewer == "" {
		return nil, errors.New("reviewer is required")
	}
	policy := domain.DefaultReviewPolicy()
	return &Service{Store: st, Clock: clock, Reviewer: reviewer, Policy: policy}, nil
}

func (s *Service) SetPolicy(policy domain.ReviewPolicy) error {
	if err := policy.Validate(); err != nil {
		return err
	}
	s.Policy = policy
	return nil
}
func (s *Service) ensure() error {
	if s == nil || s.Store == nil {
		return errors.New("service unavailable")
	}
	if s.Clock == nil {
		return errors.New("clock unavailable")
	}
	return nil
}

func (s *Service) RegisterBatch(batchID, title, source string) (domain.InspectionBatch, error) {
	if err := s.ensure(); err != nil {
		return domain.InspectionBatch{}, err
	}
	batch, err := domain.NewBatch(batchID, title, source)
	if err != nil {
		return batch, err
	}
	if err := s.Store.PutBatch(batch); err != nil {
		return batch, err
	}
	return batch, nil
}

func (s *Service) GetBatch(batchID string) (domain.InspectionBatch, error) {
	if err := s.ensure(); err != nil {
		return domain.InspectionBatch{}, err
	}
	return s.Store.GetBatch(batchID)
}
func (s *Service) GetRecord(recordID string) (domain.InspectionRecord, error) {
	if err := s.ensure(); err != nil {
		return domain.InspectionRecord{}, err
	}
	return s.Store.GetRecord(recordID)
}

func (s *Service) ValidateRecord(record domain.InspectionRecord) error {
	if err := record.Validate(); err != nil {
		return err
	}
	if strings.TrimSpace(record.Inspector) == "" {
		return errors.New("inspector is required")
	}
	return nil
}
func (s *Service) Report(batchID string) (report.Summary, error) {
	records, err := s.Store.RecordsForBatch(batchID)
	if err != nil {
		return report.Summary{}, err
	}
	return report.Summarize(batchID, records), nil
}

func (s *Service) MustHaveBatch(batchID string) (domain.InspectionBatch, error) {
	batch, err := s.GetBatch(batchID)
	if err != nil {
		return batch, fmt.Errorf("load batch %s: %w", batchID, err)
	}
	return batch, nil
}
