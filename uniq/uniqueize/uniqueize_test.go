package uniqueize_test

import (
	"testing"

	. "github.com/Petr09Mitin/technopark-go-dz1/uniq/uniqueize"
	"github.com/stretchr/testify/assert"
)

func newTrue() *bool {
	b := true
	return &b
}

func newUint(i uint) *uint {
	return &i
}

var successfulTests = map[string]struct {
	lines  []string
	flags  Flags
	output []LineData
}{
	"no flags": {
		lines: []string{
			"I love music.",
			"I love music.",
			"I love music.",
			"",
			"I love music of Kartik.",
			"I love music of Kartik.",
			"Thanks.",
			"I love music of Kartik.",
			"I love music of Kartik.",
		},
		flags: Flags{
			Count:        new(bool),
			Duplicate:    new(bool),
			Unduplicated: new(bool),
			SkipFields:   new(uint),
			SkipRunes:    new(uint),
			IgnoreCase:   new(bool),
		},
		output: []LineData{
			{Line: "I love music.", Count: 3},
			{Line: "", Count: 1},
			{Line: "I love music of Kartik.", Count: 2},
			{Line: "Thanks.", Count: 1},
			{Line: "I love music of Kartik.", Count: 2},
		},
	},
	"-c flag set": {
		lines: []string{
			"I love music.",
			"I love music.",
			"I love music.",
			"",
			"I love music of Kartik.",
			"I love music of Kartik.",
			"Thanks.",
			"I love music of Kartik.",
			"I love music of Kartik.",
		},
		flags: Flags{
			Count:        newTrue(),
			Duplicate:    new(bool),
			Unduplicated: new(bool),
			SkipFields:   new(uint),
			SkipRunes:    new(uint),
			IgnoreCase:   new(bool),
		},
		output: []LineData{
			{Line: "I love music.", Count: 3},
			{Line: "", Count: 1},
			{Line: "I love music of Kartik.", Count: 2},
			{Line: "Thanks.", Count: 1},
			{Line: "I love music of Kartik.", Count: 2},
		},
	},
	"-d flag set": {
		lines: []string{
			"I love music.",
			"I love music.",
			"I love music.",
			"",
			"I love music of Kartik.",
			"I love music of Kartik.",
			"Thanks.",
			"I love music of Kartik.",
			"I love music of Kartik.",
		},
		flags: Flags{
			Count:        new(bool),
			Duplicate:    newTrue(),
			Unduplicated: new(bool),
			SkipFields:   new(uint),
			SkipRunes:    new(uint),
			IgnoreCase:   new(bool),
		},
		output: []LineData{
			{Line: "I love music.", Count: 3},
			{Line: "I love music of Kartik.", Count: 2},
			{Line: "I love music of Kartik.", Count: 2},
		},
	},
	"-u flag set": {
		lines: []string{
			"I love music.",
			"I love music.",
			"I love music.",
			"",
			"I love music of Kartik.",
			"I love music of Kartik.",
			"Thanks.",
			"I love music of Kartik.",
			"I love music of Kartik.",
		},
		flags: Flags{
			Count:        new(bool),
			Duplicate:    new(bool),
			Unduplicated: newTrue(),
			SkipFields:   new(uint),
			SkipRunes:    new(uint),
			IgnoreCase:   new(bool),
		},
		output: []LineData{
			{Line: "", Count: 1},
			{Line: "Thanks.", Count: 1},
		},
	},
	"-i flag set": {
		lines: []string{
			"I LOVE MUSIC.",
			"I love music.",
			"I LoVe MuSiC.",
			"",
			"I love MuSIC of Kartik.",
			"I love music of kartik.",
			"Thanks.",
			"I love music of kartik.",
			"I love MuSIC of Kartik.",
		},
		flags: Flags{
			Count:        new(bool),
			Duplicate:    new(bool),
			Unduplicated: new(bool),
			SkipFields:   new(uint),
			SkipRunes:    new(uint),
			IgnoreCase:   newTrue(),
		},
		output: []LineData{
			{Line: "I LOVE MUSIC.", Count: 3},
			{Line: "", Count: 1},
			{Line: "I love MuSIC of Kartik.", Count: 2},
			{Line: "Thanks.", Count: 1},
			{Line: "I love music of kartik.", Count: 2},
		},
	},
	"-f flag set": {
		lines: []string{
			"We love music.",
			"I love music.",
			"They love music.",
			"",
			"I love music of Kartik.",
			"We love music of Kartik.",
			"Thanks.",
		},
		flags: Flags{
			Count:        new(bool),
			Duplicate:    new(bool),
			Unduplicated: new(bool),
			SkipFields:   newUint(1),
			SkipRunes:    new(uint),
			IgnoreCase:   new(bool),
		},
		output: []LineData{
			{Line: "We love music.", Count: 3},
			{Line: "", Count: 1},
			{Line: "I love music of Kartik.", Count: 2},
			{Line: "Thanks.", Count: 1},
		},
	},
	"-s flag set": {
		lines: []string{
			"I love music.",
			"A love music.",
			"C love music.",
			"",
			"I love music of Kartik.",
			"We love music of Kartik.",
			"Thanks.",
		},
		flags: Flags{
			Count:        new(bool),
			Duplicate:    new(bool),
			Unduplicated: new(bool),
			SkipFields:   new(uint),
			SkipRunes:    newUint(1),
			IgnoreCase:   new(bool),
		},
		output: []LineData{
			{Line: "I love music.", Count: 3},
			{Line: "", Count: 1},
			{Line: "I love music of Kartik.", Count: 1},
			{Line: "We love music of Kartik.", Count: 1},
			{Line: "Thanks.", Count: 1},
		},
	},
}

var failedTests = map[string]struct {
	lines  []string
	flags  Flags
	output []LineData
}{
	"invalid flags": {
		lines: []string{},
		flags: Flags{
			Count:        newTrue(),
			Duplicate:    newTrue(),
			Unduplicated: newTrue(),
			SkipFields:   newUint(1),
			SkipRunes:    newUint(1),
			IgnoreCase:   newTrue(),
		},
		output: []LineData{},
	},
}

func TestSuccessfulUniqueize(t *testing.T) {
	for name, test := range successfulTests {
		t.Run(name, func(t *testing.T) {
			result, err := Uniqueize(test.lines, test.flags)
			assert.Nil(t, err)
			assert.Equal(t, test.output, result)
		})
	}
}

func TestFailedUniqueize(t *testing.T) {
	for name, test := range failedTests {
		t.Run(name, func(t *testing.T) {
			_, err := Uniqueize(test.lines, test.flags)
			assert.NotNil(t, err)
		})
	}
}
