package handler

import (
	"encoding/json"
	"fmt"
	"github.com/zdarovich/clean-code-example/internal/domain/todo/mock"
	"github.com/zdarovich/clean-code-example/internal/model"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func Test_listTodos(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	m := mock.NewMockServiceInterface(controller)
	r := mux.NewRouter()
	MakeTodoHandlers(r, m)
	path, err := r.GetRoute("listTodos").GetPathTemplate()
	assert.Nil(t, err)
	assert.Equal(t, "/v1/todo", path)
	u := &model.Todo{
		ID: model.NewID(),
	}
	m.EXPECT().ListTodos().
		Return([]*model.Todo{u}, nil)
	ts := httptest.NewServer(listTodos(m))
	defer ts.Close()
	res, err := http.Get(ts.URL)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func Test_listTodos_NotFound(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	m := mock.NewMockServiceInterface(controller)
	ts := httptest.NewServer(listTodos(m))
	defer ts.Close()
	m.EXPECT().
		SearchTodos("dio").
		Return(nil, model.ErrNotFound)
	res, err := http.Get(ts.URL + "?name=dio")
	assert.Nil(t, err)
	assert.Equal(t, http.StatusNotFound, res.StatusCode)
}

func Test_listTodos_Search(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	m := mock.NewMockServiceInterface(controller)
	u := &model.Todo{
		ID: model.NewID(),
	}
	m.EXPECT().
		SearchTodos("ozzy").
		Return([]*model.Todo{u}, nil)
	ts := httptest.NewServer(listTodos(m))
	defer ts.Close()
	res, err := http.Get(ts.URL + "?name=ozzy")
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func Test_createTodo(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	m := mock.NewMockServiceInterface(controller)
	r := mux.NewRouter()
	MakeTodoHandlers(r, m)
	path, err := r.GetRoute("createTodo").GetPathTemplate()
	assert.Nil(t, err)
	assert.Equal(t, "/v1/todo", path)

	m.EXPECT().
		CreateTodo(gomock.Any()).
		Return(model.NewID(), nil)
	h := createTodo(m)

	ts := httptest.NewServer(h)
	defer ts.Close()
	payload := fmt.Sprintf(`{
"title": "ozzy"
}`)
	resp, _ := http.Post(ts.URL+"/v1/todo", "application/json", strings.NewReader(payload))
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var u *TodoResponse
	json.NewDecoder(resp.Body).Decode(&u)
	assert.Equal(t, "ozzy", u.Title)
}

func Test_getTodo(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	m := mock.NewMockServiceInterface(controller)
	r := mux.NewRouter()
	MakeTodoHandlers(r, m)
	path, err := r.GetRoute("getTodo").GetPathTemplate()
	assert.Nil(t, err)
	assert.Equal(t, "/v1/todo/{id}", path)
	u := &model.Todo{
		ID: model.NewID(),
	}
	m.EXPECT().
		GetTodo(u.ID).
		Return(u, nil)
	handler := getTodo(m)
	r.Handle("/v1/todo/{id}", handler)
	ts := httptest.NewServer(r)
	defer ts.Close()
	res, err := http.Get(ts.URL + "/v1/todo/" + u.ID.String())
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	var d *TodoResponse
	json.NewDecoder(res.Body).Decode(&d)
	assert.NotNil(t, d)
	assert.Equal(t, u.ID, d.ID)
}

func Test_deleteTodo(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	m := mock.NewMockServiceInterface(controller)
	r := mux.NewRouter()
	MakeTodoHandlers(r, m)
	path, err := r.GetRoute("deleteTodo").GetPathTemplate()
	assert.Nil(t, err)
	assert.Equal(t, "/v1/todo/{id}", path)
	u := &model.Todo{
		ID: model.NewID(),
	}
	m.EXPECT().DeleteTodo(u.ID).Return(nil)
	handler := deleteTodo(m)
	req, _ := http.NewRequest("DELETE", "/v1/todo/"+u.ID.String(), nil)
	r.Handle("/v1/todo/{id}", handler).Methods("DELETE", "OPTIONS")
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
}
