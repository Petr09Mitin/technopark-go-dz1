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
	SkipFields   *int
	SkipRunes    *int
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
	if count > 1 || *flags.SkipFields < 0 || *flags.SkipRunes < 0 {
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
	var currCount uint
	hasPrevLineChanged := false
	for _, line := range lines {
		line = strings.Trim(line, " ")
		currLine := line
		fields := strings.Fields(currLine)
		if *flags.SkipFields < len(fields) {
			currLine = strings.Join(fields[*flags.SkipFields:], " ")
		} else {
			currLine = ""
		}

		if *flags.SkipRunes < utf8.RuneCountInString(currLine) {
			currLine = string([]rune(currLine)[*flags.SkipRunes:])
		} else {
			currLine = ""
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
			hasPrevLineChanged = true
		}
	}

	lineData := LineData{Line: prevLine, Count: currCount}
	if !hasPrevLineChanged && len(lines) > 0 {
		lineData.Line = lines[0]
	}

	if currCount != 0 && shouldAppend(lineData, flags) {
		linesData = append(linesData, lineData)
	}

	return
}
