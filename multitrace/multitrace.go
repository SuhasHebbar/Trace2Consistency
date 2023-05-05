package multitrace

import (
	"cchkr/common"
	"cchkr/generator"
	"fmt"

	"github.com/goware/set"
)

const (
	PermuteChSz = 1024
	ResultChSz  = 10
)

func MultiTraceEntryPoint() {
	// Get trace files from the command line
	config := GetConfig()

	// Figure out the consistencies allowed by the trace files
	validTraces := ExtractTraces(config.ValidTraces)
	faultyTraces := ExtractTraces(config.FaultyTraces)

	// Final output
	// - take the intersection of consistency models output by valid traces
	// - remove all consistency models from the disallowed traces
	consistencies := set.NewStringSet()
	for _, validTrace := range validTraces {
		consistencies = set.Inter(consistencies, GetTraceConsistencies(validTrace))
	}

	for _, faultyTrace := range faultyTraces {
		consistencies = set.Diff(consistencies, GetTraceConsistencies(faultyTrace))
	}

	fmt.Println("Consistencies:")
	fmt.Println("==============")
	for _, consistency := range consistencies {
		fmt.Println(consistency)
	}
}

func ExtractTraces(traceFiles []string) []common.OpTrace {
	traces := make([]common.DistTrace, len(traceFiles))
	for i, traceFile := range traceFiles {
		traces = append(traces, ParseFile(traceFile))
	}
	return traces
}

func GetTraceConsistencies(distTrace common.DistTrace) StringSet {
	permuteCh := make(chan common.OpTrace, PermuteChSz)
	resultCh := make(chan common.VerifierResult, ResultChSz)

	g := generator.NewGenerator(distTrace, permuteCh)
	go g.RunGenerator()

	v := NewVerifier(permuteCh, resultCh)
	go func() {
		v.RunVerifier()
		close(resultCh)
	}()

	consistencies := []string{}
	for consistency := range resultCh {
		consistencies = append(consistencies, consistency)
	}

	return set.NewStringSet(consistencies...)
}
