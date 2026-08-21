package domain

import (
	"errors"
	"fmt"
)

type LedgerLink struct {
	RecordID string `json:"record_id"`
	Previous string `json:"previous"`
	Next     string `json:"next"`
	Position int    `json:"position"`
}

type RecordLedger struct {
	BatchID string       `json:"batch_id"`
	Head    string       `json:"head"`
	Tail    string       `json:"tail"`
	Links   []LedgerLink `json:"links"`
}

func NewLedger(batchID string) RecordLedger {
	return RecordLedger{BatchID: batchID, Links: make([]LedgerLink, 0)}
}

func (l RecordLedger) Validate() error {
	if l.BatchID == "" {
		return errors.New("ledger batch id is required")
	}
	if len(l.Links) == 0 {
		if l.Head != "" || l.Tail != "" {
			return errors.New("empty ledger has endpoints")
		}
		return nil
	}
	if l.Head == "" || l.Tail == "" {
		return errors.New("ledger endpoints are required")
	}
	seen := make(map[string]bool, len(l.Links))
	positions := make(map[int]bool, len(l.Links))
	for _, link := range l.Links {
		if link.RecordID == "" {
			return errors.New("ledger link record is required")
		}
		if seen[link.RecordID] {
			return fmt.Errorf("duplicate ledger record %s", link.RecordID)
		}
		seen[link.RecordID] = true
		if positions[link.Position] {
			return errors.New("duplicate ledger position")
		}
		positions[link.Position] = true
	}
	return nil
}

func (l RecordLedger) Add(recordID string) (RecordLedger, error) {
	if recordID == "" {
		return l, errors.New("record id is required")
	}
	for _, link := range l.Links {
		if link.RecordID == recordID {
			return l, errors.New("record already in ledger")
		}
	}
	updated := cloneLedger(l)
	position := len(updated.Links)
	link := LedgerLink{RecordID: recordID, Position: position}
	if position > 0 {
		link.Previous = updated.Links[position-1].RecordID
		updated.Links[position-1].Next = recordID
	} else {
		updated.Head = recordID
	}
	updated.Tail = recordID
	updated.Links = append(updated.Links, link)
	return updated, nil
}

func (l RecordLedger) InsertBefore(target, recordID string) (RecordLedger, error) {
	if target == "" || recordID == "" {
		return l, errors.New("target and record ids are required")
	}
	if target == recordID {
		return l, errors.New("record cannot precede itself")
	}
	index := l.indexOf(target)
	if index < 0 {
		return l, errors.New("target not found")
	}
	if l.indexOf(recordID) >= 0 {
		return l, errors.New("record already in ledger")
	}
	updated := cloneLedger(l)
	updated.Links = append(updated.Links[:index], append([]LedgerLink{{RecordID: recordID}}, updated.Links[index:]...)...)
	updated.relink()
	return updated, nil
}

func (l RecordLedger) Remove(recordID string) (RecordLedger, error) {
	index := l.indexOf(recordID)
	if index < 0 {
		return l, errors.New("record not found")
	}
	updated := cloneLedger(l)
	updated.Links = append(updated.Links[:index], updated.Links[index+1:]...)
	updated.relink()
	return updated, nil
}

func (l RecordLedger) indexOf(recordID string) int {
	for index, link := range l.Links {
		if link.RecordID == recordID {
			return index
		}
	}
	return -1
}

func (l *RecordLedger) relink() {
	l.Head = ""
	l.Tail = ""
	for index := range l.Links {
		l.Links[index].Position = index
		l.Links[index].Previous = ""
		l.Links[index].Next = ""
		if index > 0 {
			l.Links[index].Previous = l.Links[index-1].RecordID
			l.Links[index-1].Next = l.Links[index].RecordID
		}
	}
	if len(l.Links) > 0 {
		l.Head = l.Links[0].RecordID
		l.Tail = l.Links[len(l.Links)-1].RecordID
	}
}

func cloneLedger(l RecordLedger) RecordLedger {
	clone := l
	clone.Links = append([]LedgerLink(nil), l.Links...)
	return clone
}

func (l RecordLedger) OrderedIDs() []string {
	result := make([]string, len(l.Links))
	for i, link := range l.Links {
		result[i] = link.RecordID
	}
	return result
}

func (l RecordLedger) Previous(recordID string) (string, bool) {
	index := l.indexOf(recordID)
	if index <= 0 {
		return "", false
	}
	return l.Links[index-1].RecordID, true
}
func (l RecordLedger) Next(recordID string) (string, bool) {
	index := l.indexOf(recordID)
	if index < 0 || index+1 >= len(l.Links) {
		return "", false
	}
	return l.Links[index+1].RecordID, true
}
func (l RecordLedger) Contains(recordID string) bool { return l.indexOf(recordID) >= 0 }
func (l RecordLedger) Length() int                   { return len(l.Links) }

func BuildLedger(batchID string, records []InspectionRecord) (RecordLedger, error) {
	ledger := NewLedger(batchID)
	ordered := SortRecords(records)
	for _, record := range ordered {
		if record.BatchID != batchID {
			return ledger, errors.New("record belongs to another batch")
		}
		var err error
		ledger, err = ledger.Add(record.ID)
		if err != nil {
			return ledger, err
		}
	}
	return ledger, ledger.Validate()
}

func ReconcileLedger(l RecordLedger, records []InspectionRecord) error {
	if err := l.Validate(); err != nil {
		return err
	}
	if len(records) != len(l.Links) {
		return errors.New("ledger and records have different lengths")
	}
	for _, record := range records {
		if !l.Contains(record.ID) {
			return fmt.Errorf("record %s missing from ledger", record.ID)
		}
	}
	return nil
}
