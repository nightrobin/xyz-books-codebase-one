package method

import (
	"strings"
	"strconv"
)

func IsbnValidator(isbn string) bool {
	var isbnArr = strings.Split(isbn, "")
	switch len(isbnArr) {
		case 13:
			return validateIsbn13(isbnArr)
		case 10:
			return validateIsbn10(isbnArr)
		default:
			return false
	}
}

func validateIsbn10(isbn []string) bool {
	var isbnSum int
	for i := 0; i < len(isbn); i++ {
		var num int
		var err error
		if isbn[i] == "X" {
			num = 10
		} else {
			num, err = strconv.Atoi(isbn[i])
			if err != nil {
				return false
			}
		}
		isbnSum += (10 - i) * num
	}

	if isbnSum % 11 != 0 {
		return false
	}
	return true
}

func validateIsbn13(isbn []string) bool {
	var isbnSum int
	for i := 0; i < len(isbn); i++ {
		var multiplier int = 3
		if i % 2 == 0 {
			multiplier = 1
		}
		num, err := strconv.Atoi(isbn[i])
		if err != nil {
			return false
		}
		isbnSum += num * multiplier
	}

	if isbnSum % 10 != 0 {
		return false
	}
	return true
}