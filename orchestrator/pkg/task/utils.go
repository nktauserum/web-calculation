package task

import (
	"unicode"

	"github.com/veliashev/web-calculation/shared"
)

func IsNumeric(s string) bool {
	for _, r := range s {
		if !unicode.IsDigit(r) && r != '.' && r != ',' && r != '-' {
			return false
		}
	}
	return true
}
func Complete(task shared.Task) bool {
	return IsNumeric(task.FirstArgument) && IsNumeric(task.SecondArgument)
}
