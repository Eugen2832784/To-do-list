package cli

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"todo/internal/storage"
	"todo/internal/task"
)

type App struct {
	store *storage.Storage
	out   io.Writer
}


func New(store *storage.Storage, out io.Writer) *App {
	return &App{store: store, out: out}
}

func (a *App) usage() {
	fmt.Fprintln(a.out, `Использование:
  todo add "текст задачи"
  todo list
  todo done <id>
  todo delete <id>`)
}

func (a *App) Run(args []string) error {
	if len(args) < 1 {
		a.usage()
		return nil
	}

	tasks, err := a.store.Load()
	if err != nil {
		return err
	}

	switch args[0] {
	case "add":
		if len(args) < 2 {
			return errors.New("укажи текст задачи")
		}
		text := strings.Join(args[1:], " ")
		tasks = append(tasks, task.Task{ID: task.NextID(tasks), Text: text})
		fmt.Fprintln(a.out, "Добавлено:", text)

	case "list":
		if len(tasks) == 0 {
			fmt.Fprintln(a.out, "Задач нет")
			return nil
		}
		for _, t := range tasks {
			mark := " "
			if t.Done {
				mark = "x"
			}
			fmt.Fprintf(a.out, "%d. [%s] %s\n", t.ID, mark, t.Text)
		}
		return nil

	case "done", "delete":
		if len(args) < 2 {
			return errors.New("укажи id задачи")
		}
		id, err := strconv.Atoi(args[1])
		if err != nil {
			return fmt.Errorf("некорректный id %q", args[1])
		}
		i := task.FindIndex(tasks, id)
		if i == -1 {
			return fmt.Errorf("задача %d не найдена", id)
		}
		if args[0] == "done" {
			tasks[i].Done = true
			fmt.Fprintln(a.out, "Выполнено:", tasks[i].Text)
		} else {
			fmt.Fprintln(a.out, "Удалено:", tasks[i].Text)
			tasks = append(tasks[:i], tasks[i+1:]...)
		}

	default:
		a.usage()
		return fmt.Errorf("неизвестная команда %q", args[0])
	}

	return a.store.Save(tasks)
}
