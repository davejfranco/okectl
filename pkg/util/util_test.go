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
