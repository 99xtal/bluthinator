package stats

type number interface {
	float64 | int64
}

func Mean[T number](nums []T) float64 {
	var sum T;
	for _, num := range nums {
		sum += num
	}
	return float64(sum) / float64(len(nums))
}

func Variance[T number](nums []T) float64 {
	mean := Mean(nums)
	var sum float64;

	for _, num := range nums {
		sum += (float64(num) - mean) * (float64(num) - mean)
	}

	return sum / float64(len(nums))
}

func Covariance[T number](nums1, nums2 []T) float64 {
	if len(nums1) != len(nums2) {
		panic("slices must have the same length")
	}

	n := len(nums1)
	if n == 0 {
		panic("slices must not be empty")
	}

	mean1 := Mean(nums1)
	mean2 := Mean(nums2)
	sum := 0.0

	for i, num1 := range nums1 {
		num2 := nums2[i]
		sum += (float64(num1) - mean1) * (float64(num2) - mean2)
	}

	return sum / float64(len(nums1))
}
