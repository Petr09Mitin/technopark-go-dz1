package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/Petr09Mitin/technopark-go-dz1/uniq/uniqueize"
)

// Arguments represents the input and output files.
type Arguments struct {
	InputFile  string
	OutputFile string
}

func handleError(err error) {
	const docString string = `
Usage: 
	uniq [-c | -d | -u] [-i] [-f fields] [-s chars] [input_file [output_file]]

Parameters:

	-c: count number of occurrences

	-d: print only duplicate lines

	-u: print only unique lines

	-i: ignore case differences

	-f fields: avoid comparing the first fields fields

	-s chars: avoid comparing the first chars characters

	input_file: file to read from
	
	output_file: file to write to
`
	if err != nil {
		fmt.Println(err)
		if err.Error() == "invalid flags" {
			fmt.Print(docString)
		}
		os.Exit(0)
	}
}

// ValidateArguments validates the input and output files.
func ValidateArguments(arguments Arguments) error {
	if arguments.InputFile != "" {
		if _, err := os.Stat(arguments.InputFile); os.IsNotExist(err) {
			return errors.New("input file does not exist")
		}

		if arguments.OutputFile != "" {
			if _, err := os.Stat(arguments.OutputFile); os.IsNotExist(err) {
				return errors.New("output file does not exist")
			}
		}
	}

	return nil
}

// ParseFlags parses the flags from the command line arguments and returns the flags.
func ParseFlags() (flags uniqueize.Flags) {
	flags.Count = flag.Bool("c", false, "count number of occurrences")
	flags.Duplicate = flag.Bool("d", false, "print only duplicate lines")
	flags.Unduplicated = flag.Bool("u", false, "print only unique lines")
	flags.SkipFields = flag.Uint("f", 0, "avoid comparing the first N fields")
	flags.SkipRunes = flag.Uint("s", 0, "avoid comparing the first N characters")
	flags.IgnoreCase = flag.Bool("i", false, "ignore case differences")

	flag.Parse()

	return
}

// ParseInAndOutFiles parses the input and output files from the command line arguments and returns the files
// or stdin and stdout if files are not specified.
func ParseInAndOutFiles() (inputFile, outputFile *os.File, argumentsErr error) {
	var arguments Arguments
	arguments.InputFile = flag.Arg(0)
	arguments.OutputFile = flag.Arg(1)

	argumentsErr = ValidateArguments(arguments)

	if argumentsErr != nil {
		return
	}

	inputFile = os.Stdin
	outputFile = os.Stdout

	if arguments.InputFile != "" {
		inputFile, _ = os.Open(arguments.InputFile)
	}

	if arguments.OutputFile != "" {
		outputFile, _ = os.Create(arguments.OutputFile)
	}

	return
}

// GetReaderAndWriter returns a reader and a writer for the given files or stdin and stdout if files are not specified.
func GetReaderAndWriter(inputFile, outputFile *os.File) (reader *bufio.Reader, writer *bufio.Writer) {
	reader = bufio.NewReader(inputFile)
	writer = bufio.NewWriter(outputFile)

	return
}

// ReadInput reads the input from the reader and returns it as an array of strings.
func ReadInput(reader *bufio.Reader) (lines []string, err error) {
	for {
		line, readingErr := reader.ReadString('\n')
		if len(line) == 0 && readingErr != nil {
			if readingErr == io.EOF {
				break
			}
			err = readingErr
			return
		}

		line = strings.TrimSuffix(line, "\n")
		line = strings.TrimSuffix(line, "\r")
		lines = append(lines, line)

		if readingErr != nil {
			if readingErr == io.EOF {
				break
			}
			err = readingErr
			return
		}
	}
	return
}

// WriteOutput writes the linesData array to the writer in format specified by flags.
func WriteOutput(flags uniqueize.Flags, writer *bufio.Writer, linesData []uniqueize.LineData) (err error) {
	for i, lineData := range linesData {
		switch {
		case *flags.Count:
			fmt.Fprintf(writer, "%d %s", lineData.Count, lineData.Line)
		case *flags.Duplicate && lineData.Count > 1:
			fmt.Fprintf(writer, "%s", lineData.Line)
		case *flags.Unduplicated && lineData.Count == 1:
			fmt.Fprintf(writer, "%s", lineData.Line)
		case !*flags.Count && !*flags.Duplicate && !*flags.Unduplicated:
			fmt.Fprintf(writer, "%s", lineData.Line)
		}

		if i != len(linesData)-1 {
			fmt.Fprintf(writer, "\n")
		}
	}

	writer.Flush()
	return
}

func main() {
	flags := ParseFlags()

	inputFile, outputFile, argumentsErr := ParseInAndOutFiles()

	handleError(argumentsErr)

	reader, writer := GetReaderAndWriter(inputFile, outputFile)

	lines, err := ReadInput(reader)
	inputFile.Close()
	handleError(err)

	linesData, err := uniqueize.Uniqueize(lines, flags)
	handleError(err)

	WriteOutput(flags, writer, linesData)
	outputFile.Close()
}
