// Package task описывает задачу и базовые операции над списком задач.
package task

type Task struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
	Done bool   `json:"done"`
}

// NextID возвращает свободный id: максимальный существующий + 1.
func NextID(tasks []Task) int {
	maxID := 0
	for _, t := range tasks {
		if t.ID > maxID {
			maxID = t.ID
		}
	}
	return maxID + 1
}

// FindIndex возвращает индекс задачи с указанным id или -1, если такой нет.
func FindIndex(tasks []Task, id int) int {
	for i, t := range tasks {
		if t.ID == id {
			return i
		}
	}
	return -1
}
