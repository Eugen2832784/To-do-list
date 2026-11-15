package cli

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"

	"todo/internal/storage"
)

func newTestApp(t *testing.T) (*App, *storage.Storage, *bytes.Buffer) {
	t.Helper()
	store := storage.New(filepath.Join(t.TempDir(), "tasks.json"))
	out := &bytes.Buffer{}
	return New(store, out), store, out
}

func TestRunFlow(t *testing.T) {
	app, store, out := newTestApp(t)

	steps := [][]string{
		{"add", "купить", "молоко"},
		{"add", "погулять с собакой"},
	}
	for _, args := range steps {
		if err := app.Run(args); err != nil {
			t.Fatalf("Run(%v) вернул ошибку: %v", args, err)
		}
	}

	tasks, err := store.Load()
	if err != nil {
		t.Fatal(err)
	}
	if len(tasks) != 2 {
		t.Fatalf("ожидалось 2 задачи, получили %d", len(tasks))
	}
	if tasks[0].Text != "купить молоко" {
		t.Errorf("текст = %q, ожидалось %q", tasks[0].Text, "купить молоко")
	}

	out.Reset()
	if err := app.Run([]string{"list"}); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), "1. [ ] купить молоко") {
		t.Errorf("в выводе list нет первой задачи:\n%s", out.String())
	}

	if err := app.Run([]string{"done", "1"}); err != nil {
		t.Fatal(err)
	}
	out.Reset()
	if err := app.Run([]string{"list"}); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), "1. [x] купить молоко") {
		t.Errorf("задача 1 должна быть выполненной:\n%s", out.String())
	}

	if err := app.Run([]string{"delete", "1"}); err != nil {
		t.Fatal(err)
	}
	tasks, _ = store.Load()
	if len(tasks) != 1 || tasks[0].ID != 2 {
		t.Errorf("после удаления должна остаться задача с id 2, получили %+v", tasks)
	}
}

func TestRunEmptyList(t *testing.T) {
	app, _, out := newTestApp(t)

	if err := app.Run([]string{"list"}); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), "Задач нет") {
		t.Errorf("ожидалось сообщение о пустом списке, получили %q", out.String())
	}
}

func TestRunErrors(t *testing.T) {
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
			app, _, _ := newTestApp(t)
			if err := app.Run(tt.args); err == nil {
				t.Errorf("Run(%v) должен вернуть ошибку", tt.args)
			}
		})
	}
}

func TestListFilters(t *testing.T) {
	app, _, out := newTestApp(t)

	for _, args := range [][]string{
		{"add", "первая"},
		{"add", "вторая"},
		{"done", "1"},
	} {
		if err := app.Run(args); err != nil {
			t.Fatalf("Run(%v) вернул ошибку: %v", args, err)
		}
	}

	tests := []struct {
		name        string
		args        []string
		wantContain string
		wantAbsent  string
	}{
		{"без фильтра", []string{"list"}, "вторая", ""},
		{"только невыполненные", []string{"list", "--todo"}, "вторая", "первая"},
		{"только выполненные", []string{"list", "--done"}, "первая", "вторая"},
		{"один дефис тоже работает", []string{"list", "-todo"}, "вторая", "первая"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out.Reset()
			if err := app.Run(tt.args); err != nil {
				t.Fatal(err)
			}
			got := out.String()
			if !strings.Contains(got, tt.wantContain) {
				t.Errorf("в выводе нет %q:\n%s", tt.wantContain, got)
			}
			if tt.wantAbsent != "" && strings.Contains(got, tt.wantAbsent) {
				t.Errorf("в выводе не должно быть %q:\n%s", tt.wantAbsent, got)
			}
		})
	}
}

func TestListFilterNoMatches(t *testing.T) {
	app, _, out := newTestApp(t)

	if err := app.Run([]string{"add", "единственная"}); err != nil {
		t.Fatal(err)
	}
	out.Reset()
	if err := app.Run([]string{"list", "--done"}); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), "Подходящих задач нет") {
		t.Errorf("ожидалось сообщение об отсутствии задач, получили %q", out.String())
	}
}

func TestListFilterErrors(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{"оба флага сразу", []string{"list", "--todo", "--done"}},
		{"неизвестный флаг", []string{"list", "--foo"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app, _, _ := newTestApp(t)
			if err := app.Run(tt.args); err == nil {
				t.Errorf("Run(%v) должен вернуть ошибку", tt.args)
			}
		})
	}
}

func TestRunClear(t *testing.T) {
	app, store, out := newTestApp(t)

	for _, args := range [][]string{{"add", "первая"}, {"add", "вторая"}} {
		if err := app.Run(args); err != nil {
			t.Fatal(err)
		}
	}

	out.Reset()
	if err := app.Run([]string{"clear"}); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), "Все задачи удалены") {
		t.Errorf("нет сообщения об удалении, получили %q", out.String())
	}

	tasks, err := store.Load()
	if err != nil {
		t.Fatal(err)
	}
	if len(tasks) != 0 {
		t.Fatalf("после clear ожидался пустой список, получили %+v", tasks)
	}

	// После очистки нумерация начинается заново.
	if err := app.Run([]string{"add", "новая"}); err != nil {
		t.Fatal(err)
	}
	tasks, _ = store.Load()
	if len(tasks) != 1 || tasks[0].ID != 1 {
		t.Errorf("новая задача должна получить id 1, получили %+v", tasks)
	}
}
