package main

import (
	"fmt"
	"time"
)

const (
	inputSize       = 1000
	runs            = 100000
	confidenceLevel = 0.95
)

func main() {
	input := make([]string, inputSize)
	for i := 0; i < inputSize; i++ {
		input[i] = fmt.Sprint(i)
	}

	fmt.Printf("mean concat execution time (%d elements, %0.2f%% confidence level)\n", inputSize, confidenceLevel)
	fmt.Printf("v1\t%s\n", benchmark(v1, input))
}

// TODO: make this return means according to confidence level
func benchmark(f concat, input []string) time.Duration {
	var total time.Duration
	var start time.Time
	var elapsed time.Duration
	for i := 0; i < runs; i++ {
		start = time.Now()
		f(input)
		elapsed = time.Since(start)
		total += elapsed
	}
	return time.Duration(int64(total) / int64(runs))
}

type concat func([]string) string

func v1(ss []string) (s string) {
	var sep string
	for i := 0; i < len(ss); i++ {
		s += sep + ss[i]
		sep = " "
	}
	return s
}
