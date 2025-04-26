package task

import (
	"strings"
	"unicode"

	"github.com/nktauserum/web-calculation/shared"
	"github.com/nktauserum/web-calculation/shared/errors"
)

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

// Проверяет, является ли токен оператором
func isOperator(token string) bool {
	switch token {
	case "+", "-", "*", "/":
		return true
	default:
		return false
	}
}

// IsNumeric проверяет, является ли строка числом
func IsNumeric(s string) bool {
	for _, r := range s {
		if !unicode.IsDigit(r) && r != '.' && r != ',' && r != '-' {
			return false
		}
	}
	return true
}

// tokenize разбивает выражение на токены
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

// Complete проверяет, готова ли задача к выполнению
func Complete(task shared.Task) bool {
	return IsNumeric(task.FirstArgument) && IsNumeric(task.SecondArgument)
}
