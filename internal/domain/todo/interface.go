package todo

import (
	"github.com/zdarovich/clean-code-example/internal/model"
)

type Reader interface {
	Get(id model.ID) (*model.Todo, error)
	Search(query string) ([]*model.Todo, error)
	List() ([]*model.Todo, error)
}

type Writer interface {
	Create(e *model.Todo) (model.ID, error)
	Update(e *model.Todo) error
	Delete(id model.ID) error
}

type Repository interface {
	Reader
	Writer
}

// ServiceInterface interface
type ServiceInterface interface {
	GetTodo(id model.ID) (*model.Todo, error)
	SearchTodos(query string) ([]*model.Todo, error)
	ListTodos() ([]*model.Todo, error)
	CreateTodo(title string) (model.ID, error)
	UpdateTodo(e *model.Todo) error
	DeleteTodo(id model.ID) error
}
