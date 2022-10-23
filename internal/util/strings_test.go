package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSpecialCharacterRemoval(t *testing.T) {
	assert := assert.New(t)

	str := " Expressions /Applications/Expressions.app"

	assert.Equal("Expressions Applications Expressions app", RemoveSpecialCharacters(str), "returns the cleared up string")
}

func TestFindAllOccurrencesOfCharacters(t *testing.T) {
	assert := assert.New(t)

	str := "This is a test string 123"
	toFind := "ti2üö " // <- spaces should not be matched and case should not matter

	assert.Equal([]int{0, 2, 5, 10, 13, 16, 18, 23}, FindAllOccurrencesOfCharacters(str, toFind), "returns the right indices")
}
