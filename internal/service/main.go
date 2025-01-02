package service

import (
	"errors"
	"toDoList/internal/model"
	"toDoList/internal/storage"
)

type TodoService interface {
	GetAllTodos() ([]model.ToDo, error)
	GetTodoById(id string) (model.ToDo, error)
	GetTodoImageById(id string) (model.ToDo, error)
	AddTodo(todo model.ToDo) error
	UpdateTodo(id string, todo model.ToDo) error
	UpdateTodoImage(id string, imagePath string) error
	DeleteTodo(id string) error
}

type todoService struct {
	storage storage.Storage
}

func NewTodoService(storage storage.Storage) TodoService {
	return &todoService{storage: storage}
}

func (s *todoService) GetAllTodos() ([]model.ToDo, error) {
	return s.storage.GetTodos()
}

func (s *todoService) GetTodoById(id string) (model.ToDo, error) {
	return s.storage.GetTodoById(id)
}

func (s *todoService) GetTodoImageById(id string) (model.ToDo, error) {
	return s.storage.GetTodoImageById(id)
}

func (s *todoService) AddTodo(todo model.ToDo) error {
	if !model.IsValidStatus(todo.Status) {
		return errors.New("invalid status")
	}
	return s.storage.AddTodo(todo)
}

func (s *todoService) UpdateTodo(id string, todo model.ToDo) error {
	if !model.IsValidStatus(todo.Status) {
		return errors.New("invalid status")
	}
	return s.storage.UpdateTodo(id, todo)
}

func (s *todoService) UpdateTodoImage(id string, imagePath string) error {
	return s.storage.UpdateTodoImage(id, imagePath)
}

func (s *todoService) DeleteTodo(id string) error {
	return s.storage.DeleteTodo(id)
}
