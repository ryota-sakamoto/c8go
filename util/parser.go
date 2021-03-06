package util

import (
	"errors"
	"strconv"
	"unicode"
)

func ParseInt(s *string) (int, error) {
	t := []rune{}
	index := 0
	for _, c := range *s {
		if unicode.IsDigit(c) {
			t = append(t, c)
			index++
		} else {
			break
		}
	}

	if index == 0 {
		return 0, errors.New("not number")
	}

	*s = (*s)[index:]
	return strconv.Atoi(string(t))
}

func IsAlnum(c byte) bool {
	return ('a' <= c && c <= 'z') ||
		('A' <= c && c <= 'Z') ||
		('0' <= c && c <= '9') ||
		('_' == c)
}
