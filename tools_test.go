package toolkit

import "testing"

func TestTools_RandomString(t *testing.T) {
	var testTools Tools
	s := testTools.RandomString(10)
	println(len(s), s)
	if len(s) != 10 {
		t.Error("Wrong lenght random string returned")
	}
}
