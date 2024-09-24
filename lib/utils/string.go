package utils

import (
	"strings"
)

func ConvertToSnakeCase(s string) string {
	lower := strings.ToLower(s)
	snakeCase := strings.ReplaceAll(lower, " ", "_")

	return snakeCase
}
