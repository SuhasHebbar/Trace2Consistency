package generator

import (
	"cchkr/common"
	"testing"
)

func TestCP(t *testing.T) {
	// Client 1
	w13 := common.Operation{
		ClientId:   1,
		SequenceNo: 0,
		Op:         common.WRITE,
		Key:        "Key",
		Value:      "3",
	}
	r12 := common.Operation{
		ClientId:   1,
		SequenceNo: 1,
		Op:         common.READ,
		Key:        "Key",
		Value:      "2",
	}
	r13 := common.Operation{
		ClientId:   1,
		SequenceNo: 2,
		Op:         common.WRITE,
		Key:        "Key",
		Value:      "3",
	}
	c1 := common.OpTrace{
		w13,
		r12,
		r13,
	}

	// Client 2
	w22 := common.Operation{
		ClientId:   2,
		SequenceNo: 0,
		Op:         common.WRITE,
		Key:        "Key",
		Value:      "2",
	}
	r23 := common.Operation{
		ClientId:   2,
		SequenceNo: 1,
		Op:         common.READ,
		Key:        "Key",
		Value:      "3",
	}
	r22 := common.Operation{
		ClientId:   2,
		SequenceNo: 2,
		Op:         common.WRITE,
		Key:        "Key",
		Value:      "2",
	}
	c2 := common.OpTrace{
		w22,
		r23,
		r22,
	}

	distTrace := map[int32]common.OpTrace{
		1: c1,
		2: c2,
	}
	ch := make(chan common.OpTrace, 1000)
	g := NewGenerator(distTrace, ch)
	go g.RunGenerator()

	foundExec1 := false
	exec1 := common.OpTrace{
		w22,
		r12,
		r22,
		w13,
		r23,
		r13,
	}

	foundExec2 := false
	exec2 := common.OpTrace{
		w13,
		r13,
		r23,
		w22,
		r22,
		r12,
	}

	count := 0
	for trace := range ch {
		count++

		if len(trace) != len(exec1) {
			t.Fatalf(`Unexpected length. Trace: %v`, trace)
		}

		if checkSliceEq(trace, exec1) {
			foundExec1 = true
		}

		if checkSliceEq(trace, exec2) {
			foundExec2 = true
		}
	}

	sixFactorial := 720
	if count != sixFactorial {
		t.Fatalf(`Expected %v num permutations but found %v permutations`, sixFactorial, count)
	}

	if !(foundExec1 && foundExec2) {
		t.Fatalf(`Couldn't find required permutations!`)
	}
}

func checkSliceEq(s1, s2 []common.Operation) bool {
	for i, val := range s2 {
		if s1[i] != val {
			return false
		}
	}

	return true
}
