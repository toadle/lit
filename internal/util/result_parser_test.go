package util

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/samber/lo"
)

func TestParserCorrectFormat(t *testing.T) {
	itemFormat := "(?P<data>.+):(?P<label>.+)"
	resultStr := "/Applications/Keynote.app:Keynote"
	parsedResult := ParseResult(resultStr, itemFormat)

	assert.Equal(t, len(parsedResult.Data), 2, "has correct number of result-data")

	dataKeys := lo.Keys[string,string](parsedResult.Data)
	assert.Equal(t, dataKeys[0], "data", "has extracted the correct keys")
	assert.Equal(t, dataKeys[1], "label", "has extracted the correct keys")

	assert.Equal(t, parsedResult.Data["data"], "/Applications/Keynote.app", "has extracted the correct values")
	assert.Equal(t, parsedResult.Data["label"], "Keynote", "has extracted the correct values")

	assert.Equal(t, parsedResult.Unparsed, "/Applications/Keynote.app:Keynote", "contains the original result")
}

func TestParserEmptyFormat(t *testing.T) {
	itemFormat := ""
	resultStr := "/Applications/Keynote.app"
	parsedResult := ParseResult(resultStr, itemFormat)

	assert.Equal(t, len(parsedResult.Data), 0, "has correct number of result-data")
	assert.Equal(t, parsedResult.Unparsed, "/Applications/Keynote.app", "contains the original result")
}
