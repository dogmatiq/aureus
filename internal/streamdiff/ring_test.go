package streamdiff

import "testing"

func TestRing(t *testing.T) {
	r := newRing[int](3)
	check(t, r)

	r.Push(100)
	check(t, r, 100)

	r.Push(200)
	check(t, r, 100, 200)

	r.Push(300)
	check(t, r, 100, 200, 300)

	r.Push(400)
	check(t, r, 200, 300, 400)

	r.Push(500)
	check(t, r, 300, 400, 500)
}

func check[T comparable](t *testing.T, r *ring[T], values ...T) {
	t.Helper()

	if n := r.Len(); n != len(values) {
		t.Fatalf("unexpected length: got %d, want %d", n, len(values))
	}

	for i, v := range values {
		x := r.At(i)
		if x != v {
			t.Fatalf("[at] unexpected value at index %d: got %v, want %v", i, x, v)
		}
	}

	i := 0
	r.Each(
		func(x T) bool {
			if x != values[i] {
				t.Fatalf("[each] unexpected value at index %d: got %v, want %v", i, x, values[i])
			}
			i++
			return true
		},
	)

	for i, x := range r.Slice() {
		if x != values[i] {
			t.Fatalf("[slice] unexpected value at index %d: got %v, want %v", i, x, values[i])
		}
	}
}
