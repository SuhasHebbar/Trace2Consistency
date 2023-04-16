package generator

import (
	// "github.com/SuhasHebbar/CS739-P3/common"
	"cchkr/common"
)

func Concat(distTrace map[int32]common.OpTrace) common.OpTrace {
	serialTrace := common.OpTrace{}
	for _, trace := range distTrace {
		serialTrace = append(serialTrace, trace...)
	}
	return serialTrace
}

func Consecutive(sz int) []int {
	firstPerm := make([]int, sz)
	for i := 0; i < sz; i++ {
		firstPerm[i] = i
	}
	return firstPerm
}

// Adapted from https://en.cppreference.com/w/cpp/algorithm/next_permutation
func NextPermutation(nums []int) bool {
	last := len(nums)
	left := IsSortedUntil(nums, RevAccess)

	if left != last {
		right := UpperBound(nums[len(nums)-left:], RevAccess(nums, left), RevAccess)
		RevSwap(nums, left, right)
	}

	Reverse(nums[len(nums)-left:])

	return left != last
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
// N when target is at least largest element
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

func RevAccess(nums []int, idx int) int {
	return nums[len(nums)-1-idx]
}

func RevSwap(nums []int, first int, second int) {
	n := len(nums)
	rFirst := n - 1 - first
	rSecond := n - 1 - second
	nums[rFirst], nums[rSecond] = nums[rSecond], nums[rFirst]
}

// Adapted from https://en.cppreference.com/w/cpp/algorithm/reverse
func Reverse(nums []int) {
	first := 0
	last := len(nums)

	// empty
	if first == last {
		return
	}

	last--
	// swap the first and last pointer coming closer to the middle
	for first < last {
		nums[first], nums[last] = nums[last], nums[first]
		first++
		last--
	}
}

func Comp(first int, second int) bool {
	return first < second
}
