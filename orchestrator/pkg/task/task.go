package task

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/veliashev/web-calculation/shared"
	"github.com/veliashev/web-calculation/shared/errors"
)

// Операция в выражении
type Operation rune

const (
	Add      Operation = '+'
	Subtract Operation = '-'
	Multiply Operation = '*'
	Divide   Operation = '/'
)

type ExpressionQueue struct {
	Expressions []shared.Expression
	NextID      int64
}

// Очередь
type TasksQueue struct {
	Tasks  []shared.Task // Список задач
	NextID int64         // Следующий доступный ID
}

type Queue struct {
	Tasks       map[int64]shared.Task
	Expressions map[int64]shared.Expression
}

// Добавляет задачу в очередь. Возвращает ID переданной задачи в очереди

func (q *Queue) AddTask(task shared.Task) int64 {
	// Находим максимальный ID существующих задач
	maxID := int64(0)
	for id := range q.Tasks {
		if id > maxID {
			maxID = id
		}
	}

	// Присваиваем новому task ID
	newID := maxID + 1
	q.Tasks[newID] = task // Добавляем задачу в очередь

	return newID // Возвращаем новый ID
}

func (q *Queue) AddExpression(expression shared.Expression) int64 {
	// Находим максимальный ID существующих выражений
	maxID := int64(0)
	for id := range q.Expressions {
		if id > maxID {
			maxID = id
		}
	}

	// Присваиваем новому expression ID
	newID := maxID + 1
	q.Expressions[newID] = expression // Добавляем выражение в очередь

	return newID // Возвращаем новый ID
}

// Помечает задачу как выполненную
func (q *Queue) Done(id int64, result float64) {
	task, exists := q.Tasks[id]
	if !exists {
		return
	}

	task.Status = true
	task.Result = result
	q.Tasks[id] = task

	if err := q.UpdateTasks(); err != nil {
		log.Printf("error updating task %d: %v", id, err)
		return
	}
	if err := q.UpdateExpressions(); err != nil {
		log.Printf("error updating expressions on task %d: %v", id, err)
		return
	}
}

// Получаем абсолютно все задачи из очереди
func (q *Queue) GetTasks() map[int64]shared.Task {
	return q.Tasks
}

func (q *Queue) GetExpressions() map[int64]shared.Expression {
	err := q.UpdateExpressions()
	if err != nil {
		log.Printf("error updating expressions: %v", err)
	}
	return q.Expressions
}

func (q *Queue) UpdateExpressions() error {
	expressions := q.Expressions

	for _, exp := range expressions {
		if !IsNumeric(exp.Result) {
			id, err := strconv.ParseInt(strings.TrimPrefix(exp.Result, "id"), 10, 64)
			if err != nil {
				log.Printf("Ошибка в парсинге %s", strings.TrimPrefix(exp.Result, "id"))
				return err
			}

			relatedTask, exists := q.Tasks[id]
			if !exists {
				log.Printf("Задача ID: %d равна nil\n", id)
				return err
			}

			if !relatedTask.Status {
				log.Printf("Задача ID: %d ещё не выполнена\n", id)
				return err
			}

			currentExpression := q.Expressions[exp.ID]

			log.Printf("Результат с id %d", exp.ID)
			currentExpression.Result = strconv.FormatFloat(relatedTask.Result, 'f', -1, 64)
			log.Printf(" преобразован в значение (%s)\n", currentExpression.Result)
			currentExpression.Status = true
			log.Printf("Выражение %d выполнено!\n", currentExpression.ID)
			q.Expressions[exp.ID] = currentExpression
		}
	}

	return nil
}

func (q *Queue) FindExpression(id int64) *shared.Expression {
	expr, exists := q.Expressions[id]
	if exists {
		return &expr
	}
	return nil
}

func (q *Queue) FindTask(id int64) *shared.Task {
	task, exists := q.Tasks[id]
	if exists {
		return &task
	}
	return nil
}

func (q *Queue) UpdateTasks() error {
	for _, task := range q.Tasks {
		if !IsNumeric(task.FirstArgument) {
			log.Printf("Аргумент 1 не является числом: %s\n", task.FirstArgument)
			id, err := strconv.ParseInt(strings.TrimPrefix(task.FirstArgument, "id"), 10, 64)
			if err != nil {
				log.Printf("Ошибка в парсинге %s", strings.TrimPrefix(task.FirstArgument, "id"))
				return err
			}

			relatedTask := q.FindTask(id)
			if relatedTask == nil {
				log.Printf("Задача ID %d равна nil\n", id)
				continue
			}

			if !relatedTask.Status {
				log.Printf("Задача ID %d ещё не выполнена\n", id)
				continue
			}

			log.Printf("Аргумент 1 (%s)", task.FirstArgument)
			task.FirstArgument = strconv.FormatFloat(relatedTask.Result, 'f', -1, 64)
			log.Printf(" преобразован в значение (%s)\n", task.FirstArgument)
			q.Tasks[task.ID] = task
		}

		if !IsNumeric(task.SecondArgument) {
			log.Printf("Аргумент 2 не является числом: %s\n", task.SecondArgument)
			id, err := strconv.ParseInt(strings.TrimPrefix(task.SecondArgument, "id"), 10, 64)
			if err != nil {
				log.Printf("Ошибка в парсинге %s", strings.TrimPrefix(task.FirstArgument, "id"))
				return err
			}

			relatedTask := q.FindTask(id)
			if relatedTask == nil {
				log.Printf("Задача ID %d равна nil\n", id)
				continue
			}

			if !relatedTask.Status {
				log.Printf("Задача ID %d ещё не выполнена\n", id)
				continue
			}

			log.Printf("Аргумент 2 (%s)", task.SecondArgument)
			task.SecondArgument = strconv.FormatFloat(relatedTask.Result, 'f', -1, 64)
			log.Printf(" преобразован в значение (%s)\n", task.SecondArgument)
			q.Tasks[task.ID] = task
		}
	}
	return nil
}

func convertToRPN(tokens []string) ([]string, error) {
	var stack []string
	var output []string
	var precedence = map[rune]int{
		'+': 1,
		'-': 1,
		'*': 2,
		'/': 2,
	}

	parenthesesCount := 0

	for i, token := range tokens {
		switch token {
		case "(":
			parenthesesCount++
			stack = append(stack, token)
		case ")":
			parenthesesCount--
			if parenthesesCount < 0 {
				return nil, errors.ErrMismatchedParentheses
			}
			for len(stack) > 0 && stack[len(stack)-1] != "(" {
				output = append(output, stack[len(stack)-1])
				stack = stack[:len(stack)-1]
			}
			if len(stack) > 0 {
				stack = stack[:len(stack)-1] // Remove "("
			}
		case "+", "-", "*", "/":
			if token == "/" && i < len(tokens)-1 && tokens[i+1] == "0" {
				return nil, errors.ErrDivisionByZero
			}
			// Проверяем, является ли минус унарным
			if token == "-" && (i == 0 || tokens[i-1] == "(" || isOperator(tokens[i-1])) {
				// Добавляем 0 перед унарным минусом
				output = append(output, "0")
			}
			for len(stack) > 0 {
				top := stack[len(stack)-1]
				if top != "(" && precedence[rune(top[0])] >= precedence[rune(token[0])] {
					output = append(output, stack[len(stack)-1])
					stack = stack[:len(stack)-1]
				} else {
					break
				}
			}
			stack = append(stack, token)
		default:
			if !IsNumeric(token) && !strings.HasPrefix(token, "id") {
				return nil, errors.ErrInvalidNumber
			}
			output = append(output, token)
		}
	}

	if parenthesesCount != 0 {
		return nil, errors.ErrMismatchedParentheses
	}

	for len(stack) > 0 {
		output = append(output, stack[len(stack)-1])
		stack = stack[:len(stack)-1]
	}

	return output, nil
}

func (q *Queue) generateTasksFromRPN(output []string) ([]shared.Task, map[int]string, error) {
	var tasks []shared.Task
	var operandStack []string
	taskIDs := make(map[int]string)
	nextID := int64(1)

	for _, token := range output {
		if isOperator(token) {
			if len(operandStack) < 2 {
				return nil, nil, errors.ErrNotEnoughOperands
			}
			arg2 := operandStack[len(operandStack)-1]
			arg1 := operandStack[len(operandStack)-2]
			operandStack = operandStack[:len(operandStack)-2]

			if token == "/" && arg2 == "0" {
				return nil, nil, errors.ErrDivisionByZero
			}

			task := shared.Task{
				ID:             nextID,
				FirstArgument:  arg1,
				SecondArgument: arg2,
				Operator:       rune(token[0]),
				Status:         false,
			}
			tasks = append(tasks, task)
			taskIDs[int(nextID)] = fmt.Sprintf("id%d", nextID)
			operandStack = append(operandStack, fmt.Sprintf("id%d", nextID))
			nextID++
		} else {
			operandStack = append(operandStack, token)
		}
	}

	return tasks, taskIDs, nil
}

// Проверяет, является ли токен оператором
func isOperator(token string) bool {
	switch token {
	case "+", "-", "*", "/":
		return true
	default:
		return false
	}
}

func (q *Queue) ParseExpression(expression string) (int64, error) {
	tokens := tokenize(expression)
	output, err := convertToRPN(tokens)
	if err != nil {
		return 0, err
	}

	tasks, _, err := q.generateTasksFromRPN(output)
	if err != nil {
		return 0, err
	}

	if len(tasks) == 0 {
		return 0, errors.ErrInvalidExpression
	}

	// Добавляем задачи в очередь
	for _, task := range tasks {
		q.Tasks[task.ID] = task
	}

	// Добавляем выражение в массив выражений
	nextID := int64(len(q.Expressions) + 1)
	expr := shared.Expression{
		ID:     nextID,
		Status: false,
		Result: fmt.Sprintf("id%d", tasks[len(tasks)-1].ID),
	}
	q.Expressions[nextID] = expr

	return expr.ID, nil
}
