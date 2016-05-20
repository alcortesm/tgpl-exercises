package main

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/alcortesm/sample"
)

const (
	inputSize       = 1000
	runs            = 1000
	confidenceLevel = 0.95
)

func main() {
	input := sliceOfManyShortStrings(inputSize)
	benchmarks := runAllBenchmarks(input)
	printResults(benchmarks)

	os.Exit(0)
}

func sliceOfManyShortStrings(howMany int) []string {
	input := make([]string, howMany)
	for i := 0; i < howMany; i++ {
		input[i] = fmt.Sprint(i)
	}

	return input
}

func runAllBenchmarks(input []string) []result {
	results := make([]result, len(versions))
	var err error
	for i, fn := range versions {
		results[i], err = benchmark(fn, input)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}

	return results
}

type result [2]time.Duration

type concat func([]string) string

var versions = []concat{v1, v2, v3, v4, v5}
var descriptions = []string{
	"+ string operator",
	"same with range",
	"strings.Join",
	"bytes.Buffer",
	`Sprintf("%v")`,
}

// Runs `f` function over `input` a number of times (`runs`) and
// calculates the confidence intervals of the mean duration of each run
// for a `confidence` confidence level.
func benchmark(fn concat, input []string) (result, error) {
	var execTimes = make([]float64, runs)
	for i := 0; i < runs; i++ {
		start := time.Now()
		fn(input)
		elapsed := time.Since(start)
		execTimes[i] = float64(elapsed)
	}

	meanInterval, err := sample.MeanConfidenceIntervals(execTimes, confidenceLevel)
	if err != nil {
		return result{}, err
	}

	asDuration := result{
		time.Duration(meanInterval[0]),
		time.Duration(meanInterval[1]),
	}

	return asDuration, nil
}

const (
	realSep = " "
)

func v1(input []string) string {
	var output, sep string
	for i := 0; i < len(input); i++ {
		output += sep + input[i]
		sep = realSep
	}

	return "[" + output + "]"
}

func v2(input []string) string {
	var output, sep string
	for _, s := range input {
		output += sep + s
		sep = realSep
	}

	return "[" + output + "]"
}

func v3(input []string) string {
	return "[" + strings.Join(input, " ") + "]"
}

func v4(input []string) string {
	var buf bytes.Buffer
	var sep string

	buf.WriteString("[")

	for _, s := range input {
		buf.WriteString(sep)
		buf.WriteString(s)
		sep = realSep
	}

	buf.WriteString("]")

	return buf.String()
}

func v5(input []string) string {
	return fmt.Sprintf("%v", input)
}

func printResults(results []result) {
	fmt.Printf("Mean concatenation time of %d short strings.\n", inputSize)
	fmt.Printf("(calculated over %d runs with a %0.2f confidence level)\n", runs, confidenceLevel)
	for i, result := range results {
		fmt.Printf("% 17s : %s - %s\n", descriptions[i], result[0], result[1])
	}
}
