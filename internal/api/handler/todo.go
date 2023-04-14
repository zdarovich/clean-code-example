package handler

import (
	"encoding/json"
	"github.com/zdarovich/clean-code-example/internal/domain/todo"
	"github.com/zdarovich/clean-code-example/internal/model"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func listTodos(service todo.ServiceInterface) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		errorMessage := "Error reading todos"
		var data []*model.Todo
		var err error
		title := r.URL.Query().Get("title")
		switch {
		case title == "":
			data, err = service.ListTodos()
		default:
			data, err = service.SearchTodos(title)
		}
		w.Header().Set("Content-Type", "application/json")
		if err != nil && err != model.ErrNotFound {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage))
			return
		}

		if data == nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(errorMessage))
			return
		}
		var toJ []*TodoResponse
		for _, d := range data {
			toJ = append(toJ, &TodoResponse{
				ID:        d.ID,
				Title:     d.Title,
				Completed: d.Completed,
			})
		}
		if err := json.NewEncoder(w).Encode(toJ); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage))
		}
	})
}

func createTodo(service todo.ServiceInterface) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		errorMessage := "Error adding todo"
		var input CreateTodoReqeust
		err := json.NewDecoder(r.Body).Decode(&input)
		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage))
			return
		}
		id, err := service.CreateTodo(input.Title)
		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage))
			return
		}
		toJ := &TodoResponse{
			ID:    id,
			Title: input.Title,
		}

		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(toJ); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage))
			return
		}
	})
}

func getTodo(service todo.ServiceInterface) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		errorMessage := "Error reading todo"
		data, done := getTodoById(w, r, errorMessage, service)
		if done {
			return
		}
		toJ := &TodoResponse{
			ID:        data.ID,
			Title:     data.Title,
			Completed: data.Completed,
		}
		if err := json.NewEncoder(w).Encode(toJ); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage))
		}
	})
}

func getTodoById(w http.ResponseWriter, r *http.Request, errorMessage string, service todo.ServiceInterface) (*model.Todo, bool) {
	vars := mux.Vars(r)
	id, err := model.StringToID(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(errorMessage))
		return nil, true
	}
	data, err := service.GetTodo(id)
	w.Header().Set("Content-Type", "application/json")
	if err != nil && err != model.ErrNotFound {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(errorMessage))
		return nil, true
	}

	if data == nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(errorMessage))
		return nil, true
	}
	return data, false
}

func deleteTodo(service todo.ServiceInterface) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		errorMessage := "Error removing todo"
		vars := mux.Vars(r)
		id, err := model.StringToID(vars["id"])
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage))
			return
		}
		err = service.DeleteTodo(id)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage))
			return
		}
	})
}

func completeTodo(service todo.ServiceInterface) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		errorMessage := "Error completing todo"
		data, done := getTodoById(w, r, errorMessage, service)
		if done {
			return
		}
		var input struct {
			Completed bool `json:"completed"`
		}
		err := json.NewDecoder(r.Body).Decode(&input)
		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage))
			return
		}
		data.Completed = input.Completed
		err = service.UpdateTodo(data)
		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage))
			return
		}
		toJ := &TodoResponse{
			ID:        data.ID,
			Title:     data.Title,
			Completed: data.Completed,
		}

		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(toJ); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage))
			return
		}
	})
}

func MakeTodoHandlers(r *mux.Router, service todo.ServiceInterface) {
	r.Handle("/v1/todo", listTodos(service)).Methods("GET", "OPTIONS").Name("listTodos")

	r.Handle("/v1/todo", createTodo(service)).Methods("POST", "OPTIONS").Name("createTodo")

	r.Handle("/v1/todo/{id}", getTodo(service)).Methods("GET", "OPTIONS").Name("getTodo")

	r.Handle("/v1/todo/{id}", deleteTodo(service)).Methods("DELETE", "OPTIONS").Name("deleteTodo")

	r.Handle("/v1/todo/{id}", completeTodo(service)).Methods("PATCH", "OPTIONS").Name("completeTodo")
}
