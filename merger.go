package main

import "sort"

func MergeAndSortUnique(nums1, nums2 []int) []int {
	nums := append(nums1, nums2...)
	sort.Ints(nums)

	// remove duplicates
	uniqueNums := []int{}
	seen := map[int]bool{}

	for _, num := range nums {
		if !seen[num] {
			seen[num] = true
			uniqueNums = append(uniqueNums, num)
		}
	}

	return uniqueNums
}
