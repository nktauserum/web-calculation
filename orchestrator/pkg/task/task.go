package task

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"github.com/nktauserum/web-calculation/shared"
	"github.com/nktauserum/web-calculation/shared/errors"
)

// Операция в выражении
type Operation rune

const (
	Add      Operation = '+'
	Subtract Operation = '-'
	Multiply Operation = '*'
	Divide   Operation = '/'
)

type Queue struct {
	db *sql.DB
}

// NewQueue создает новую очередь с SQLite
func NewQueue(dbPath string) (*Queue, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("ошибка открытия базы данных: %w", err)
	}

	db.Exec("PRAGMA journal_mode = WAL;")

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ошибка подключения к базе данных: %w", err)
	}

	q := &Queue{db: db}
	if err := q.initTables(); err != nil {
		return nil, fmt.Errorf("ошибка инициализации таблиц: %w", err)
	}

	return q, nil
}

// Close закрывает соединение с базой данных
func (q *Queue) Close() error {
	return q.db.Close()
}

// initTables создает необходимые таблицы, если они не существуют
func (q *Queue) initTables() error {
	// Создаем таблицу для задач
	_, err := q.db.Exec(`
		CREATE TABLE IF NOT EXISTS tasks (
			id INTEGER PRIMARY KEY,
			first_argument TEXT NOT NULL,
			second_argument TEXT NOT NULL,
			operator TEXT NOT NULL,
			status BOOLEAN NOT NULL DEFAULT 0,
			result REAL
		)
	`)
	if err != nil {
		return err
	}

	// Создаем таблицу для выражений
	_, err = q.db.Exec(`
		CREATE TABLE IF NOT EXISTS expressions (
			id INTEGER PRIMARY KEY,
			status BOOLEAN NOT NULL DEFAULT 0,
			result TEXT NOT NULL
		)
	`)
	return err
}

// AddTask добавляет задачу в очередь. Возвращает ID переданной задачи в очереди
func (q *Queue) AddTask(task shared.Task) int64 {
	// Находим максимальный ID существующих задач
	var maxID int64
	err := q.db.QueryRow("SELECT COALESCE(MAX(id), 0) FROM tasks").Scan(&maxID)
	if err != nil {
		log.Printf("Ошибка при получении максимального ID задачи: %v", err)
		return 0
	}

	// Присваиваем новый ID
	newID := maxID + 1
	task.ID = newID

	// Добавляем задачу в базу данных
	_, err = q.db.Exec(
		"INSERT INTO tasks (id, first_argument, second_argument, operator, status, result) VALUES (?, ?, ?, ?, ?, ?)",
		task.ID, task.FirstArgument, task.SecondArgument, string(task.Operator), task.Status, task.Result,
	)
	if err != nil {
		log.Printf("Ошибка при добавлении задачи: %v", err)
		return 0
	}

	return newID
}

func (q *Queue) AddExpression(expression shared.Expression) int64 {
	// Находим максимальный ID существующих выражений
	var maxID int64
	err := q.db.QueryRow("SELECT COALESCE(MAX(id), 0) FROM expressions").Scan(&maxID)
	if err != nil {
		log.Printf("Ошибка при получении максимального ID выражения: %v", err)
		return 0
	}

	// Присваиваем новый ID
	newID := maxID + 1
	expression.ID = newID

	// Добавляем выражение в базу данных
	_, err = q.db.Exec(
		"INSERT INTO expressions (id, status, result) VALUES (?, ?, ?)",
		expression.ID, expression.Status, expression.Result,
	)
	if err != nil {
		log.Printf("Ошибка при добавлении выражения: %v", err)
		return 0
	}

	return newID
}

// Done помечает задачу как выполненную
func (q *Queue) Done(id int64, result float64) {
	// Получаем задачу из базы данных
	var task shared.Task
	var operatorStr string
	err := q.db.QueryRow(
		"SELECT id, first_argument, second_argument, operator, status FROM tasks WHERE id = ?",
		id,
	).Scan(&task.ID, &task.FirstArgument, &task.SecondArgument, &operatorStr, &task.Status)

	if err != nil {
		log.Printf("Ошибка при получении задачи %d: %v", id, err)
		return
	}

	if len(operatorStr) > 0 {
		task.Operator = rune(operatorStr[0])
	}

	// Обновляем статус и результат задачи
	task.Status = true
	task.Result = result
	_, err = q.db.Exec(
		"UPDATE tasks SET status = ?, result = ? WHERE id = ?",
		task.Status, task.Result, task.ID,
	)
	if err != nil {
		log.Printf("Ошибка при обновлении задачи %d: %v", id, err)
		return
	}

	if err := q.UpdateTasks(); err != nil {
		log.Printf("Ошибка при обновлении задач после задачи %d: %v", id, err)
		return
	}
	if err := q.UpdateExpressions(); err != nil {
		log.Printf("Ошибка при обновлении выражений после задачи %d: %v", id, err)
		return
	}
}

// GetTasks получает абсолютно все задачи из очереди
func (q *Queue) GetTasks() map[int64]shared.Task {
	tasks := make(map[int64]shared.Task)

	rows, err := q.db.Query("SELECT id, first_argument, second_argument, operator, status, result FROM tasks")
	if err != nil {
		log.Printf("Ошибка при получении задач: %v", err)
		return tasks
	}
	defer rows.Close()

	for rows.Next() {
		var task shared.Task
		var operatorStr string
		err := rows.Scan(&task.ID, &task.FirstArgument, &task.SecondArgument, &operatorStr, &task.Status, &task.Result)
		if err != nil {
			log.Printf("Ошибка при сканировании задачи: %v", err)
			continue
		}

		if len(operatorStr) > 0 {
			task.Operator = rune(operatorStr[0])
		}

		tasks[task.ID] = task
	}

	return tasks
}

func (q *Queue) GetExpressions() map[int64]shared.Expression {
	err := q.UpdateExpressions()
	if err != nil {
		log.Printf("Ошибка при обновлении выражений: %v", err)
	}

	expressions := make(map[int64]shared.Expression)

	rows, err := q.db.Query("SELECT id, status, result FROM expressions")
	if err != nil {
		log.Printf("Ошибка при получении выражений: %v", err)
		return expressions
	}
	defer rows.Close()

	for rows.Next() {
		var expr shared.Expression
		err := rows.Scan(&expr.ID, &expr.Status, &expr.Result)
		if err != nil {
			log.Printf("Ошибка при сканировании выражения: %v", err)
			continue
		}

		expressions[expr.ID] = expr
	}

	return expressions
}

func (q *Queue) UpdateExpressions() error {
	rows, err := q.db.Query("SELECT id, status, result FROM expressions")
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var expr shared.Expression
		err := rows.Scan(&expr.ID, &expr.Status, &expr.Result)
		if err != nil {
			log.Printf("Ошибка при сканировании выражения: %v", err)
			continue
		}

		if !IsNumeric(expr.Result) {
			id, err := strconv.ParseInt(strings.TrimPrefix(expr.Result, "id"), 10, 64)
			if err != nil {
				log.Printf("Ошибка в парсинге %s", strings.TrimPrefix(expr.Result, "id"))
				continue
			}

			var relatedTask shared.Task
			var operatorStr string
			err = q.db.QueryRow(
				"SELECT id, first_argument, second_argument, operator, status, result FROM tasks WHERE id = ?",
				id,
			).Scan(&relatedTask.ID, &relatedTask.FirstArgument, &relatedTask.SecondArgument, &operatorStr, &relatedTask.Status, &relatedTask.Result)

			if err != nil {
				log.Printf("Задача ID: %d не найдена: %v\n", id, err)
				continue
			}

			if len(operatorStr) > 0 {
				relatedTask.Operator = rune(operatorStr[0])
			}

			if !relatedTask.Status {
				log.Printf("Задача ID: %d ещё не выполнена\n", id)
				continue
			}

			log.Printf("Результат с id %d", expr.ID)
			expr.Result = strconv.FormatFloat(relatedTask.Result, 'f', -1, 64)
			log.Printf(" преобразован в значение (%s)\n", expr.Result)
			expr.Status = true
			log.Printf("Выражение %d выполнено!\n", expr.ID)

			_, err = q.db.Exec(
				"UPDATE expressions SET status = ?, result = ? WHERE id = ?",
				expr.Status, expr.Result, expr.ID,
			)
			if err != nil {
				log.Printf("Ошибка при обновлении выражения %d: %v", expr.ID, err)
				continue
			}
		}
	}

	return nil
}

func (q *Queue) FindExpression(id int64) *shared.Expression {
	var expr shared.Expression
	err := q.db.QueryRow(
		"SELECT id, status, result FROM expressions WHERE id = ?",
		id,
	).Scan(&expr.ID, &expr.Status, &expr.Result)

	if err != nil {
		if err != sql.ErrNoRows {
			log.Printf("Ошибка при поиске выражения %d: %v", id, err)
		}
		return nil
	}

	return &expr
}

func (q *Queue) FindTask(id int64) *shared.Task {
	var task shared.Task
	var operatorStr string
	err := q.db.QueryRow(
		"SELECT id, first_argument, second_argument, operator, status, result FROM tasks WHERE id = ?",
		id,
	).Scan(&task.ID, &task.FirstArgument, &task.SecondArgument, &operatorStr, &task.Status, &task.Result)

	if err != nil {
		if err != sql.ErrNoRows {
			log.Printf("Ошибка при поиске задачи %d: %v", id, err)
		}
		return nil
	}

	if len(operatorStr) > 0 {
		task.Operator = rune(operatorStr[0])
	}

	return &task
}

func (q *Queue) UpdateTasks() error {
	rows, err := q.db.Query("SELECT id, first_argument, second_argument, operator, status, result FROM tasks")
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var task shared.Task
		var operatorStr string
		err := rows.Scan(&task.ID, &task.FirstArgument, &task.SecondArgument, &operatorStr, &task.Status, &task.Result)
		if err != nil {
			log.Printf("Ошибка при сканировании задачи: %v", err)
			continue
		}

		if len(operatorStr) > 0 {
			task.Operator = rune(operatorStr[0])
		}

		updated := false

		if !IsNumeric(task.FirstArgument) {
			log.Printf("Аргумент 1 не является числом: %s\n", task.FirstArgument)
			id, err := strconv.ParseInt(strings.TrimPrefix(task.FirstArgument, "id"), 10, 64)
			if err != nil {
				log.Printf("Ошибка в парсинге %s", task.FirstArgument)
				continue
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
			updated = true
		}

		if !IsNumeric(task.SecondArgument) {
			log.Printf("Аргумент 2 не является числом: %s\n", task.SecondArgument)
			id, err := strconv.ParseInt(strings.TrimPrefix(task.SecondArgument, "id"), 10, 64)
			if err != nil {
				log.Printf("Ошибка в парсинге %s", task.SecondArgument)
				continue
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
			updated = true
		}

		if updated {
			_, err = q.db.Exec(
				"UPDATE tasks SET first_argument = ?, second_argument = ? WHERE id = ?",
				task.FirstArgument, task.SecondArgument, task.ID,
			)
			if err != nil {
				log.Printf("Ошибка при обновлении задачи %d: %v", task.ID, err)
				continue
			}
		}
	}

	if err = rows.Err(); err != nil {
		return fmt.Errorf("ошибка при итерации по строкам: %v", err)
	}

	return nil
}

func (q *Queue) generateTasksFromRPN(output []string) ([]shared.Task, map[int]string, error) {
	var tasks []shared.Task
	var operandStack []string
	taskIDs := make(map[int]string)

	// Получаем максимальный ID существующих задач
	var nextID int64
	err := q.db.QueryRow("SELECT COALESCE(MAX(id), 0) + 1 FROM tasks").Scan(&nextID)
	if err != nil {
		log.Printf("Ошибка при получении следующего ID задачи: %v", err)
		nextID = 1
	}

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

	// Добавляем задачи в базу данных
	for _, task := range tasks {
		_, err := q.db.Exec(
			"INSERT INTO tasks (id, first_argument, second_argument, operator, status, result) VALUES (?, ?, ?, ?, ?, ?)",
			task.ID, task.FirstArgument, task.SecondArgument, string(task.Operator), task.Status, task.Result,
		)
		if err != nil {
			log.Printf("Ошибка при добавлении задачи %d: %v", task.ID, err)
			return 0, err
		}
	}

	// Получаем максимальный ID существующих выражений
	var nextExprID int64
	err = q.db.QueryRow("SELECT COALESCE(MAX(id), 0) + 1 FROM expressions").Scan(&nextExprID)
	if err != nil {
		log.Printf("Ошибка при получении следующего ID выражения: %v", err)
		nextExprID = 1
	}

	// Добавляем выражение в базу данных
	expr := shared.Expression{
		ID:     nextExprID,
		Status: false,
		Result: fmt.Sprintf("id%d", tasks[len(tasks)-1].ID),
	}

	_, err = q.db.Exec(
		"INSERT INTO expressions (id, status, result) VALUES (?, ?, ?)",
		expr.ID, expr.Status, expr.Result,
	)
	if err != nil {
		log.Printf("Ошибка при добавлении выражения %d: %v", expr.ID, err)
		return 0, err
	}

	return expr.ID, nil
}
