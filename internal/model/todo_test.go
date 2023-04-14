package model_test

import (
	"github.com/zdarovich/clean-code-example/internal/model"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTodoValidate(t *testing.T) {
	type test struct {
		title string
		want  error
	}

	tests := []test{
		{
			title: "test",
			want:  nil,
		},
		{
			title: "",
			want:  model.ErrInvalidEntity,
		},
	}
	for _, tc := range tests {
		_, err := model.NewTodo(tc.title)
		assert.Equal(t, err, tc.want)
	}

}
