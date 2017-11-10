package util

const maxUint = ^uint(0)
const minUint = 0
const maxInt = int(maxUint >> 1)
const minInt = -maxInt - 1

// MinInt returns the minimun
func MinInt(nums ...int) int {
	if len(nums) == 0 {
		return minInt
	}

	min := nums[0]
	for i := 1; i < len(nums); i++ {
		if nums[i] < min {
			min = nums[i]
		}
	}
	return min
}

// MaxInt returns the maximun
func MaxInt(nums ...int) int {
	if len(nums) == 0 {
		return maxInt
	}

	max := nums[0]
	for i := 1; i < len(nums); i++ {
		if nums[i] > max {
			max = nums[i]
		}
	}
	return max
}

// IsIntIn checks if an integer is in a list of integers
func IsIntIn(n int, list ...int) bool {
	for _, i := range list {
		if i == n {
			return true
		}
	}

	return false
}

// IsStringIn checks if a string is in a list of strings
func IsStringIn(s string, list ...string) bool {
	for _, str := range list {
		if str == s {
			return true
		}
	}

	return false
}
