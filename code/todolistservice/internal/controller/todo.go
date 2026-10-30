package controller

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"mime"
	"net/http"
	"time"

	"example.com/todolistservice/internal/model"
	"example.com/todolistservice/internal/service"
)

type TodoService interface {
	List(context.Context) ([]model.Item, error)
	Create(context.Context, string) (model.Item, error)
}

type ReadinessChecker interface {
	Ready(context.Context) error
}

type Controller struct {
	service   TodoService
	readiness ReadinessChecker
}

func New(service TodoService, readiness ReadinessChecker) http.Handler {
	c := Controller{service: service, readiness: readiness}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/items", c.listItems)
	mux.HandleFunc("POST /api/items", c.createItem)
	mux.HandleFunc("GET /health/live", c.liveness)
	mux.HandleFunc("GET /health/ready", c.ready)
	return mux
}

func (c Controller) listItems(w http.ResponseWriter, r *http.Request) {
	items, err := c.service.List(r.Context())
	if err != nil {
		serverError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, items)
}

func (c Controller) createItem(w http.ResponseWriter, r *http.Request) {
	mediaType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil || mediaType != "application/json" {
		http.Error(w, "Content-Type must be application/json", http.StatusUnsupportedMediaType)
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, 4096)
	var input struct {
		Description string `json:"description"`
	}
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&input); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		http.Error(w, "invalid trailing JSON", http.StatusBadRequest)
		return
	}
	item, err := c.service.Create(r.Context(), input.Description)
	if errors.Is(err, service.ErrInvalidDescription) {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err != nil {
		serverError(w, err)
		return
	}
	w.Header().Set("Location", fmt.Sprintf("/api/items/%d", item.ID))
	writeJSON(w, http.StatusCreated, item)
}

func (c Controller) liveness(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (c Controller) ready(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()
	if err := c.readiness.Ready(ctx); err != nil {
		http.Error(w, "database unavailable", http.StatusServiceUnavailable)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(value); err != nil {
		slog.Error("response write failed", "error", err)
	}
}

func serverError(w http.ResponseWriter, err error) {
	slog.Error("request failed", "error", err)
	http.Error(w, "internal server error", http.StatusInternalServerError)
}
