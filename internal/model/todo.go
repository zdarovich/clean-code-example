package model

import (
	"time"
)

type Todo struct {
	ID        ID
	Title     string
	Completed bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewTodo(title string) (*Todo, error) {
	u := &Todo{
		ID:        NewID(),
		Title:     title,
		CreatedAt: time.Now(),
	}
	err := u.Validate()
	if err != nil {
		return nil, ErrInvalidEntity
	}
	return u, nil
}

func (u *Todo) Complete() error {
	u.Completed = true
	return nil
}

func (u *Todo) Validate() error {
	if u.Title == "" {
		return ErrInvalidEntity
	}

	return nil
}
