package util

import "testing"

// Write test for FatalIfError function

func TestRandomString(t *testing.T) {
	// Write test for RandomString function
	random := RandomString(10)
	if len(random) != 10 {
		t.Errorf("RandomString function failed")
	}
}

func TestCreateZipFile(t *testing.T) {
	// Write test for CreateZipFile function
	err := CreateZipFile("stack.zip")
	if err != nil {
		t.Errorf("CreateZipFile function failed")
	}
}
