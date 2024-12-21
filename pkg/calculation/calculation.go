package calculation

import (
	"strconv"
	"unicode"

	"github.com/veliashev/rpn/pkg/errors"
)

// Calc принимает строку с математическим выражением и возвращает результат вычисления
func Calc(expression string) (float64, error) {
	// Разбиваем выражение на токены (числа и операторы)
	tokens := tokenize(expression)
	// Вычисляем результат на основе полученных токенов
	return evaluate(tokens)
}

// tokenize разбивает строку выражения на отдельные токены (числа и операторы)
func tokenize(expression string) []string {
	var tokens []string // слайс для хранения токенов
	var current string  // текущее накапливаемое число

	// Перебираем каждый символ в выражении
	for _, char := range expression {
		if unicode.IsDigit(char) || char == '.' {
			current += string(char)
		} else if char == ',' {
			current += string(".")
		} else {
			if current != "" {
				tokens = append(tokens, current)
				current = ""
			}
			if !unicode.IsSpace(char) {
				tokens = append(tokens, string(char))
			}
		}
	}

	if current != "" {
		tokens = append(tokens, current)
	}
	return tokens
}

// Вычисляет результат выражения
func evaluate(tokens []string) (float64, error) {
	var numbers []float64  // стек для чисел
	var operators []string // стек для операторов

	// Обрабатываем каждый токен
	for _, token := range tokens {
		switch token {
		case "(":
			// Открывающая скобка просто добавляется в стек операторов
			operators = append(operators, token)
		case ")":
			// При закрывающей скобке вычисляем все операторы до открывающей
			for len(operators) > 0 && operators[len(operators)-1] != "(" {
				if err := applyOperator(&numbers, &operators); err != nil {
					return 0, err
				}
			}
			if len(operators) == 0 {
				return 0, errors.ErrMismatchedParentheses
			}
			// Удаляем открывающую скобку
			operators = operators[:len(operators)-1]
		case "+", "-", "*", "/":
			// Для операторов вычисляем все операторы с большим или равным приоритетом
			for len(operators) > 0 && precedence(operators[len(operators)-1]) >= precedence(token) {
				if err := applyOperator(&numbers, &operators); err != nil {
					return 0, err
				}
			}
			operators = append(operators, token)
		default:
			// Преобразуем строку в число
			num, err := strconv.ParseFloat(token, 64)
			if err != nil {
				return 0, errors.ErrInvalidNumber
			}
			numbers = append(numbers, num)
		}
	}

	// Применяем оставшиеся операторы
	for len(operators) > 0 {
		if err := applyOperator(&numbers, &operators); err != nil {
			return 0, err
		}
	}

	// Проверяем корректность результата
	if len(numbers) != 1 {
		return 0, errors.ErrInvalidExpression
	}
	return numbers[0], nil
}

// precedence возвращает приоритет оператора
func precedence(op string) int {
	switch op {
	case "+", "-":
		return 1 // низкий приоритет
	case "*", "/":
		return 2 // высокий приоритет
	default:
		return 0 // для скобок и неизвестных операторов
	}
}

// applyOperator применяет оператор к двум последним числам в стеке
func applyOperator(numbers *[]float64, operators *[]string) error {
	if len(*numbers) < 2 {
		return errors.ErrNotEnoughOperands
	}
	// Извлекаем два последних числа
	b, a := (*numbers)[len(*numbers)-1], (*numbers)[len(*numbers)-2]
	*numbers = (*numbers)[:len(*numbers)-2]

	// Извлекаем оператор
	op := (*operators)[len(*operators)-1]
	*operators = (*operators)[:len(*operators)-1]

	// Выполняем операцию
	var result float64
	switch op {
	case "+":
		result = a + b
	case "-":
		result = a - b
	case "*":
		result = a * b
	case "/":
		if b == 0 {
			return errors.ErrDivisionByZero
		}
		result = a / b
	default:
		return errors.ErrUnknownOperator
	}

	// Добавляем результат обратно в стек чисел
	*numbers = append(*numbers, result)
	return nil
}
