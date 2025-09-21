package util

import (
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"strings"
)

func LowCase(input string) string {
	lCaser := cases.Lower(language.English)
	lowerCaseInput := lCaser.String(input)
	return lowerCaseInput
}

func UpperCaseFirstLetter(input string) string {
	return strings.ToUpper(string(input[0])) + string(input[1:])
}
