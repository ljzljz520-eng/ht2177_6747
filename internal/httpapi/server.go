package httpapi

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"storeledger/internal/domain"
	"storeledger/internal/parser"
	"storeledger/internal/service"
)

type Server struct {
	Service *service.Service
	Mux     *http.ServeMux
}

func NewServer(svc *service.Service) *Server {
	server := &Server{Service: svc, Mux: http.NewServeMux()}
	server.routes()
	return server
}
func (s *Server) routes() {
	s.Mux.HandleFunc("/health", s.handleHealth)
	s.Mux.HandleFunc("/records", s.handleRecords)
	s.Mux.HandleFunc("/batches/import", s.handleImport)
	s.Mux.HandleFunc("/records/review", s.handleReview)
	s.Mux.HandleFunc("/records/note", s.handleNote)
	s.Mux.HandleFunc("/records/publish", s.handlePublish)
}
func (s *Server) Handler() http.Handler { return logging(s.Mux) }

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
func writeError(w http.ResponseWriter, status int, err error) {
	writeJSON(w, status, map[string]string{"error": err.Error()})
}

func (s *Server) handleHealth(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, s.Service.Health())
}

func (s *Server) handleRecords(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, errMethod)
		return
	}
	query := domain.Query{BatchID: r.URL.Query().Get("batch"), StoreID: r.URL.Query().Get("store"), Status: r.URL.Query().Get("status"), Text: r.URL.Query().Get("q")}
	query.Page, _ = strconv.Atoi(r.URL.Query().Get("page"))
	query.PageSize, _ = strconv.Atoi(r.URL.Query().Get("page_size"))
	page, err := s.Service.QueryRecords(query)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, page)
}

func (s *Server) handleImport(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, errMethod)
		return
	}
	batch := r.URL.Query().Get("batch")
	title := r.URL.Query().Get("title")
	source := r.URL.Query().Get("source")
	result, err := parser.ParseCSV(batch, title, source, r.Body)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	imported, err := s.Service.ImportAndValidate(result)
	if err != nil {
		writeError(w, http.StatusUnprocessableEntity, err)
		return
	}
	writeJSON(w, http.StatusCreated, imported)
}

type reviewRequest struct {
	ID      string `json:"id"`
	Note    string `json:"note"`
	Approve bool   `json:"approve"`
}

func (s *Server) handleReview(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, errMethod)
		return
	}
	var request reviewRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	record, event, err := s.Service.ReviewRecord(request.ID, request.Note, request.Approve)
	if err != nil {
		writeError(w, http.StatusUnprocessableEntity, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"record": record, "audit": event})
}

type noteRequest struct {
	ID     string `json:"id"`
	Author string `json:"author"`
	Body   string `json:"body"`
}

func (s *Server) handleNote(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, errMethod)
		return
	}
	var request noteRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	note, err := s.Service.AddNote(request.ID, request.Author, request.Body)
	if err != nil {
		writeError(w, http.StatusUnprocessableEntity, err)
		return
	}
	writeJSON(w, http.StatusCreated, note)
}

func (s *Server) handlePublish(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, errMethod)
		return
	}
	id := strings.TrimSpace(r.URL.Query().Get("id"))
	text, err := s.Service.Publish(id)
	if err != nil {
		writeError(w, http.StatusNotFound, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"summary": text})
}

var errMethod = &methodError{}

type methodError struct{}

func (*methodError) Error() string { return "method not allowed" }
func logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { next.ServeHTTP(w, r) })
}
