package controller

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"example.com/todolistservice/internal/model"
)

type fakeService struct{}

func (fakeService) List(context.Context) ([]model.Item, error) {
	return []model.Item{}, nil
}

func (fakeService) Create(_ context.Context, description string) (model.Item, error) {
	return model.Item{ID: 7, Description: description}, nil
}

type fakeReadiness struct{}

func (fakeReadiness) Ready(context.Context) error { return nil }

func TestCreateResponse(t *testing.T) {
	request := httptest.NewRequest(http.MethodPost, "/api/items", strings.NewReader(`{"description":"pane"}`))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	New(fakeService{}, fakeReadiness{}).ServeHTTP(response, request)
	if response.Code != http.StatusCreated || response.Header().Get("Location") != "/api/items/7" {
		t.Fatalf("status = %d, location = %q", response.Code, response.Header().Get("Location"))
	}
	if !strings.Contains(response.Body.String(), `"description":"pane"`) {
		t.Fatalf("body = %q", response.Body.String())
	}
}

func TestInvalidJSON(t *testing.T) {
	request := httptest.NewRequest(http.MethodPost, "/api/items", strings.NewReader(`{"description":"pane","extra":1}`))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	New(fakeService{}, fakeReadiness{}).ServeHTTP(response, request)
	if response.Code != http.StatusBadRequest {
		t.Fatalf("status = %d", response.Code)
	}
}
