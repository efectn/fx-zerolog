package fxzerolog

import "testing"

func Test_AssertEqual(t *testing.T) {
	t.Parallel()
	assertEqual(nil, []string{}, []string{})
	assertEqual(t, []string{}, []string{})
}
