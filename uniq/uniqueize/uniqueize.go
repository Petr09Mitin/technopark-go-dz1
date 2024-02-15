package uniqueize

import (
	"strings"
	"unicode/utf8"
)

type Flags struct {
	Count        *bool
	Duplicate    *bool
	Unduplicated *bool
	SkipFields   *uint
	SkipRunes    *uint
	IgnoreCase   *bool
}

type LineData struct {
	Line  string
	Count uint
}

func shouldAppend(lineData LineData, flags Flags) bool {
	switch {
	case *flags.Count:
		return true
	case *flags.Duplicate && lineData.Count > 1:
		return true
	case *flags.Unduplicated && lineData.Count == 1:
		return true
	case !*flags.Count && !*flags.Duplicate && !*flags.Unduplicated:
		return true
	}

	return false
}

func Uniqueize(lines []string, flags Flags) (linesData []LineData) {
	prevLine := ""
	prevCurrLine := ""
	var currCount uint = 0
	for _, line := range lines {
		currLine := line
		if *flags.SkipFields > 0 && *flags.SkipFields < uint(utf8.RuneCountInString(currLine)) {
			currLine = strings.Join(strings.Fields(currLine)[*flags.SkipFields:], " ")
		}

		if *flags.SkipRunes > 0 && *flags.SkipRunes < uint(utf8.RuneCountInString(currLine)) {
			currLine = string([]rune(currLine)[*flags.SkipRunes:])
		}

		if *flags.IgnoreCase {
			currLine = strings.ToLower(currLine)
		}

		if currLine == prevCurrLine {
			currCount++
		} else {
			lineData := LineData{Line: prevLine, Count: currCount}
			if currCount != 0 && shouldAppend(lineData, flags) {
				linesData = append(linesData, lineData)
			}
			currCount = 1
			prevCurrLine = currLine
			prevLine = line
		}
	}

	lineData := LineData{Line: prevLine, Count: currCount}
	if currCount != 0 && shouldAppend(lineData, flags) {
		linesData = append(linesData, lineData)
	}

	return
}
