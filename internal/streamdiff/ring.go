package streamdiff

type ring[T any] struct {
	values     []T
	pivot, cap int
}

func newRing[T any](n int) *ring[T] {
	return &ring[T]{cap: n}
}

func (r *ring[T]) Push(v T) {
	if len(r.values) < r.cap {
		r.values = append(r.values, v)
	} else {
		r.values[r.pivot] = v
	}

	r.pivot++
	if r.pivot == r.cap {
		r.pivot = 0
	}
}

func (r *ring[T]) Len() int {
	return len(r.values)
}

func (r *ring[T]) At(i int) T {
	i += r.pivot
	if i >= len(r.values) {
		i -= len(r.values)
	}
	return r.values[i]
}

func (r *ring[T]) Each(fn func(T) bool) bool {
	for _, v := range r.values[r.pivot:] {
		if !fn(v) {
			return false
		}
	}
	for _, v := range r.values[:r.pivot] {
		if !fn(v) {
			return false
		}
	}
	return true
}

func (r *ring[T]) Slice() []T {
	slice := make([]T, 0, len(r.values))
	slice = append(slice, r.values[r.pivot:]...)
	slice = append(slice, r.values[:r.pivot]...)
	return slice
}
