package task

import "testing"

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
			if got := NextID(tt.tasks); got != tt.want {
				t.Errorf("NextID() = %d, ожидалось %d", got, tt.want)
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
			if got := FindIndex(tasks, tt.id); got != tt.want {
				t.Errorf("FindIndex(%d) = %d, ожидалось %d", tt.id, got, tt.want)
			}
		})
	}
}
