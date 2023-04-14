package main

import (
	"database/sql"
	"fmt"
	"github.com/zdarovich/clean-code-example/internal/api/handler"
	"github.com/zdarovich/clean-code-example/internal/api/middleware"
	"github.com/zdarovich/clean-code-example/internal/domain/todo"
	"github.com/zdarovich/clean-code-example/internal/repository"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/zdarovich/clean-code-example/config"
)

func main() {

	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?parseTime=true", config.DB_USER, config.DB_PASSWORD, config.DB_HOST, config.DB_DATABASE)
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		log.Panic(err.Error())
	}
	defer db.Close()

	todoRepo := repository.NewTodoMySQL(db)
	todoService := todo.NewService(todoRepo)

	r := mux.NewRouter()
	r.Use(middleware.Cors)

	handler.MakeTodoHandlers(r, todoService)

	http.Handle("/", r)
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	logger := log.New(os.Stderr, "logger: ", log.Lshortfile)
	srv := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		Addr:         ":" + strconv.Itoa(config.API_PORT),
		Handler:      context.ClearHandler(http.DefaultServeMux),
		ErrorLog:     logger,
	}
	err = srv.ListenAndServe()
	if err != nil {
		log.Panic(err.Error())
	}
}
