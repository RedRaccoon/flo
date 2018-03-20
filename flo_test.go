package main

import "testing"

func TestReverseToReturnReversedInputString(t *testing.T) {
	actualResult := binaryorder("1")
	var expectedResult = "111"

	if actualResult != expectedResult {
		t.Fatalf("Expected %s but got %s", expectedResult, actualResult)
	}
}
