package util

import (
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func LowCase(input string) string {
	lCaser := cases.Lower(language.English)
	lowerCaseInput := lCaser.String(input)
	return lowerCaseInput
}
