package utils

import (
	"fmt"
	"math/rand"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	alphabet        = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	numerals        = "0123456789"
	special         = "~!@#$%^&*()-_+={}[]\\|<,>.?/\"';:`"
	asciiChars      = alphabet + numerals + special
	DefaultLanguage = "idn"
	EnglishLanguage = "en"
)

func init() {
	rand.New(rand.NewSource(time.Now().UnixNano() + int64(os.Getpid())))
}

func Generate(template string) (string, error) {
	re := regexp.MustCompile(`\[([a-zA-Z0-9\-\\]+)\](\{([0-9]+)\})?`)
	return re.ReplaceAllStringFunc(template, func(s string) string {
		match := re.FindStringSubmatch(s)
		ranges := match[1]
		lengthStr := match[3]
		length, err := strconv.Atoi(lengthStr)
		if err != nil {
			length = 1
		}
		var sb strings.Builder
		for i := 0; i < length; i++ {
			if err := generateChar(&sb, ranges); err != nil {
				return ""
			}
		}
		return sb.String()
	}), nil
}

func generateChar(sb *strings.Builder, ranges string) error {
	switch ranges {
	case `\w`:
		ranges = asciiChars
	case `\d`:
		ranges = numerals
	}
	re := regexp.MustCompile(`\\?([a-zA-Z0-9])-\\?([a-zA-Z0-9])`)
	ranges = re.ReplaceAllStringFunc(ranges, func(s string) string {
		match := re.FindStringSubmatch(s)
		from, to := match[1][0], match[2][0]
		if from > to {
			from, to = to, from
		}
		return alphabetSlice(from, to)
	})
	if len(ranges) == 0 {
		return fmt.Errorf("empty range in expression: %s", ranges)
	}
	_, err := sb.WriteString(string(ranges[seedAndReturnRandom(len(ranges))]))
	return err
}

func alphabetSlice(from, to byte) string {
	leftPos := strings.IndexByte(asciiChars, from)
	rightPos := strings.LastIndexByte(asciiChars, to)
	if leftPos > rightPos {
		return ""
	}
	return asciiChars[leftPos : rightPos+1]
}

func seedAndReturnRandom(n int) int {
	return rand.Intn(n)
}
