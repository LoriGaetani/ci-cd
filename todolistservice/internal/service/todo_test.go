package service

import (
	"context"
	"errors"
	"strings"
	"testing"

	"example.com/todolistservice/internal/model"
)

type fakeRepository struct {
	created string
}

func (f *fakeRepository) List(context.Context) ([]model.Item, error) {
	return []model.Item{}, nil
}

func (f *fakeRepository) Create(_ context.Context, description string) (model.Item, error) {
	f.created = description
	return model.Item{ID: 1, Description: description}, nil
}

func TestCreateValidatesAndTrims(t *testing.T) {
	repository := &fakeRepository{}
	svc := New(repository)
	for _, description := range []string{"  ", strings.Repeat("x", 256)} {
		if _, err := svc.Create(context.Background(), description); !errors.Is(err, ErrInvalidDescription) {
			t.Fatalf("Create(%q) error = %v", description, err)
		}
	}
	item, err := svc.Create(context.Background(), "  comprare il pane  ")
	if err != nil || item.Description != "comprare il pane" || repository.created != "comprare il pane" {
		t.Fatalf("item = %+v, created = %q, err = %v", item, repository.created, err)
	}
}
