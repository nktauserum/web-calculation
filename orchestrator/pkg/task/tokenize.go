package task

import "unicode"

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
