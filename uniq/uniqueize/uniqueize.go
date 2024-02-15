package uniqueize

import (
	"errors"
	"strings"
	"unicode/utf8"
)

// Flags represents the flags for the uniq command
// Count: count number of occurrences (-c)
// Duplicate: print only duplicate lines (-d)
// Unduplicated: print only unique lines (-u)
// SkipFields: avoid comparing the first N fields (-f num)
// SkipRunes: avoid comparing the first N characters (-s num)
// IgnoreCase: ignore case differences (-i)
type Flags struct {
	Count        *bool
	Duplicate    *bool
	Unduplicated *bool
	SkipFields   *uint
	SkipRunes    *uint
	IgnoreCase   *bool
}

// LineData represents the line and its appearance count.
type LineData struct {
	Line  string
	Count uint
}

// validateFlags checks so that only one of the flags -c, -d or -u is set.
func validateFlags(flags Flags) error {
	count := 0
	if *flags.Count {
		count++
	}
	if *flags.Duplicate {
		count++
	}
	if *flags.Unduplicated {
		count++
	}
	if count > 1 {
		return errors.New("invalid flags")
	}

	return nil
}

// shouldAppend checks if the line should be appended to the output according to the flags.
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

// Uniqueize transforms input lines into []lineData according to the flags.
func Uniqueize(lines []string, flags Flags) (linesData []LineData, err error) {
	flagsErr := validateFlags(flags)
	if flagsErr != nil {
		err = flagsErr
		return
	}

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
