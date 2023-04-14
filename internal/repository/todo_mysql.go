package repository

import (
	"database/sql"
	"fmt"
	"github.com/zdarovich/clean-code-example/internal/model"
	"time"
)

// TodoMySQL mysql repo
type TodoMySQL struct {
	db *sql.DB
}

// NewTodoMySQL create new repository
func NewTodoMySQL(db *sql.DB) *TodoMySQL {
	return &TodoMySQL{
		db: db,
	}
}

func (r *TodoMySQL) Create(e *model.Todo) (model.ID, error) {
	stmt, err := r.db.Prepare(`
		insert into todo (id, title, completed, created_at) 
		values(?,?,?,?)`)
	if err != nil {
		return e.ID, err
	}
	_, err = stmt.Exec(
		e.ID,
		e.Title,
		e.Completed,
		time.Now().Format("2006-01-02"),
	)
	if err != nil {
		return e.ID, err
	}
	err = stmt.Close()
	if err != nil {
		return e.ID, err
	}
	return e.ID, nil
}

func (r *TodoMySQL) Get(id model.ID) (*model.Todo, error) {
	return getTodo(id, r.db)
}

func getTodo(id model.ID, db *sql.DB) (*model.Todo, error) {
	stmt, err := db.Prepare(`select id, title, completed, created_at from todo where id = ?`)
	if err != nil {
		return nil, err
	}
	var u model.Todo
	rows, err := stmt.Query(id)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		err = rows.Scan(&u.ID, &u.Title, &u.Completed, &u.CreatedAt)
	}
	return &u, nil
}

func (r *TodoMySQL) Update(e *model.Todo) error {
	e.UpdatedAt = time.Now()
	_, err := r.db.Exec("update todo set title = ?, completed = ?, updated_at = ? where id = ?", e.Title, e.Completed, e.UpdatedAt.Format("2006-01-02"), e.ID)
	if err != nil {
		return err
	}
	return nil
}

func (r *TodoMySQL) Search(query string) ([]*model.Todo, error) {
	stmt, err := r.db.Prepare(`select id from todo where title like ?`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	var ids []model.ID
	rows, err := stmt.Query("%" + query + "%")
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var i model.ID
		err = rows.Scan(&i)
		if err != nil {
			return nil, err
		}
		ids = append(ids, i)
	}
	if len(ids) == 0 {
		return nil, fmt.Errorf("not found")
	}
	var todos []*model.Todo
	for _, id := range ids {
		u, err := getTodo(id, r.db)
		if err != nil {
			return nil, err
		}
		todos = append(todos, u)
	}
	return todos, nil
}

func (r *TodoMySQL) List() ([]*model.Todo, error) {
	stmt, err := r.db.Prepare(`select id from todo`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	var ids []model.ID
	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var i model.ID
		err = rows.Scan(&i)
		if err != nil {
			return nil, err
		}
		ids = append(ids, i)
	}
	if len(ids) == 0 {
		return nil, fmt.Errorf("not found")
	}
	var todos []*model.Todo
	for _, id := range ids {
		u, err := getTodo(id, r.db)
		if err != nil {
			return nil, err
		}
		todos = append(todos, u)
	}
	return todos, nil
}

func (r *TodoMySQL) Delete(id model.ID) error {
	_, err := r.db.Exec("delete from todo where id = ?", id)
	if err != nil {
		return err
	}
	return nil
}
