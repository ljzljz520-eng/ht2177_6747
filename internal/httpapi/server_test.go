package httpapi

import (
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"storeledger/internal/service"
	"storeledger/internal/store"
)

func TestHealthEndpoint(t *testing.T) {
	svc := serviceForHTTP(t)
	request := httptest.NewRequest("GET", "/health", nil)
	response := httptest.NewRecorder()
	NewServer(svc).Handler().ServeHTTP(response, request)
	if response.Code != 200 || !strings.Contains(response.Body.String(), "ok") {
		t.Fatalf("response %d %s", response.Code, response.Body.String())
	}
}
func serviceForHTTP(t *testing.T) *service.Service {
	t.Helper()
	st, err := store.Open(filepath.Join(t.TempDir(), "ledger.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { st.Close() })
	svc, err := service.New(st, service.DeterministicClock("2024-01-01T00:00:00Z"), "reviewer")
	if err != nil {
		t.Fatal(err)
	}
	return svc
}
