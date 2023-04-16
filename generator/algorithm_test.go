package generator

import (
	"testing"

	// "github.com/SuhasHebbar/CS739-P3/common"
	"cchkr/common"
)

func TestConcat(t *testing.T) {
	trace := common.OpTrace{
		common.Operation{
			ClientId:   1,
			SequenceNo: 1,
			Op:         common.WRITE,
			Key:        "Hello",
			Value:      "World",
		},
		common.Operation{
			ClientId:   1,
			SequenceNo: 2,
			Op:         common.READ,
			Key:        "Hello",
			Value:      "World",
		},
	}

	cted := Concat(map[int32]common.OpTrace{
		1: trace,
	})

	if len(cted) != 2 {
		t.Fatalf(`Incorrect length. Expected %v, Got %v`, 2, len(cted))
	}

	for i := range cted {
		if cted[i] != trace[i] {
			t.Fatalf(`Output differs at position %v`, i)
		}
	}
}

func TestConsecutive(t *testing.T) {
	n := 13
	con := Consecutive(n)

	if len(con) != n {
		t.Fatalf(`Incorrect length. Expected %v, Got %v`, n, len(con))
	}

	for i, val := range con {
		if i != val {
			t.Fatalf(`Incorrect value at idx %v`, i)
		}
	}
}

func TestIsSortedUntil(t *testing.T) {
	// empty
	{
		input := []int{}
		output := 0
		su := IsSortedUntil(input, RevAccess)
		if su != output {
			t.Fatalf(`Expected %v, Got %v`, output, su)
		}
	}

	// non-increasing
	{
		input := []int{3, 3, 2}
		output := len(input)
		su := IsSortedUntil(input, RevAccess)
		if su != output {
			t.Fatalf(`Expected %v, Got %v`, output, su)
		}
	}

	// peak case
	{
		input := []int{2, 3, 3, 2}
		output := 3
		su := IsSortedUntil(input, RevAccess)
		if su != output {
			t.Fatalf(`Expected %v, Got %v`, output, su)
		}
	}
}
