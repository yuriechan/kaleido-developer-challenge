package http

import (
	"backend/internal/domain"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type itemService interface {
	GetItem(ctx context.Context, id string) (*domain.Item, error)
	ListItem(ctx context.Context, item *domain.Item) error
	PurchaseItem(ctx context.Context, item *domain.Item) error
}

type Server struct {
	iSvc itemService
}

func New(iSvc itemService) *Server {
	return &Server{
		iSvc: iSvc,
	}
}

func (s *Server) ListItem(w http.ResponseWriter, r *http.Request) {
	var item domain.Item
	s.unmarshalItem(&item, w, r)
	if err := s.iSvc.ListItem(r.Context(), &item); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("500 - Something bad happened!: %s", err.Error())))
	}
}

func (s *Server) PurchaseItem(w http.ResponseWriter, r *http.Request) {
	var item domain.Item
	s.unmarshalItem(&item, w, r)
	if err := s.iSvc.PurchaseItem(r.Context(), &item); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("500 - Something bad happened!: %s", err.Error())))
	}
}

func (s *Server) GetItem(w http.ResponseWriter, r *http.Request) {
	var item domain.Item
	s.unmarshalItem(&item, w, r)
	resp, err := s.iSvc.GetItem(r.Context(), item.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("500 - Something bad happened!: %s", err.Error())))
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (s *Server) unmarshalItem(item *domain.Item, w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - Something bad happened!"))
	}
	if err := json.Unmarshal(b, &item); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - Something bad happened!"))
	}
}
