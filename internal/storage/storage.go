// Package storage сохраняет и загружает задачи из JSON-файла.
package storage

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

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
	// Создаём папку для файла, если её ещё нет (например, ~/.todo).
	if err := os.MkdirAll(filepath.Dir(s.path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0o644)
}

// DefaultPath возвращает путь к файлу с задачами.
// Если задана переменная окружения TODO_FILE, используется она,
// иначе ~/.todo/tasks.json.
func DefaultPath() (string, error) {
	if p := os.Getenv("TODO_FILE"); p != "" {
		return p, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".todo", "tasks.json"), nil
}
