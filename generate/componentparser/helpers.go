package componentparser

import (
	"strings"
)

/*
Function for parsing values
*/
func extractOptions(structString string) map[string]string {
	optionMap := map[string]string{}

	for index := strings.Index(structString, "::"); index != -1; index = strings.Index(structString, "::") {
		before, after := structString[:index], structString[index+2:]

		beforeSplits := multiSplit(before, []string{" ", "\n", "\t"})
		afterSplits := multiSplit(after, []string{" ", "\n", "\t"})

		optionMap[beforeSplits[len(beforeSplits)-1]] = afterSplits[0]
		structString = after
	}

	return optionMap
}

/*
Function for splitting by multiple values in sequence
*/
func multiSplit(data string, seps []string) []string {
	splits := []string{data}
	var newSplits []string

	for _, sep := range seps {
		newSplits = []string{}
		for _, splitData := range splits {
			newSplits = append(newSplits, strings.Split(splitData, sep)...)
		}
		splits = newSplits
	}

	return splits
}
