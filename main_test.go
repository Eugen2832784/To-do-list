package main

import (
	"path/filepath"
	"reflect"
	"testing"
)

// useTempDataFile направляет программу на временный файл,
// чтобы тесты не трогали настоящий tasks.json.
func useTempDataFile(t *testing.T) {
	t.Helper()
	old := dataFile
	dataFile = filepath.Join(t.TempDir(), "tasks.json")
	t.Cleanup(func() { dataFile = old })
}

func TestNextID(t *testing.T) {
	tests := []struct {
		name  string
		tasks []Task
		want  int
	}{
		{"пустой список", []Task{}, 1},
		{"одна задача", []Task{{ID: 1}}, 2},
		{"после удаления остаётся дырка", []Task{{ID: 2}, {ID: 5}}, 6},
		{"неупорядоченные id", []Task{{ID: 7}, {ID: 3}}, 8},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := nextID(tt.tasks); got != tt.want {
				t.Errorf("nextID() = %d, ожидалось %d", got, tt.want)
			}
		})
	}
}

func TestFindIndex(t *testing.T) {
	tasks := []Task{{ID: 2}, {ID: 5}, {ID: 9}}

	tests := []struct {
		name string
		id   int
		want int
	}{
		{"первая задача", 2, 0},
		{"последняя задача", 9, 2},
		{"не существует", 4, -1},
		{"нулевой id", 0, -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := findIndex(tasks, tt.id); got != tt.want {
				t.Errorf("findIndex(%d) = %d, ожидалось %d", tt.id, got, tt.want)
			}
		})
	}
}

func TestSaveAndLoad(t *testing.T) {
	useTempDataFile(t)

	want := []Task{
		{ID: 1, Text: "первая", Done: false},
		{ID: 2, Text: "вторая", Done: true},
	}

	if err := save(want); err != nil {
		t.Fatalf("save() вернул ошибку: %v", err)
	}

	got, err := load()
	if err != nil {
		t.Fatalf("load() вернул ошибку: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("load() = %+v, ожидалось %+v", got, want)
	}
}

func TestLoadMissingFile(t *testing.T) {
	useTempDataFile(t)

	got, err := load()
	if err != nil {
		t.Fatalf("load() без файла вернул ошибку: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("ожидался пустой список, получили %+v", got)
	}
}

func TestRunFlow(t *testing.T) {
	useTempDataFile(t)

	steps := [][]string{
		{"add", "купить", "молоко"},
		{"add", "погулять с собакой"},
	}
	for _, args := range steps {
		if err := run(args); err != nil {
			t.Fatalf("run(%v) вернул ошибку: %v", args, err)
		}
	}

	tasks, err := load()
	if err != nil {
		t.Fatal(err)
	}
	if len(tasks) != 2 {
		t.Fatalf("ожидалось 2 задачи, получили %d", len(tasks))
	}
	if tasks[0].Text != "купить молоко" {
		t.Errorf("текст = %q, ожидалось %q", tasks[0].Text, "купить молоко")
	}

	if err := run([]string{"done", "1"}); err != nil {
		t.Fatal(err)
	}
	tasks, _ = load()
	if !tasks[0].Done {
		t.Error("задача 1 должна быть выполненной")
	}

	if err := run([]string{"delete", "1"}); err != nil {
		t.Fatal(err)
	}
	tasks, _ = load()
	if len(tasks) != 1 || tasks[0].ID != 2 {
		t.Errorf("после удаления должна остаться задача с id 2, получили %+v", tasks)
	}
}

func TestRunErrors(t *testing.T) {
	useTempDataFile(t)

	tests := []struct {
		name string
		args []string
	}{
		{"add без текста", []string{"add"}},
		{"done без id", []string{"done"}},
		{"done с неверным id", []string{"done", "abc"}},
		{"done несуществующей задачи", []string{"done", "99"}},
		{"delete несуществующей задачи", []string{"delete", "99"}},
		{"неизвестная команда", []string{"foo"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := run(tt.args); err == nil {
				t.Errorf("run(%v) должен вернуть ошибку", tt.args)
			}
		})
	}
}
