package domain

import "testing"

func TestNumReg(t *testing.T) {

	s := "High (C3)"

	found := num.Find([]byte(s))

	if string(found) != "3" {
		t.Fatalf("wrong answer. expected: %s, got: %s", "3", string(found))
	}

}
