package list

import (
	"lit/internal/config"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLabelsWithoutSensibleResultData(t *testing.T) {
	assert := assert.New(t)

	sourceConfig := config.SearchConfig{
		Format: "",
		Labels: config.MultiLineLabelsConfig{
			Title:       "",
			Description: "",
		},
	}

	resultItem := NewSearchListItem("TestData", sourceConfig)

	assert.Equal("", resultItem.title(), "returns nothing as label")
}

func TestLabelsWithoutConfigValues(t *testing.T) {
	assert := assert.New(t)

	sourceConfig := config.SearchConfig{
		Format: "(?P<description>.+):(?P<title>.+)",
		Labels: config.MultiLineLabelsConfig{
			Title:       "",
			Description: "",
		},
	}

	resultItem := NewSearchListItem("TestData:TestLabel", sourceConfig)

	assert.Equal("TestLabel", resultItem.title(), "returns the corrent label")
}

func TestLabelsWithConfigValues(t *testing.T) {
	assert := assert.New(t)

	sourceConfig := config.SearchConfig{
		Format: "(?P<data>.+):(?P<label>.+)",
		Labels: config.MultiLineLabelsConfig{
			Title:       "This is the {label}",
			Description: "",
		},
	}

	resultItem := NewSearchListItem("TestData:TestLabel", sourceConfig)

	assert.Equal("This is the TestLabel", resultItem.title(), "returns the corrent label")
}
