package shell

import (
	"fmt"
	"strings"
	"regexp"
)

func SetCommandParameters(cmdStr string, params map[string]string) string {
	re := regexp.MustCompile("{(?P<param>[a-z]+)}")
	resultStr := cmdStr
	for _, match := range re.FindAllStringSubmatch(cmdStr, -1) {
		for _, paramName := range match {
			resultStr = strings.Replace(resultStr, fmt.Sprintf("{%s}", paramName), params[paramName], 1)
		}
	}
	return strings.TrimSpace(resultStr)
}
