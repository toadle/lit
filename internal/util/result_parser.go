package util

import (
	"regexp"
)

type ParsedResult struct {
	Unparsed string
	Data map[string]string
}

func ParseResult(resultStr string, itemFormat string) ParsedResult {
	re := regexp.MustCompile(itemFormat)
	groupNames := re.SubexpNames()

	resultData := map[string]string{}

	for _, match := range re.FindAllStringSubmatch(resultStr, -1) {
		for groupIdx, group := range match {
			name := groupNames[groupIdx]
			if name != "" {
				resultData[name] = group
			}
		}
	}

	return ParsedResult{Unparsed: resultStr, Data: resultData}
}
