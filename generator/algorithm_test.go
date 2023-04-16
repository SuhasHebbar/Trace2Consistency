package generator

import (
	"testing"

	"github.com/SuhasHebbar/CS739-P3/common"
)

func TestConcat(t *testing.T) {
	trace := common.OpTrace{
		common.Operation{
			ClientId:   1,
			SequenceNo: 1,
			Op:         WRITE,
			Key:        "Hello",
			Value:      "World",
		},
		common.Operation{
			ClientId:   1,
			SequenceNo: 2,
			Op:         READ,
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
