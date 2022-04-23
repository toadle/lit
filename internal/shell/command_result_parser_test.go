package shell

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/samber/lo"
)

func TestParserCorrectFormat(t *testing.T) {
	assert := assert.New(t)

	itemFormat := "(?P<data>.+):(?P<label>.+)"
	resultStr := "/Applications/Keynote.app:Keynote"
	parsedResult := ParseCommandResult(resultStr, itemFormat)

	assert.Equal(len(parsedResult.Params), 2, "has correct number of result-data")

	dataKeys := lo.Keys[string,string](parsedResult.Params)
	assert.Equal(lo.Contains[string](dataKeys, "data"), true, "has extracted the 'data' key")
	assert.Equal(lo.Contains[string](dataKeys, "label"), true, "has extracted the 'label' key")

	assert.Equal(parsedResult.Params["data"], "/Applications/Keynote.app", "has extracted the correct values")
	assert.Equal(parsedResult.Params["label"], "Keynote", "has extracted the correct values")

	assert.Equal(parsedResult.Unparsed, "/Applications/Keynote.app:Keynote", "contains the original result")
}

func TestParserEmptyFormat(t *testing.T) {
	assert := assert.New(t)

	itemFormat := "(?P<data>.+):(?P<label>.+)"
	resultStr := "/Applications/Keynote.app"
	parsedResult := ParseCommandResult(resultStr, itemFormat)

	assert.Equal(len(parsedResult.Params), 0, "has correct number of result-data")
	assert.Equal(parsedResult.Unparsed, "/Applications/Keynote.app", "contains the original result")
}
