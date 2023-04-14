package todo

import (
	"fmt"
	"github.com/zdarovich/clean-code-example/internal/model"
	"strings"
)

type inmem struct {
	m map[model.ID]*model.Todo
}

func newInmem() *inmem {
	var m = map[model.ID]*model.Todo{}
	return &inmem{
		m: m,
	}
}

func (r *inmem) Create(e *model.Todo) (model.ID, error) {
	r.m[e.ID] = e
	return e.ID, nil
}

func (r *inmem) Get(id model.ID) (*model.Todo, error) {
	if r.m[id] == nil {
		return nil, model.ErrNotFound
	}
	return r.m[id], nil
}

func (r *inmem) Update(e *model.Todo) error {
	_, err := r.Get(e.ID)
	if err != nil {
		return err
	}
	r.m[e.ID] = e
	return nil
}

func (r *inmem) Search(query string) ([]*model.Todo, error) {
	var d []*model.Todo
	for _, j := range r.m {
		if strings.Contains(strings.ToLower(j.Title), query) {
			d = append(d, j)
		}
	}
	if len(d) == 0 {
		return nil, model.ErrNotFound
	}

	return d, nil
}

func (r *inmem) List() ([]*model.Todo, error) {
	var d []*model.Todo
	for _, j := range r.m {
		d = append(d, j)
	}
	return d, nil
}

func (r *inmem) Delete(id model.ID) error {
	if r.m[id] == nil {
		return fmt.Errorf("not found")
	}
	r.m[id] = nil
	return nil
}
