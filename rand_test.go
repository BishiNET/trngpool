package truerand

import "testing"

func TestRand(t *testing.T) {
	r := New("")
	r.Refresh()

	t.Log(r.Get())
	t.Log(r.Get())
	t.Log(r.Get())
	t.Log(r.Slices())
}
