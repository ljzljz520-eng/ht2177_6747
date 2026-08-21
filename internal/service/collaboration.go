package service

import (
	"errors"
	"storeledger/internal/domain"
	"strings"
)

func (s *Service) AddNote(recordID, author, body string) (domain.CollaborationNote, error) {
	if err := s.ensure(); err != nil {
		return domain.CollaborationNote{}, err
	}
	if _, err := s.GetRecord(recordID); err != nil {
		return domain.CollaborationNote{}, err
	}
	if strings.TrimSpace(author) == "" || strings.TrimSpace(body) == "" {
		return domain.CollaborationNote{}, errors.New("author and body are required")
	}
	note := domain.CollaborationNote{ID: recordID + ":" + author + ":" + s.Clock.Now(), RecordID: recordID, Author: author, Body: strings.TrimSpace(body), CreatedAt: s.Clock.Now()}
	return note, s.Store.PutNote(note)
}

func (s *Service) Notes(recordID string) ([]domain.CollaborationNote, error) {
	if err := s.ensure(); err != nil {
		return nil, err
	}
	return s.Store.NotesForRecord(recordID)
}

func (s *Service) Attach(recordID, name, checksum, kind string) (domain.Attachment, error) {
	if err := s.ensure(); err != nil {
		return domain.Attachment{}, err
	}
	if _, err := s.GetRecord(recordID); err != nil {
		return domain.Attachment{}, err
	}
	if name == "" || checksum == "" {
		return domain.Attachment{}, errors.New("attachment name and checksum are required")
	}
	item := domain.Attachment{ID: recordID + ":" + name, RecordID: recordID, Name: name, Checksum: checksum, Kind: kind}
	return item, s.Store.PutAttachment(item)
}

func (s *Service) Publish(recordID string) (string, error) {
	record, err := s.GetRecord(recordID)
	if err != nil {
		return "", err
	}
	notes, err := s.Notes(recordID)
	if err != nil {
		return "", err
	}
	return publishText(record, notes), nil
}

func publishText(record domain.InspectionRecord, notes []domain.CollaborationNote) string {
	var out strings.Builder
	out.WriteString(record.ID)
	out.WriteString(" ")
	out.WriteString(record.Status)
	out.WriteString(" score=")
	out.WriteString(strconvInt(record.Score))
	for _, note := range notes {
		out.WriteString(" | ")
		out.WriteString(note.Author)
		out.WriteString(": ")
		out.WriteString(note.Body)
	}
	return out.String()
}
func strconvInt(value int) string {
	if value == 0 {
		return "0"
	}
	negative := value < 0
	if negative {
		value = -value
	}
	digits := make([]byte, 0, 4)
	for value > 0 {
		digits = append(digits, byte('0'+value%10))
		value /= 10
	}
	for i, j := 0, len(digits)-1; i < j; i, j = i+1, j-1 {
		digits[i], digits[j] = digits[j], digits[i]
	}
	if negative {
		digits = append([]byte{'-'}, digits...)
	}
	return string(digits)
}
