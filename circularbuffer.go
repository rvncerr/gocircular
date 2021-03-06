package gocontainers

import "errors"

// CircularBuffer is the basic class in gocontainers.
// There are no public members in this struct.
type CircularBuffer struct {
	buffer   []interface{}
	capacity int
	shift    int
	size     int
}

// NewCircularBuffer is the constructor function for CircularBuffer.
func NewCircularBuffer(capacity int) CircularBuffer {
	var cb CircularBuffer

	cb.buffer = make([]interface{}, capacity)
	cb.capacity = capacity
	cb.shift = 0
	cb.size = 0

	return cb
}

// At returns element from CircularBuffer by index.
func (cb *CircularBuffer) At(index int) (interface{}, error) {
	if 0 <= index && index < cb.size {
		return cb.buffer[(cb.shift+index)%cb.capacity], nil
	}
	return nil, errors.New("index out of bounds")
}

// Back returns the back element in CircularBuffer.
// In case of empty CircularBuffer nil returns.
func (cb *CircularBuffer) Back() (interface{}, error) {
	if cb.Empty() {
		return nil, errors.New("empty buffer")
	}
	v, e := cb.At(cb.Size() - 1)
	if e != nil {
		return nil, e
	}
	return v, nil
}

// Capacity returns the maximum possible number elements in CircularBuffer.
func (cb *CircularBuffer) Capacity() int {
	return cb.capacity
}

// Clear removes all the data from CircularBuffer.
func (cb *CircularBuffer) Clear() {
	for i := 0; i < cb.size; i++ {
		cb.buffer[(cb.shift+i)%cb.capacity] = nil
	}
	cb.size = 0
}

// Do calls function f on each element of the CircularBuffer.
func (cb *CircularBuffer) Do(f func(interface{}) error) error {
	for i := 0; i < cb.size; i++ {
		v, e := cb.At(i)
		if e != nil {
			return e
		}
		e = f(v)
		if e != nil {
			return e
		}
	}
	return nil
}

// Empty checks if CircularBuffer has no elements.
func (cb *CircularBuffer) Empty() bool {
	return cb.size == 0
}

// Front returns the front element in CircularBuffer.
// In case of empty CircularBuffer nil returns.
func (cb *CircularBuffer) Front() (interface{}, error) {
	return cb.At(0)
}

// Full checks if CircularBuffer is full.
func (cb *CircularBuffer) Full() bool {
	return cb.size == cb.capacity
}

// PopBack removes back element from CircularBuffer.
func (cb *CircularBuffer) PopBack() {
	if !cb.Empty() {
		cb.buffer[(cb.shift+cb.size-1)%cb.capacity] = nil
		cb.size = cb.size - 1
	}
}

// PopFront removes front element from CircularBuffer.
func (cb *CircularBuffer) PopFront() {
	if !cb.Empty() {
		cb.buffer[cb.shift%cb.capacity] = nil
		cb.size = cb.size - 1
		cb.shift = (cb.shift + 1) % cb.capacity
	}
}

// PushBack appends new element into CircularBuffer.
// If CircularBuffer is full, PopFront() will be called.
func (cb *CircularBuffer) PushBack(value interface{}) {
	if cb.Full() {
		cb.PopFront()
	}
	cb.buffer[(cb.size+cb.shift)%cb.capacity] = value
	cb.size = cb.size + 1
}

// PushFront appends new element into CircularBuffer.
// If CircularBuffer is full, PopBack() will be called.
func (cb *CircularBuffer) PushFront(value interface{}) {
	if cb.Full() {
		cb.PopBack()
	}
	index := (cb.shift + cb.capacity - 1) % cb.capacity
	cb.buffer[index] = value
	cb.shift = index
	cb.size = cb.size + 1
}

// Resize affects capacity of CircularBuffer. TODO: Better algorithm.
func (cb *CircularBuffer) Resize(size int) {
	cb.shiftToZero()
	if size > cb.size {
		if len(cb.buffer) < size {
			abuffer := make([]interface{}, size-len(cb.buffer))
			cb.buffer = append(cb.buffer, abuffer...)
		}
	} else {
		cb.size = size
	}
	cb.capacity = size
}

// shiftToZero makes shift zero. TODO: Make private.
func (cb *CircularBuffer) shiftToZero() {
	var swap = func(i, j int) {
		temp := cb.buffer[i]
		cb.buffer[i] = cb.buffer[j]
		cb.buffer[j] = temp
	}
	var revert = func(i, j int) {
		for k := i; k < (i+j)/2; k++ {
			swap(k, j+i-k-1)
		}
	}
	revert(0, cb.shift)
	revert(cb.shift, cb.capacity)
	revert(0, cb.capacity)
	cb.shift = 0
}

// Size returns number of elements in CircularBuffer.
func (cb *CircularBuffer) Size() int {
	return cb.size
}

// ToArray converts CircularBuffer to Array. TODO: Better algorithm?
func (cb *CircularBuffer) ToArray() []interface{} {
	array := make([]interface{}, cb.size)
	for i := 0; i < cb.size; i++ {
		array[i], _ = cb.At(i)
	}
	return array
}
