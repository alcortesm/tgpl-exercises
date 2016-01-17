package main

import (
	"fmt"
	"os"
	"time"

	"github.com/alcortesm/sample"
)

const (
	inputSize       = 1000
	runs            = 1000
	confidenceLevel = 0.95
)

func main() {
	input := make([]string, inputSize)
	for i := 0; i < inputSize; i++ {
		input[i] = fmt.Sprint(i)
	}

	durationConfidenceIntervals, err := benchmark(v1, input)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Printf("mean concat execution time (%d elements, %d runs, %0.2f%% confidence level)\n",
		inputSize, runs, confidenceLevel)
	fmt.Printf("v1\t[ %s - %s ]\n", durationConfidenceIntervals[0], durationConfidenceIntervals[1])
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

func v1(input []string) (output string) {
	var sep string
	for i := 0; i < len(input); i++ {
		output += sep + input[i]
		sep = " "
	}
	return output
}
