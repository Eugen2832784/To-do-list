package storage

import (
	"path/filepath"
	"reflect"
	"testing"

	"todo/internal/task"
)

func newTestStorage(t *testing.T) *Storage {
	t.Helper()
	return New(filepath.Join(t.TempDir(), "tasks.json"))
}

func TestSaveAndLoad(t *testing.T) {
	s := newTestStorage(t)

	want := []task.Task{
		{ID: 1, Text: "первая", Done: false},
		{ID: 2, Text: "вторая", Done: true},
	}

	if err := s.Save(want); err != nil {
		t.Fatalf("Save() вернул ошибку: %v", err)
	}

	got, err := s.Load()
	if err != nil {
		t.Fatalf("Load() вернул ошибку: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Load() = %+v, ожидалось %+v", got, want)
	}
}

func TestLoadMissingFile(t *testing.T) {
	s := newTestStorage(t)

	got, err := s.Load()
	if err != nil {
		t.Fatalf("Load() без файла вернул ошибку: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("ожидался пустой список, получили %+v", got)
	}
}
