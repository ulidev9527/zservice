package zservice

func MaxInt(nums ...int) int {
	if len(nums) == 0 {
		return 0
	}
	a := nums[0]
	for _, b := range nums {
		if a < b {
			a = b
		}
	}
	return a
}
func MaxInt64(nums ...int64) int64 {
	if len(nums) == 0 {
		return 0
	}
	a := nums[0]
	for _, b := range nums {
		if a < b {
			a = b
		}
	}
	return a
}

func MinInt(nums ...int) int {
	if len(nums) == 0 {
		return 0
	}
	a := nums[0]
	for _, b := range nums {
		if a > b {
			a = b
		}
	}
	return a
}
func MinInt64(nums ...int64) int64 {
	if len(nums) == 0 {
		return 0
	}
	a := nums[0]
	for _, b := range nums {
		if a > b {
			a = b
		}
	}
	return a
}
