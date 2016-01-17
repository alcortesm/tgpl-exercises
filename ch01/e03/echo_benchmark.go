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
	// generate input slice of strings
	input := make([]string, inputSize)
	for i := 0; i < inputSize; i++ {
		input[i] = fmt.Sprint(i)
	}

	// run the benchmarks of all the functions in `versions`
	b := make([][2]time.Duration, len(versions))
	var err error
	for i, v := range versions {
		b[i], err = benchmark(v, input)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}

	// print benchmarks results
	fmt.Printf("mean concat execution time (%d elements, %d runs, %0.2f%% confidence level)\n",
		inputSize, runs, confidenceLevel)
	for i, _ := range b {
		fmt.Printf("% 17s : %s - %s\n", descriptions[i], b[i][0], b[i][1])
	}

	os.Exit(0)
}

// Runs `f` function over `input` a number of times (`runs`) and
// calculates the confidence intervals of the mean duration of each run
// for a `confidence` confidence level.
func benchmark(f concat, input []string) ([2]time.Duration, error) {
	var execTimes = make([]float64, runs)
	for i := 0; i < runs; i++ {
		start := time.Now()
		f(input)
		elapsed := time.Since(start)
		execTimes[i] = float64(elapsed)
	}

	durations, err := sample.New(execTimes)
	if err != nil {
		return [2]time.Duration{}, err
	}

	meanConfidenceIntervals, err := durations.MeanConfidenceIntervals(confidenceLevel)
	if err != nil {
		return [2]time.Duration{}, err
	}

	asDuration := [2]time.Duration{
		time.Duration(meanConfidenceIntervals[0]),
		time.Duration(meanConfidenceIntervals[1]),
	}

	return asDuration, nil
}

type concat func([]string) string

var versions []concat = []concat{v1, v2, v3, v4, v5}
var descriptions []string = []string{"+ string operator", "same with range", "strings.Join", "bytes.Buffer", `Sprintf("%q")`}

func v1(input []string) string {
	var output, sep string
	for i := 0; i < len(input); i++ {
		output += sep + input[i]
		sep = " "
	}
	return "[" + output + "]"
}

func v2(input []string) string {
	var output, sep string
	for _, s := range input {
		output += sep + s
		sep = " "
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
		sep = " "
	}
	buf.WriteString("]")
	return buf.String()
}

func v5(input []string) string {
	return fmt.Sprintf("%v", input)
}
