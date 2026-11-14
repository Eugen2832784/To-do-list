package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var dataFile = "tasks.json"

type Task struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
	Done bool   `json:"done"`
}

func load() ([]Task, error) {
	data, err := os.ReadFile(dataFile)
	if errors.Is(err, os.ErrNotExist) {
		return []Task{}, nil
	}
	if err != nil {
		return nil, err
	}
	var tasks []Task
	if err := json.Unmarshal(data, &tasks); err != nil {
		return nil, err
	}
	return tasks, nil
}

func save(tasks []Task) error {
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(dataFile, data, 0o644)
}

func nextID(tasks []Task) int {
	maxID := 0
	for _, t := range tasks {
		if t.ID > max {
			max = t.ID
		}
	}
	return max + 1
}

func findIndex(tasks []Task, id int) int {
	for i, t := range tasks {
		if t.ID == id {
			return i
		}
	}
	return -1
}

func usage() {
	fmt.Println(`Использование:
  todo add "текст задачи"
  todo list
  todo done <id>
  todo delete <id>`)
}

func run(args []string) error {
	if len(args) < 1 {
		usage()
		return nil
	}

	tasks, err := load()
	if err != nil {
		return err
	}

	switch args[0] {
	case "add":
		if len(args) < 2 {
			return errors.New("укажи текст задачи")
		}
		text := strings.Join(args[1:], " ")
		tasks = append(tasks, Task{ID: nextID(tasks), Text: text})
		fmt.Println("Добавлено:", text)

	case "list":
		if len(tasks) == 0 {
			fmt.Println("Задач нет")
			return nil
		}
		for _, t := range tasks {
			mark := " "
			if t.Done {
				mark = "x"
			}
			fmt.Printf("%d. [%s] %s\n", t.ID, mark, t.Text)
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
		i := findIndex(tasks, id)
		if i == -1 {
			return fmt.Errorf("задача %d не найдена", id)
		}
		if args[0] == "done" {
			tasks[i].Done = true
			fmt.Println("Выполнено:", tasks[i].Text)
		} else {
			fmt.Println("Удалено:", tasks[i].Text)
			tasks = append(tasks[:i], tasks[i+1:]...)
		}

	default:
		usage()
		return fmt.Errorf("неизвестная команда %q", args[0])
	}

	return save(tasks)
}

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "Ошибка:", err)
		os.Exit(1)
	}
}
