package generator

import (
	"github.com/SuhasHebbar/CS739-P3/common"
	"golang.org/x/exp/slog"
)

type Generator struct {
	serialTrace common.OpTrace
	verifierCh  chan common.OpTrace
	currPerm    []int
}

func NewGenerator(distTrace map[int32]common.OpTrace, verifierCh chan common.OpTrace) *Generator {
	serialTrace := Concat(distTrace)
	firstPerm := Consecutive(len(distTrace))

	return &Generator{
		distTrace:  serialTrace,
		currPerm:   firstPerm,
		verifierCh: verifierCh,
	}
}

func (g *Generator) RunGenerator() {
	slog.Debug("Starting generator ...")

	for {
		// Send the permutation to the verifier
		currTrace := g.getCurrentTrace()
		slog.Debug("Sending ...", "Trace", currTrace)
		g.verifierCh <- currTrace

		// Figure out the next permutation
		if !NextPermutation(g.currPerm) {
			slog.Debug("Exhausted all permutations")
			break
		}
	}
}

func (g *Generator) getCurrentTrace() common.OpTrace {
	currTrace := make(common.OpTrace, len(g.serialTrace))
	for i, idx := range g.ids {
		currTrace[i] = g.serialTrace[idx]
	}
	return currTrace
}
