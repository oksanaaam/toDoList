package storage

import (
	"context"
	"fmt"
	"toDoList/internal/model"

	"github.com/jackc/pgx/v4"
)

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

func NewPostgresConnection(connString string) (*postgresStorage, error) {
	conn, err := pgx.Connect(context.Background(), connString)
	if err != nil {
		return nil, fmt.Errorf("could not connect to db: %v", err)
	}
	return &postgresStorage{conn: conn}, nil
}

func (s *postgresStorage) AddTodo(todo model.ToDo) error {
	_, err := s.conn.Exec(context.Background(),
		"INSERT INTO todos (id, title, status, image_path, reminder_time) VALUES ($1, $2, $3, $4, $5)",
		todo.ID, todo.Title, todo.Status, todo.ImagePath, todo.ReminderTime)
	return err
}

func (s *postgresStorage) GetTodos() ([]model.ToDo, error) {
	rows, err := s.conn.Query(context.Background(),
		"SELECT id, title, status, image_path, reminder_time FROM todos")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []model.ToDo
	for rows.Next() {
		var todo model.ToDo
		if err := rows.Scan(&todo.ID, &todo.Title, &todo.Status, &todo.ImagePath, &todo.ReminderTime); err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}
	return todos, nil
}

func (s *postgresStorage) GetTodoById(id string) (model.ToDo, error) {
	var todo model.ToDo
	err := s.conn.QueryRow(context.Background(),
		"SELECT id, title, status, image_path, reminder_time FROM todos WHERE id = $1", id).
		Scan(&todo.ID, &todo.Title, &todo.Status, &todo.ImagePath, &todo.ReminderTime)
	if err != nil {
		return model.ToDo{}, err
	}
	return todo, nil
}

func (s *postgresStorage) GetTodoImageById(id string) (model.ToDo, error) {
	var todo model.ToDo
	err := s.conn.QueryRow(context.Background(), "SELECT image_path FROM todos WHERE id = $1", id).Scan(&todo.ImagePath)
	if err != nil {
		return model.ToDo{}, err
	}
	return todo, nil
}

func (s *postgresStorage) UpdateTodo(id string, todo model.ToDo) error {
	_, err := s.conn.Exec(context.Background(), "UPDATE todos SET title = $1, status = $2 WHERE id = $3", todo.Title, todo.Status, id)
	return err
}

func (s *postgresStorage) UpdateTodoImage(id string, imagePath string) error {
	_, err := s.conn.Exec(context.Background(), "UPDATE todos SET image_path = $1 WHERE id = $2", imagePath, id)
	return err
}

func (s *postgresStorage) DeleteTodo(id string) error {
	_, err := s.conn.Exec(context.Background(), "DELETE FROM todos WHERE id = $1", id)
	return err
}

func (s *postgresStorage) Close() {
	s.conn.Close(context.Background())
}
