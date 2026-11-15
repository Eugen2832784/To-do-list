package storage

import (
	"encoding/json"
	"errors"
	"os"

	"todo/internal/task"
)

type Storage struct {
	path string
}

func New(path string) *Storage {
	return &Storage{path: path}
}

// Load читает задачи из файла. Если файла ещё нет, возвращает пустой список.
func (s *Storage) Load() ([]task.Task, error) {
	data, err := os.ReadFile(s.path)
	if errors.Is(err, os.ErrNotExist) {
		return []task.Task{}, nil
	}
	if err != nil {
		return nil, err
	}

	var tasks []task.Task
	if err := json.Unmarshal(data, &tasks); err != nil {
		return nil, err
	}
	return tasks, nil
}

// Save записывает задачи в файл.
func (s *Storage) Save(tasks []task.Task) error {
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0o644)
}
