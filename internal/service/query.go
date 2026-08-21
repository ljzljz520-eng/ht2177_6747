package service

import (
	"storeledger/internal/domain"
	"strings"
)

func (s *Service) QueryRecords(query domain.Query) (domain.Page, error) {
	if err := s.ensure(); err != nil {
		return domain.Page{}, err
	}
	records, err := s.Store.ListRecords()
	if err != nil {
		return domain.Page{}, err
	}
	filtered := domain.ApplyQuery(records, query)
	return domain.MakePage(filtered, query), nil
}

func (s *Service) FindByID(recordID string) (domain.InspectionRecord, error) {
	return s.GetRecord(recordID)
}

func (s *Service) SearchText(text string, page, pageSize int) (domain.Page, error) {
	return s.QueryRecords(domain.Query{Text: strings.ToLower(text), Page: page, PageSize: pageSize})
}

func (s *Service) SearchStore(storeID string, page, pageSize int) (domain.Page, error) {
	return s.QueryRecords(domain.Query{StoreID: storeID, Page: page, PageSize: pageSize})
}

func (s *Service) ListBatchRecords(batchID string) ([]domain.InspectionRecord, error) {
	if err := s.ensure(); err != nil {
		return nil, err
	}
	return s.Store.RecordsForBatch(batchID)
}

func Paginate(records []domain.InspectionRecord, page, pageSize int) domain.Page {
	return domain.MakePage(domain.SortRecords(records), domain.Query{Page: page, PageSize: pageSize})
}
