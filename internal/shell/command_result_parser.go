package shell

import (
	"regexp"
)

type CommandResult struct {
	Unparsed string
	Params map[string]string
}

func ParseCommandResult(resultStr string, itemFormat string) CommandResult {
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

	return CommandResult{Unparsed: resultStr, Params: resultData}
}
