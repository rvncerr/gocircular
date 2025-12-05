// Package gocircular provides a generic circular buffer implementation.
//
// A circular buffer (also known as a ring buffer) is a fixed-size data structure
// that uses a single, continuous buffer as if it were connected end-to-end.
// When the buffer is full, new elements overwrite the oldest elements.
package gocircular

import "iter"

// Buffer is a generic circular buffer with a fixed capacity.
// It supports efficient O(1) operations at both ends.
type Buffer[T any] struct {
	buf []T // pre-allocated to capacity; len(buf) is the capacity
	r   int // read position (front index)
	n   int // number of elements
}

// New creates a new circular buffer with the given capacity.
// Panics if capacity is less than 1.
func New[T any](capacity int) *Buffer[T] {
	if capacity < 1 {
		panic("gocircular: capacity must be at least 1")
	}
	return &Buffer[T]{
		buf: make([]T, capacity),
	}
}

// PushBack adds an element to the back of the buffer.
// If the buffer is full, the front element is overwritten.
// Returns true if an element was overwritten.
func (b *Buffer[T]) PushBack(v T) bool {
	c := len(b.buf)
	full := b.n == c
	if full {
		b.r = (b.r + 1) % c
	} else {
		b.n++
	}
	b.buf[(b.r+b.n-1)%c] = v
	return full
}

// PushFront adds an element to the front of the buffer.
// If the buffer is full, the back element is overwritten.
// Returns true if an element was overwritten.
func (b *Buffer[T]) PushFront(v T) bool {
	c := len(b.buf)
	full := b.n == c
	b.r = (b.r + c - 1) % c
	if !full {
		b.n++
	}
	b.buf[b.r] = v
	return full
}

// PopBack removes and returns the back element.
// Returns the zero value and false if the buffer is empty.
func (b *Buffer[T]) PopBack() (T, bool) {
	var zero T
	if b.n == 0 {
		return zero, false
	}
	c := len(b.buf)
	i := (b.r + b.n - 1) % c
	v := b.buf[i]
	b.buf[i] = zero // clear for GC
	b.n--
	return v, true
}

// PopFront removes and returns the front element.
// Returns the zero value and false if the buffer is empty.
func (b *Buffer[T]) PopFront() (T, bool) {
	var zero T
	if b.n == 0 {
		return zero, false
	}
	c := len(b.buf)
	v := b.buf[b.r]
	b.buf[b.r] = zero // clear for GC
	b.r = (b.r + 1) % c
	b.n--
	return v, true
}

// Front returns the front element without removing it.
// Returns the zero value and false if the buffer is empty.
func (b *Buffer[T]) Front() (T, bool) {
	var zero T
	if b.n == 0 {
		return zero, false
	}
	return b.buf[b.r], true
}

// Back returns the back element without removing it.
// Returns the zero value and false if the buffer is empty.
func (b *Buffer[T]) Back() (T, bool) {
	var zero T
	if b.n == 0 {
		return zero, false
	}
	return b.buf[(b.r+b.n-1)%len(b.buf)], true
}

// At returns the element at the given index (0-based from front).
// Returns the zero value and false if the index is out of bounds.
func (b *Buffer[T]) At(i int) (T, bool) {
	var zero T
	if i < 0 || i >= b.n {
		return zero, false
	}
	return b.buf[(b.r+i)%len(b.buf)], true
}

// Set updates the element at the given index (0-based from front).
// Returns false if the index is out of bounds.
func (b *Buffer[T]) Set(i int, v T) bool {
	if i < 0 || i >= b.n {
		return false
	}
	b.buf[(b.r+i)%len(b.buf)] = v
	return true
}

// Len returns the number of elements in the buffer.
func (b *Buffer[T]) Len() int {
	return b.n
}

// Cap returns the capacity of the buffer.
func (b *Buffer[T]) Cap() int {
	return len(b.buf)
}

// Empty returns true if the buffer has no elements.
func (b *Buffer[T]) Empty() bool {
	return b.n == 0
}

// Full returns true if the buffer is at capacity.
func (b *Buffer[T]) Full() bool {
	return b.n == len(b.buf)
}

// Clear removes all elements from the buffer.
func (b *Buffer[T]) Clear() {
	var zero T
	c := len(b.buf)
	for i := 0; i < b.n; i++ {
		b.buf[(b.r+i)%c] = zero
	}
	b.r = 0
	b.n = 0
}

// ToSlice returns a new slice containing all elements in order (front to back).
func (b *Buffer[T]) ToSlice() []T {
	s := make([]T, b.n)
	c := len(b.buf)
	for i := 0; i < b.n; i++ {
		s[i] = b.buf[(b.r+i)%c]
	}
	return s
}

// All returns an iterator over all elements from front to back.
// The iterator yields index and value pairs.
func (b *Buffer[T]) All() iter.Seq2[int, T] {
	return func(yield func(int, T) bool) {
		c := len(b.buf)
		for i := 0; i < b.n; i++ {
			if !yield(i, b.buf[(b.r+i)%c]) {
				return
			}
		}
	}
}

// Values returns an iterator over all values from front to back.
func (b *Buffer[T]) Values() iter.Seq[T] {
	return func(yield func(T) bool) {
		c := len(b.buf)
		for i := 0; i < b.n; i++ {
			if !yield(b.buf[(b.r+i)%c]) {
				return
			}
		}
	}
}

// Backward returns an iterator over all elements from back to front.
// The iterator yields index and value pairs (indices are still front-based).
func (b *Buffer[T]) Backward() iter.Seq2[int, T] {
	return func(yield func(int, T) bool) {
		c := len(b.buf)
		for i := b.n - 1; i >= 0; i-- {
			if !yield(i, b.buf[(b.r+i)%c]) {
				return
			}
		}
	}
}

// Resize changes the buffer capacity.
// If shrinking, elements at the back are discarded to fit the new capacity.
// Panics if n is less than 1.
func (b *Buffer[T]) Resize(n int) {
	if n < 1 {
		panic("gocircular: capacity must be at least 1")
	}
	c := len(b.buf)
	if n == c {
		return
	}
	// When shrinking, discard excess elements from the back
	if n < b.n {
		b.n = n
	}
	b.resize(n)
}

// resize reallocates the internal buffer with the new capacity.
func (b *Buffer[T]) resize(n int) {
	nb := make([]T, n)
	c := len(b.buf)
	for i := 0; i < b.n; i++ {
		nb[i] = b.buf[(b.r+i)%c]
	}
	b.buf = nb
	b.r = 0
}

// Clone creates a deep copy of the buffer.
func (b *Buffer[T]) Clone() *Buffer[T] {
	c := len(b.buf)
	b2 := &Buffer[T]{
		buf: make([]T, c),
		r:   0,
		n:   b.n,
	}
	for i := 0; i < b.n; i++ {
		b2.buf[i] = b.buf[(b.r+i)%c]
	}
	return b2
}

// Do calls the function f on each element from front to back.
// If f returns false, iteration stops early.
func (b *Buffer[T]) Do(f func(T) bool) {
	c := len(b.buf)
	for i := 0; i < b.n; i++ {
		if !f(b.buf[(b.r+i)%c]) {
			return
		}
	}
}
