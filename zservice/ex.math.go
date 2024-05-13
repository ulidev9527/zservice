package zservice

func MaxInt(a int, b int, nums ...int) int {

	if a < b {
		a = b
	}

	len := len(nums)
	if len == 0 {
		return a
	}

	b = nums[0]
	return MaxInt(a, b, append(nums[:1], nums[1:]...)...)
}

func MinInt(a int, b int, nums ...int) int {

	if a > b {
		a = b
	}

	len := len(nums)
	if len == 0 {
		return a
	}

	b = nums[0]
	return MinInt(a, b, append(nums[:1], nums[1:]...)...)
}
