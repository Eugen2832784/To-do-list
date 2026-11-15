package storage

import (
	"path/filepath"
	"reflect"
	"runtime"
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

func TestSaveCreatesMissingDir(t *testing.T) {
	path := filepath.Join(t.TempDir(), "nested", "dir", "tasks.json")
	s := New(path)

	want := []task.Task{{ID: 1, Text: "задача"}}
	if err := s.Save(want); err != nil {
		t.Fatalf("Save() в несуществующую папку вернул ошибку: %v", err)
	}

	got, err := s.Load()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Load() = %+v, ожидалось %+v", got, want)
	}
}

func TestDefaultPathFromEnv(t *testing.T) {
	custom := filepath.Join(t.TempDir(), "custom.json")
	t.Setenv("TODO_FILE", custom)

	got, err := DefaultPath()
	if err != nil {
		t.Fatal(err)
	}
	if got != custom {
		t.Errorf("DefaultPath() = %q, ожидалось %q", got, custom)
	}
}

func TestDefaultPathInHome(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("тест использует переменную HOME")
	}
	home := t.TempDir()
	t.Setenv("TODO_FILE", "")
	t.Setenv("HOME", home)

	got, err := DefaultPath()
	if err != nil {
		t.Fatal(err)
	}
	want := filepath.Join(home, ".todo", "tasks.json")
	if got != want {
		t.Errorf("DefaultPath() = %q, ожидалось %q", got, want)
	}
}
