package merger

import (
	"reflect"
	"testing"
)

func TestMerger(t *testing.T) {
	nums1 := []int{1, 3, 5}
	nums2 := []int{2, 4, 6, 6, 3}

	result := MergeAndSortUnique(nums1, nums2)
	expected := []int{1, 2, 3, 4, 5, 6}

	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("Expected %v, got %v", expected, result)
	}
}
