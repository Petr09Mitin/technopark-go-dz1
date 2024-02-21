package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/Petr09Mitin/technopark-go-dz1/uniq/uniqueize"
)

// Arguments represents the input and output files.
type Arguments struct {
	InputFile  string
	OutputFile string
}

const terminator = "ENDLINE"

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
	}
}

// ValidateArguments validates the input and output files.
func ValidateArguments(arguments Arguments) error {
	if arguments.InputFile != "" {
		if _, err := os.Stat(arguments.InputFile); os.IsNotExist(err) {
			return errors.New("input file does not exist")
		} else if err != nil {
			return err
		}

		if arguments.OutputFile != "" {
			if _, err := os.Stat(arguments.OutputFile); err != nil && !os.IsNotExist(err) {
				return err
			}
		}
	} else {
		fmt.Println("Enter lines one by one. When you're finished, enter", terminator)
	}

	return nil
}

// ParseFlags parses the flags from the command line arguments and returns the flags.
func ParseFlags() (flags uniqueize.Flags) {
	flags.Count = flag.Bool("c", false, "count number of occurrences")
	flags.Duplicate = flag.Bool("d", false, "print only duplicate lines")
	flags.Unduplicated = flag.Bool("u", false, "print only unique lines")
	flags.SkipFields = flag.Int("f", 0, "avoid comparing the first N fields")
	flags.SkipRunes = flag.Int("s", 0, "avoid comparing the first N characters")
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
		inputFile, argumentsErr = os.Open(arguments.InputFile)
	}

	if arguments.OutputFile != "" {
		outputFile, argumentsErr = os.Create(arguments.OutputFile)
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
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err = scanner.Err(); err != nil {
		return
	}

	return
}

// WriteOutput writes the linesData array to the writer in format specified by flags.
func WriteOutput(flags uniqueize.Flags, writer *bufio.Writer, linesData []uniqueize.LineData) (err error) {
	for _, lineData := range linesData {
		switch {
		case *flags.Count:
			_, err = fmt.Fprintln(writer, lineData.Count, lineData.Line)
		case *flags.Duplicate && lineData.Count > 1:
			_, err = fmt.Fprintln(writer, lineData.Line)
		case *flags.Unduplicated && lineData.Count == 1:
			_, err = fmt.Fprintln(writer, lineData.Line)
		case !*flags.Count && !*flags.Duplicate && !*flags.Unduplicated:
			_, err = fmt.Fprintln(writer, lineData.Line)
		}

		if err != nil {
			return
		}
	}

	writer.Flush()
	return
}

func main() {
	flags := ParseFlags()

	inputFile, outputFile, argumentsErr := ParseInAndOutFiles()

	if argumentsErr != nil {
		handleError(argumentsErr)
		return
	}

	reader, writer := GetReaderAndWriter(inputFile, outputFile)
	defer inputFile.Close()
	defer outputFile.Close()

	lines, readErr := ReadInput(reader)
	if readErr != nil {
		handleError(readErr)
		return
	}

	linesData, uniqErr := uniqueize.Uniqueize(lines, flags)
	if uniqErr != nil {
		handleError(uniqErr)
		return
	}

	writeErr := WriteOutput(flags, writer, linesData)
	if writeErr != nil {
		handleError(writeErr)
		return
	}
}
