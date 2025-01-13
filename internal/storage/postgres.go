package storage

import (
	"context"
	"fmt"
	"time"
	"toDoList/internal/model"

	"github.com/jackc/pgx/v4"
)

const maxRetries = 5               // for db
const retryDelay = 2 * time.Second // between retries

type Storage interface {
	GetTodos() ([]model.ToDo, error)
	GetTodoById(id string) (model.ToDo, error)
	GetTodoImageById(id string) (model.ToDo, error)
	AddTodo(todo model.ToDo) error
	UpdateTodo(id string, todo model.ToDo) error
	UpdateTodoImage(id string, imagePath string) error
	DeleteTodo(id string) error
	Close()
}

type postgresStorage struct {
	conn *pgx.Conn
}

// executes a function with retries on error
func retryWrapper(retries int, retryDelay time.Duration, operation func() error) error {
	var err error
	for retries > 0 {
		err = operation()
		if err == nil {
			return nil
		}
		retries--
		time.Sleep(retryDelay)
		retryDelay *= 2 // delay for each subsequent attempt
	}
	return fmt.Errorf("operation failed after multiple retries: %v", err)
}

func NewPostgresDb(connString string) (*postgresStorage, error) {
	var conn *pgx.Conn
	err := retryWrapper(maxRetries, retryDelay, func() error {
		var err error
		conn, err = pgx.Connect(context.Background(), connString)
		return err
	})
	if err != nil {
		return nil, fmt.Errorf("could not connect to db after multiple retries: %v", err)
	}
	return &postgresStorage{conn: conn}, nil
}

func (s *postgresStorage) AddTodo(todo model.ToDo) error {
	// repeat attempts to execute the SQL query
	return retryWrapper(maxRetries, retryDelay, func() error {
		_, err := s.conn.Exec(context.Background(),
			"INSERT INTO todos (id, title, status, image_path, reminder_time) VALUES ($1, $2, $3, $4, $5)",
			todo.ID, todo.Title, todo.Status, todo.ImagePath, todo.ReminderTime)
		return err
	})
}

func (s *postgresStorage) GetTodos() ([]model.ToDo, error) {
	var todos []model.ToDo
	err := retryWrapper(maxRetries, retryDelay, func() error {
		rows, err := s.conn.Query(context.Background(),
			"SELECT id, title, status, image_path, reminder_time FROM todos")
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var todo model.ToDo
			if err := rows.Scan(&todo.ID, &todo.Title, &todo.Status, &todo.ImagePath, &todo.ReminderTime); err != nil {
				return err
			}
			todos = append(todos, todo)
		}
		return nil
	})
	return todos, err
}

func (s *postgresStorage) GetTodoById(id string) (model.ToDo, error) {
	var todo model.ToDo
	err := retryWrapper(maxRetries, retryDelay, func() error {
		return s.conn.QueryRow(context.Background(),
			"SELECT id, title, status, image_path, reminder_time FROM todos WHERE id = $1", id).
			Scan(&todo.ID, &todo.Title, &todo.Status, &todo.ImagePath, &todo.ReminderTime)
	})
	return todo, err
}

func (s *postgresStorage) GetTodoImageById(id string) (model.ToDo, error) {
	var todo model.ToDo
	err := retryWrapper(maxRetries, retryDelay, func() error {
		return s.conn.QueryRow(context.Background(), "SELECT image_path FROM todos WHERE id = $1", id).
			Scan(&todo.ImagePath)
	})
	return todo, err
}

func (s *postgresStorage) UpdateTodo(id string, todo model.ToDo) error {
	return retryWrapper(maxRetries, retryDelay, func() error {
		_, err := s.conn.Exec(context.Background(), "UPDATE todos SET title = $1, status = $2 WHERE id = $3", todo.Title, todo.Status, id)
		return err
	})
}

func (s *postgresStorage) UpdateTodoImage(id string, imagePath string) error {
	return retryWrapper(maxRetries, retryDelay, func() error {
		_, err := s.conn.Exec(context.Background(), "UPDATE todos SET image_path = $1 WHERE id = $2", imagePath, id)
		return err
	})
}

func (s *postgresStorage) DeleteTodo(id string) error {
	return retryWrapper(maxRetries, retryDelay, func() error {
		_, err := s.conn.Exec(context.Background(), "DELETE FROM todos WHERE id = $1", id)
		return err
	})
}

func (s *postgresStorage) Close() {
	s.conn.Close(context.Background())
}
