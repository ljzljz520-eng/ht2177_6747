package service

import (
	"time"
)

func Health(st interface{ Healthy() bool }) map[string]string {
	if st == nil {
		return map[string]string{"status": "down"}
	}
	if !st.Healthy() {
		return map[string]string{"status": "down"}
	}
	return map[string]string{"status": "ok"}
}

func DeterministicClock(value string) StaticClock {
	if value == "" {
		value = time.Unix(0, 0).UTC().Format(time.RFC3339)
	}
	return StaticClock{Value: value}
}

func (s *Service) Health() map[string]string {
	if s == nil || s.Store == nil {
		return map[string]string{"status": "down"}
	}
	return Health(s.Store)
}
