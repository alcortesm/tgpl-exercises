// Dup2 prints the count and text of lines that appear more than once
// in the input.  It reads from stdin or from a list of named files.
package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func main() {
	counts := NewFreqCountWithTags()

	files := os.Args[1:]
	if len(files) == 0 {
		countLines(io.Reader(os.Stdin), "/dev/stdin", counts)
	} else {
		// TODO: parallelize this using gorutines
		for _, arg := range files {
			f, err := os.Open(arg)
			if err != nil {
				fmt.Fprintf(os.Stderr, "dup2: %v\n", err)
				continue
			}

			countLines(f, arg, counts)

			err = f.Close()
			if err != nil {
				fmt.Fprintf(os.Stderr, "dup2: %v\n", err)
				continue
			}
		}
	}

	printResults(counts)
}

func countLines(f io.Reader, fileName string, counts *FreqCountWithTags) {
	input := bufio.NewScanner(f)
	for input.Scan() {
		counts.Add(input.Text(), fileName)
	}

	if err := input.Err(); err != nil {
		fmt.Fprintln(os.Stdout, input.Err())
		os.Exit(0)
	}
}

func printResults(results *FreqCountWithTags) {
	lines, counts, files := results.GetAll()

	for i := 0; i < len(lines); i++ {
		line := lines[i]
		count := counts[i]
		fileList := files[i]
		fmt.Println(line, count, fileList)
	}
}
