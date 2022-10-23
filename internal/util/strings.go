package util

import (
	"regexp"
	"strings"
)

var nonAlphanumericRegex = regexp.MustCompile(`[^\p{L}\p{N} ]+`)

func RemoveSpecialCharacters(str string) string {
	str = nonAlphanumericRegex.ReplaceAllString(str, " ")
	str = strings.ReplaceAll(str, "  ", " ")
	str = strings.TrimSpace(str)
	return str
}

func FindAllOccurrencesOfCharacters(str string, toFind string) []int {
	targetRunes := []rune(strings.ToLower(str))
	findRunes := []rune(strings.ToLower(toFind))
	i := []int{}

	for idx, tr := range targetRunes {
		for _, fr := range findRunes {
			if fr != ' ' && tr == fr {
				i = append(i, idx)
			}
		}
	}

	return i
}
