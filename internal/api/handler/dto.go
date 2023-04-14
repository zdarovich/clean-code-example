package handler

import "github.com/zdarovich/clean-code-example/internal/model"

type TodoResponse struct {
	ID        model.ID `json:"id"`
	Title     string   `json:"title"`
	Completed bool     `json:"completed"`
}

type CreateTodoReqeust struct {
	Title string `json:"title"`
}
