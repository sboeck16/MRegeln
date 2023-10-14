package main

import (
	"testing"
)

func TestUtils(t *testing.T) {
	arr1 := []string{"A", "B"}
	arr2 := []string{"B", "A"}
	if !sameArr(arr1, arr2) {
		t.Error("array checking went wrong")
	}
}
