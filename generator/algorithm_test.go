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

func TestUpperBound(t *testing.T) {
	// empty
	{
		input := []int{}
		output := 0
		su := UpperBound(input, 0, RevAccess)
		if su != output {
			t.Fatalf(`Expected %v, Got %v`, output, su)
		}
	}

	// largest element
	{
		input := []int{4, 3, 3, 2}
		output := 4
		su := UpperBound(input, 4, RevAccess)
		if su != output {
			t.Fatalf(`Expected %v, Got %v`, output, su)
		}
	}

	// skip equal
	{
		input := []int{4, 3, 3, 2}
		output := 3
		su := UpperBound(input, 3, RevAccess)
		if su != output {
			t.Fatalf(`Expected %v, Got %v`, output, su)
		}
	}

	// smaller than smallest
	{
		input := []int{4, 3, 3, 2}
		output := 0
		su := UpperBound(input, 1, RevAccess)
		if su != output {
			t.Fatalf(`Expected %v, Got %v`, output, su)
		}
	}
}

func TestReverse(t *testing.T) {
	// empty
	{
		input := []int{}
		Reverse(input)
		if len(input) != 0 {
			t.Fatalf(`Expected %v, Got %v`, []int{}, input)
		}
	}

	// increasing odd
	{
		input := []int{1, 2, 3}
		output := []int{3, 2, 1}
		Reverse(input)
		if len(input) != len(output) {
			t.Fatalf(`Expected %v, Got %v`, output, input)
		}

		for i, val := range output {
			if input[i] != val {
				t.Fatalf(`Expected %v, Got %v`, output, input)
			}
		}
	}

	// decerasing even
	{
		input := []int{4, 3, 2, 1}
		output := []int{1, 2, 3, 4}
		Reverse(input)
		if len(input) != len(output) {
			t.Fatalf(`Expected %v, Got %v`, output, input)
		}

		for i, val := range output {
			if input[i] != val {
				t.Fatalf(`Expected %v, Got %v`, output, input)
			}
		}
	}
}

func TestNextPermutation(t *testing.T) {
	// permutations of 1,2,3 are
	perms := [][]int{
		[]int{0, 1, 2},
		[]int{0, 2, 1},
		[]int{1, 0, 2},
		[]int{1, 2, 0},
		[]int{2, 0, 1},
		[]int{2, 1, 0},
	}

	// start with the initial permutation
	perm := Consecutive(3)

	// run other iterations
	for i := 0; i <= len(perms); i++ {
		j := i % len(perms)

		// check that the array slices are equal
		actual := perm
		expected := perms[j]
		if len(actual) != len(expected) {
			t.Fatalf(`Incorrect length for iter %v. Expected %v, Got %v`, i, expected, actual)
		}

		for k, val := range expected {
			if actual[k] != val {
				t.Fatalf(`Incorrect value for iter %v. Expected %v, Got %v`, i, expected, actual)
			}
		}

		if !(i == len(perms) || NextPermutation(perm) == (i < len(perms)-1)) {
			t.Fatalf(`Expected true return value`)
		}
	}
}

func CheckIntSliceEqual(expected []int, actual []int, t *testing.T) {
}
