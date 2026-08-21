package domain

import (
	"errors"
	"fmt"
)

type ReviewPolicy struct {
	MinimumScore       int
	RequireFindingNote bool
}

func DefaultReviewPolicy() ReviewPolicy {
	return ReviewPolicy{MinimumScore: 70, RequireFindingNote: true}
}

func (p ReviewPolicy) Validate() error {
	if p.MinimumScore < 0 || p.MinimumScore > 100 {
		return errors.New("minimum score is invalid")
	}
	return nil
}

func (p ReviewPolicy) CanApprove(record InspectionRecord, note string) error {
	if err := p.Validate(); err != nil {
		return err
	}
	if record.Status != RecordImported && record.Status != RecordReviewed {
		return fmt.Errorf("record %s is not reviewable", record.ID)
	}
	if record.Score < p.MinimumScore {
		return fmt.Errorf("score %d is below policy threshold %d", record.Score, p.MinimumScore)
	}
	if p.RequireFindingNote && len(record.Findings) > 0 && note == "" {
		return errors.New("finding note is required")
	}
	return nil
}

func Approve(record InspectionRecord, reviewer, note, timestamp string) (InspectionRecord, AuditEvent, error) {
	if reviewer == "" {
		return record, AuditEvent{}, errors.New("reviewer is required")
	}
	record.Status = RecordReviewed
	event := AuditEvent{ID: record.ID + ":" + reviewer + ":approve", RecordID: record.ID, Reviewer: reviewer, Decision: "approved", Note: note, CreatedAt: timestamp}
	return record, event, nil
}

func Reject(record InspectionRecord, reviewer, note, timestamp string) (InspectionRecord, AuditEvent, error) {
	if reviewer == "" {
		return record, AuditEvent{}, errors.New("reviewer is required")
	}
	record.Status = RecordImported
	event := AuditEvent{ID: record.ID + ":" + reviewer + ":reject", RecordID: record.ID, Reviewer: reviewer, Decision: "rejected", Note: note, CreatedAt: timestamp}
	return record, event, nil
}

func CanArchive(record InspectionRecord) bool { return record.Status == RecordReviewed }

func ArchiveRecord(record InspectionRecord) (InspectionRecord, error) {
	if !CanArchive(record) {
		return record, errors.New("record must be reviewed before archive")
	}
	record.Status = RecordArchived
	return record, nil
}
