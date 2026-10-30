package service

import (
	"context"
	"errors"
	"strings"
	"unicode/utf8"

	"example.com/todolistservice/internal/model"
)

var ErrInvalidDescription = errors.New("description must contain 1 to 255 characters")

type Repository interface {
	List(context.Context) ([]model.Item, error)
	Create(context.Context, string) (model.Item, error)
}

type TodoService struct {
	repository Repository
}

func New(repository Repository) *TodoService {
	return &TodoService{repository: repository}
}

func (s *TodoService) List(ctx context.Context) ([]model.Item, error) {
	return s.repository.List(ctx)
}

func (s *TodoService) Create(ctx context.Context, description string) (model.Item, error) {
	description = strings.TrimSpace(description)
	if description == "" || utf8.RuneCountInString(description) > 255 {
		return model.Item{}, ErrInvalidDescription
	}
	return s.repository.Create(ctx, description)
}
