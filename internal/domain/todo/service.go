package todo

import (
	"github.com/zdarovich/clean-code-example/internal/model"
	"strings"
	"time"
)

type Service struct {
	repo Repository
}

func NewService(r Repository) *Service {
	return &Service{
		repo: r,
	}
}

func (s *Service) CreateTodo(title string) (model.ID, error) {
	e, err := model.NewTodo(title)
	if err != nil {
		return e.ID, err
	}
	return s.repo.Create(e)
}

func (s *Service) GetTodo(id model.ID) (*model.Todo, error) {
	return s.repo.Get(id)
}

func (s *Service) SearchTodos(query string) ([]*model.Todo, error) {
	return s.repo.Search(strings.ToLower(query))
}

func (s *Service) ListTodos() ([]*model.Todo, error) {
	return s.repo.List()
}

func (s *Service) DeleteTodo(id model.ID) error {
	u, err := s.GetTodo(id)
	if u == nil {
		return model.ErrNotFound
	}
	if err != nil {
		return err
	}
	return s.repo.Delete(id)
}

func (s *Service) UpdateTodo(e *model.Todo) error {
	err := e.Validate()
	if err != nil {
		return model.ErrInvalidEntity
	}
	e.UpdatedAt = time.Now()
	return s.repo.Update(e)
}
