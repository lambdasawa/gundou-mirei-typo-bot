package validate

import (
	"strings"
)

var (
	firstNames = [][]string{
		{"郡", "群"},
		{"道", "堂"},
	}
	lastNames = [][]string{
		{"美"},
		{"玲", "怜", "鈴", "麗"},
	}
	suffix     = "先生"
	validTexts = map[string]interface{}{
		"郡道美玲": struct{}{},
		"郡道先生": struct{}{},
		"美玲先生": struct{}{},
	}
)

func IsValidText(text string) bool {
	invalidTexts := GetInvalidTexts()
	for _, t := range invalidTexts {
		if strings.Contains(text, t) {
			return false
		}
	}

	return true
}

func GetInvalidTexts() []string {
	var ts []string

	invalidFullNames := makeTextPermutations(
		append(
			append([][]string{},
				firstNames...,
			),
			lastNames...,
		),
	)
	ts = append(ts, invalidFullNames...)

	ts = append(ts, makeTextPermutations([][]string{
		makeTextPermutations(firstNames),
		{suffix},
	})...)

	// ts = append(ts, makeTextPermutations([][]string{
	// 	makeTextPermutations(lastNames),
	// 	{suffix},
	// })...)

	var its []string
	for _, t := range ts {
		if _, ok := validTexts[t]; !ok {
			its = append(its, t)
		}
	}

	return its
}

func makeTextPermutations(textTable [][]string) []string {
	return makeTextsInternal(textTable[0], textTable[1:])
}

func makeTextsInternal(texts []string, textTable [][]string) []string {
	if len(textTable) == 0 {
		return texts
	}

	var result []string
	for _, prefix := range texts {
		for _, t := range textTable[0] {
			result = append(result, prefix+t)
		}
	}

	return makeTextsInternal(result, textTable[1:])
}
