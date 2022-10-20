package shell

import (
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func TestParserCorrectFormat(t *testing.T) {
	assert := assert.New(t)

	itemFormat := "(?P<data>.+):(?P<label>.+)"
	resultStr := "/Applications/Keynote.app:Keynote"
	parsedResult := ParseCommandResult(resultStr, itemFormat)

	assert.Equal(len(parsedResult.Params), 2, "has correct number of result-data")

	dataKeys := lo.Keys[string, string](parsedResult.Params)
	assert.Equal(true, lo.Contains[string](dataKeys, "data"), "has extracted the 'data' key")
	assert.Equal(true, lo.Contains[string](dataKeys, "label"), "has extracted the 'label' key")

	assert.Equal("/Applications/Keynote.app", parsedResult.Params["data"], "has extracted the correct values")
	assert.Equal("Keynote", parsedResult.Params["label"], "has extracted the correct values")

	assert.Equal("/Applications/Keynote.app:Keynote", parsedResult.Unparsed, "contains the original result")
}

func TestParserEmptyFormat(t *testing.T) {
	assert := assert.New(t)

	itemFormat := "(?P<data>.+):(?P<label>.+)"
	resultStr := "/Applications/Keynote.app"
	parsedResult := ParseCommandResult(resultStr, itemFormat)

	assert.Equal(0, len(parsedResult.Params), "has correct number of result-data")
	assert.Equal("/Applications/Keynote.app", parsedResult.Unparsed, "contains the original result")
}
