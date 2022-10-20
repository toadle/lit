package shell

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetterWithSensibleInput(t *testing.T) {
	assert := assert.New(t)

	cmdStr := "trans --brief {input} {someother} {notpresent}"
	var params = map[string]string{
		"input":     "test",
		"someother": "test2",
	}

	cmdStrParameterized := SetCommandParameters(cmdStr, params)

	assert.Equal("trans --brief test test2", cmdStrParameterized, "sets params correctly")
}
