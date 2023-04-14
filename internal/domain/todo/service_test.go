package todo

import (
	"github.com/zdarovich/clean-code-example/internal/model"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func newFixtureTodo() *model.Todo {
	return &model.Todo{
		ID:        model.NewID(),
		Title:     "todo1",
		CreatedAt: time.Now(),
	}
}

func Test_Create(t *testing.T) {
	repo := newInmem()
	m := NewService(repo)
	u := newFixtureTodo()
	_, err := m.CreateTodo(u.Title)
	assert.Nil(t, err)
	assert.False(t, u.CreatedAt.IsZero())
	assert.True(t, u.UpdatedAt.IsZero())
}

func Test_SearchAndFind(t *testing.T) {
	repo := newInmem()
	m := NewService(repo)
	u1 := newFixtureTodo()
	u2 := newFixtureTodo()
	u2.Title = "todo2"

	uID, _ := m.CreateTodo(u1.Title)
	_, _ = m.CreateTodo(u2.Title)

	t.Run("search", func(t *testing.T) {
		c, err := m.SearchTodos("todo1")
		assert.Nil(t, err)
		assert.Equal(t, 1, len(c))
		assert.Equal(t, "todo1", c[0].Title)

		c, err = m.SearchTodos("notfound")
		assert.Equal(t, model.ErrNotFound, err)
		assert.Nil(t, c)
	})
	t.Run("list all", func(t *testing.T) {
		all, err := m.ListTodos()
		assert.Nil(t, err)
		assert.Equal(t, 2, len(all))
	})

	t.Run("get", func(t *testing.T) {
		saved, err := m.GetTodo(uID)
		assert.Nil(t, err)
		assert.Equal(t, u1.Title, saved.Title)
	})
}

func Test_Update(t *testing.T) {
	repo := newInmem()
	m := NewService(repo)
	u := newFixtureTodo()
	id, err := m.CreateTodo(u.Title)
	assert.Nil(t, err)
	saved, _ := m.GetTodo(id)
	saved.Title = "Dio"
	assert.Nil(t, m.UpdateTodo(saved))
	updated, err := m.GetTodo(id)
	assert.Nil(t, err)
	assert.Equal(t, "Dio", updated.Title)
	assert.False(t, updated.UpdatedAt.IsZero())
}

func TestDelete(t *testing.T) {
	repo := newInmem()
	m := NewService(repo)
	u1 := newFixtureTodo()
	u2 := newFixtureTodo()
	u2ID, _ := m.CreateTodo(u2.Title)

	err := m.DeleteTodo(u1.ID)
	assert.Equal(t, model.ErrNotFound, err)

	err = m.DeleteTodo(u2ID)
	assert.Nil(t, err)
	_, err = m.GetTodo(u2ID)
	assert.Equal(t, model.ErrNotFound, err)
}
