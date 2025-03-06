package task

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/nktauserum/web-calculation/shared"
)

func TestQueue_AddTask(t *testing.T) {
	q := &Queue{
		Tasks:       make(map[int64]shared.Task),
		Expressions: make(map[int64]shared.Expression),
	}

	task := shared.Task{
		FirstArgument:  "5",
		SecondArgument: "3",
		Operator:       '+',
	}

	id := q.AddTask(task)
	assert.Equal(t, int64(1), id)
	assert.Equal(t, task, q.Tasks[id])
}

func TestQueue_ParseExpression(t *testing.T) {
	tests := []struct {
		name        string
		expression  string
		wantErr     bool
		taskCount   int
		resultValue string
	}{
		{
			name:        "Simple addition",
			expression:  "2 + 3",
			wantErr:     false,
			taskCount:   1,
			resultValue: "id1",
		},
		{
			name:        "Complex expression",
			expression:  "2 + 3 * 4",
			wantErr:     false,
			taskCount:   2,
			resultValue: "id2",
		},
		{
			name:        "Invalid expression",
			expression:  "2 + + 3",
			wantErr:     true,
			taskCount:   0,
			resultValue: "",
		},
		{
			name:        "Division by zero",
			expression:  "5 / 0",
			wantErr:     true,
			taskCount:   0,
			resultValue: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := &Queue{
				Tasks:       make(map[int64]shared.Task),
				Expressions: make(map[int64]shared.Expression),
			}

			id, err := q.ParseExpression(tt.expression)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, len(q.Tasks), tt.taskCount)
			assert.Equal(t, tt.resultValue, q.Expressions[id].Result)
		})
	}
}

func TestQueue_Done(t *testing.T) {
	q := &Queue{
		Tasks:       make(map[int64]shared.Task),
		Expressions: make(map[int64]shared.Expression),
	}

	task := shared.Task{
		ID:             1,
		FirstArgument:  "5",
		SecondArgument: "3",
		Operator:       '+',
		Status:         false,
	}

	q.Tasks[1] = task
	q.Done(1, 8.0)

	updatedTask := q.Tasks[1]
	assert.True(t, updatedTask.Status)
	assert.Equal(t, 8.0, updatedTask.Result)
}

func TestQueue_UpdateTasks(t *testing.T) {
	q := &Queue{
		Tasks:       make(map[int64]shared.Task),
		Expressions: make(map[int64]shared.Expression),
	}

	// Add first task
	task1 := shared.Task{
		ID:             1,
		FirstArgument:  "5",
		SecondArgument: "3",
		Operator:       '+',
		Status:         true,
		Result:         8.0,
	}
	q.Tasks[1] = task1

	// Add second task that depends on first task
	task2 := shared.Task{
		ID:             2,
		FirstArgument:  "id1",
		SecondArgument: "2",
		Operator:       '*',
		Status:         false,
	}
	q.Tasks[2] = task2

	err := q.UpdateTasks()
	assert.NoError(t, err)
	assert.Equal(t, "8", q.Tasks[2].FirstArgument)
}

func TestConvertToRPN(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected []string
		wantErr  bool
	}{
		{
			name:     "Simple addition",
			input:    []string{"2", "+", "3"},
			expected: []string{"2", "3", "+"},
			wantErr:  false,
		},
		{
			name:     "Expression with parentheses",
			input:    []string{"(", "2", "+", "3", ")", "*", "4"},
			expected: []string{"2", "3", "+", "4", "*"},
			wantErr:  false,
		},
		{
			name:     "Unmatched parentheses",
			input:    []string{"(", "2", "+", "3"},
			expected: nil,
			wantErr:  true,
		},
		{
			name:     "Division by zero",
			input:    []string{"5", "/", "0"},
			expected: nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := convertToRPN(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}
