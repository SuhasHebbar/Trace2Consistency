package generator

import (
	"sort"

	"github.com/SuhasHebbar/CS739-P3/common"
)

func Concat(distTrace map[int32]common.OpTrace) common.OpTrace {
	serialTrace := make(common.OpTrace)
	for _, trace := range distTrace {
		serialTrace = append(serialTrace, trace...)
	}
	return serialTrace
}

func Consecutive(sz int) []int {
	firstPerm := make([]int, sz)
	for i := 0; i < traceLen; i++ {
		firstPerm[i] = i
	}
	return firstPerm
}

// Adapted from https://en.cppreference.com/w/cpp/algorithm/next_permutation
func NextPermutation(nums []int) bool {
	first := 0
	last := len(nums)
	left := IsSortedUntil(nums, RevAccess)

	if left != last {
		right := UpperBound(RevSplice(nums, first, left), RevAccess(left), RevAccess)
		RevSwap(nums, left, right)
	}

	sort.Reverse(sort.IntSlice(nums[len(nums)-left : len(nums)-first]))
}

// Returns the index which is strictly less than the next element from the end
// Returns -1 for non-increasing numbers
// Adapted from https://en.cppreference.com/w/cpp/algorithm/is_sorted_until
func IsSortedUntil(nums []int, access func([]int, int) int) int {
	first := 0
	last := len(nums)

	// slice empty check
	if first != last {
		next := first
		for {
			// get the next element
			next++

			// check if we reached the end
			if next == last {
				break
			}

			// the check
			// we need strict inequality since it is sorted otherwise
			if Comp(access(nums, next), access(nums, first)) {
				return next
			}

			first = next
		}
	}

	return last
}

// Gets the first index strictly greater than the target
// N when current element is the largest
// Adapted from https://en.cppreference.com/w/cpp/algorithm/upper_bound
func UpperBound(nums []int, target int, access func([]int, int) int) int {
	first := 0
	last := len(nums)
	count := last - first
	var step, curr int

	for count > 0 {
		curr = first
		step = count / 2
		curr += step

		if !Comp(target, access(nums, curr)) {
			curr++
			first = curr
			count -= step + 1
		} else {
			count = step
		}
	}

	return first
}

func RevAccess(nums []int, int idx) int {
	return nums[len(nums)-1-idx]
}

func RevSplice(nums []int, first int, last int) []int {
	n := len(nums)
	return nums[n-first : n-last]
}

func RevSwap(nums []int, first int, second int) {
	n := len(nums)
	rFirst := n - first
	rSecond := n - second
	nums[rFirst], nums[rSecond] = nums[rSecond], nums[rFirst]
}

func Comp(first int, second int) bool {
	return first < second
}
