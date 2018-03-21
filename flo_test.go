package main

import "testing"

func TestBinaryOrderSimple(t *testing.T) {
	res := binaryorder("1")
	var expected = "11"

	if res != expected {
		t.Fatalf("Expected %s but got %s", res, expected)
	}

}

func TestBinaryOrderMaxSize(t *testing.T) {
	tsFileMaxSize := getTsFileMaxSize()
	res := "1"
	for i := 0; i < tsFileMaxSize; i++ {
		res += "1"
	}

	res = binaryorder(res)
	expected := "2"

	if res != expected {
		t.Fatalf("Expected %s but got %s", res, expected)
	}

}
