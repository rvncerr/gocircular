package gocircular

import (
	"slices"
	"testing"
)

func TestNew(t *testing.T) {
	b := New[int](5)
	if b.Len() != 0 {
		t.Errorf("new buffer Len() = %d, want 0", b.Len())
	}
	if b.Cap() != 5 {
		t.Errorf("new buffer Cap() = %d, want 5", b.Cap())
	}
	if !b.Empty() {
		t.Error("new buffer should be empty")
	}
}

func TestNew_Panic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("New(0) should panic")
		}
	}()
	New[int](0)
}

func TestPushBack(t *testing.T) {
	b := New[int](3)

	if overwritten := b.PushBack(1); overwritten {
		t.Error("PushBack on non-full buffer should not overwrite")
	}
	if overwritten := b.PushBack(2); overwritten {
		t.Error("PushBack on non-full buffer should not overwrite")
	}
	if overwritten := b.PushBack(3); overwritten {
		t.Error("PushBack on non-full buffer should not overwrite")
	}

	want := []int{1, 2, 3}
	if got := b.ToSlice(); !slices.Equal(got, want) {
		t.Errorf("ToSlice() = %v, want %v", got, want)
	}

	// Overflow
	if overwritten := b.PushBack(4); !overwritten {
		t.Error("PushBack on full buffer should overwrite")
	}

	want = []int{2, 3, 4}
	if got := b.ToSlice(); !slices.Equal(got, want) {
		t.Errorf("after overflow ToSlice() = %v, want %v", got, want)
	}
}

func TestPushFront(t *testing.T) {
	b := New[int](3)

	if overwritten := b.PushFront(1); overwritten {
		t.Error("PushFront on non-full buffer should not overwrite")
	}
	if overwritten := b.PushFront(2); overwritten {
		t.Error("PushFront on non-full buffer should not overwrite")
	}
	if overwritten := b.PushFront(3); overwritten {
		t.Error("PushFront on non-full buffer should not overwrite")
	}

	want := []int{3, 2, 1}
	if got := b.ToSlice(); !slices.Equal(got, want) {
		t.Errorf("ToSlice() = %v, want %v", got, want)
	}

	// Overflow
	if overwritten := b.PushFront(4); !overwritten {
		t.Error("PushFront on full buffer should overwrite")
	}

	want = []int{4, 3, 2}
	if got := b.ToSlice(); !slices.Equal(got, want) {
		t.Errorf("after overflow ToSlice() = %v, want %v", got, want)
	}
}

func TestPopBack(t *testing.T) {
	b := New[int](3)

	// Pop from empty
	if _, ok := b.PopBack(); ok {
		t.Error("PopBack on empty buffer should return false")
	}

	b.PushBack(1)
	b.PushBack(2)
	b.PushBack(3)

	v, ok := b.PopBack()
	if !ok || v != 3 {
		t.Errorf("PopBack() = %d, %v; want 3, true", v, ok)
	}

	v, ok = b.PopBack()
	if !ok || v != 2 {
		t.Errorf("PopBack() = %d, %v; want 2, true", v, ok)
	}

	if b.Len() != 1 {
		t.Errorf("Len() = %d, want 1", b.Len())
	}
}

func TestPopFront(t *testing.T) {
	b := New[int](3)

	// Pop from empty
	if _, ok := b.PopFront(); ok {
		t.Error("PopFront on empty buffer should return false")
	}

	b.PushBack(1)
	b.PushBack(2)
	b.PushBack(3)

	v, ok := b.PopFront()
	if !ok || v != 1 {
		t.Errorf("PopFront() = %d, %v; want 1, true", v, ok)
	}

	v, ok = b.PopFront()
	if !ok || v != 2 {
		t.Errorf("PopFront() = %d, %v; want 2, true", v, ok)
	}

	if b.Len() != 1 {
		t.Errorf("Len() = %d, want 1", b.Len())
	}
}

func TestFront(t *testing.T) {
	b := New[string](2)

	if _, ok := b.Front(); ok {
		t.Error("Front on empty buffer should return false")
	}

	b.PushBack("hello")
	v, ok := b.Front()
	if !ok || v != "hello" {
		t.Errorf("Front() = %q, %v; want 'hello', true", v, ok)
	}

	// Front shouldn't remove the element
	if b.Len() != 1 {
		t.Error("Front should not remove element")
	}
}

func TestBack(t *testing.T) {
	b := New[string](2)

	if _, ok := b.Back(); ok {
		t.Error("Back on empty buffer should return false")
	}

	b.PushBack("hello")
	b.PushBack("world")

	v, ok := b.Back()
	if !ok || v != "world" {
		t.Errorf("Back() = %q, %v; want 'world', true", v, ok)
	}

	// Back shouldn't remove the element
	if b.Len() != 2 {
		t.Error("Back should not remove element")
	}
}

func TestAt(t *testing.T) {
	b := New[int](4)
	b.PushBack(10)
	b.PushBack(20)
	b.PushBack(30)

	tests := []struct {
		index int
		want  int
		ok    bool
	}{
		{0, 10, true},
		{1, 20, true},
		{2, 30, true},
		{-1, 0, false},
		{3, 0, false},
	}

	for _, tt := range tests {
		v, ok := b.At(tt.index)
		if ok != tt.ok || (ok && v != tt.want) {
			t.Errorf("At(%d) = %d, %v; want %d, %v", tt.index, v, ok, tt.want, tt.ok)
		}
	}
}

func TestSet(t *testing.T) {
	b := New[int](3)
	b.PushBack(1)
	b.PushBack(2)
	b.PushBack(3)

	if !b.Set(1, 42) {
		t.Error("Set(1, 42) should return true")
	}

	v, _ := b.At(1)
	if v != 42 {
		t.Errorf("At(1) = %d, want 42", v)
	}

	if b.Set(-1, 0) {
		t.Error("Set(-1, 0) should return false")
	}
	if b.Set(3, 0) {
		t.Error("Set(3, 0) should return false")
	}
}

func TestEmptyFull(t *testing.T) {
	b := New[int](2)

	if !b.Empty() {
		t.Error("new buffer should be empty")
	}
	if b.Full() {
		t.Error("new buffer should not be full")
	}

	b.PushBack(1)
	if b.Empty() {
		t.Error("buffer with element should not be empty")
	}
	if b.Full() {
		t.Error("buffer with 1/2 elements should not be full")
	}

	b.PushBack(2)
	if b.Empty() {
		t.Error("full buffer should not be empty")
	}
	if !b.Full() {
		t.Error("buffer at capacity should be full")
	}
}

func TestClear(t *testing.T) {
	b := New[int](3)
	b.PushBack(1)
	b.PushBack(2)
	b.PushBack(3)

	b.Clear()

	if !b.Empty() {
		t.Error("buffer should be empty after Clear")
	}
	if b.Len() != 0 {
		t.Errorf("Len() = %d, want 0 after Clear", b.Len())
	}
	if b.Cap() != 3 {
		t.Errorf("Cap() = %d, want 3 after Clear", b.Cap())
	}
}

func TestToSlice(t *testing.T) {
	b := New[int](5)
	b.PushBack(1)
	b.PushBack(2)
	b.PushBack(3)

	got := b.ToSlice()
	want := []int{1, 2, 3}
	if !slices.Equal(got, want) {
		t.Errorf("ToSlice() = %v, want %v", got, want)
	}

	// Test after wrap-around
	b.PushBack(4)
	b.PushBack(5)
	b.PushBack(6) // Overwrites 1
	b.PushBack(7) // Overwrites 2

	got = b.ToSlice()
	want = []int{3, 4, 5, 6, 7}
	if !slices.Equal(got, want) {
		t.Errorf("after wraparound ToSlice() = %v, want %v", got, want)
	}
}

func TestAll(t *testing.T) {
	b := New[int](3)
	b.PushBack(10)
	b.PushBack(20)
	b.PushBack(30)

	var indices []int
	var values []int
	for i, v := range b.All() {
		indices = append(indices, i)
		values = append(values, v)
	}

	wantIndices := []int{0, 1, 2}
	wantValues := []int{10, 20, 30}
	if !slices.Equal(indices, wantIndices) {
		t.Errorf("indices = %v, want %v", indices, wantIndices)
	}
	if !slices.Equal(values, wantValues) {
		t.Errorf("values = %v, want %v", values, wantValues)
	}
}

func TestValues(t *testing.T) {
	b := New[int](3)
	b.PushBack(10)
	b.PushBack(20)
	b.PushBack(30)

	var values []int
	for v := range b.Values() {
		values = append(values, v)
	}

	want := []int{10, 20, 30}
	if !slices.Equal(values, want) {
		t.Errorf("values = %v, want %v", values, want)
	}
}

func TestBackward(t *testing.T) {
	b := New[int](3)
	b.PushBack(10)
	b.PushBack(20)
	b.PushBack(30)

	var indices []int
	var values []int
	for i, v := range b.Backward() {
		indices = append(indices, i)
		values = append(values, v)
	}

	wantIndices := []int{2, 1, 0}
	wantValues := []int{30, 20, 10}
	if !slices.Equal(indices, wantIndices) {
		t.Errorf("indices = %v, want %v", indices, wantIndices)
	}
	if !slices.Equal(values, wantValues) {
		t.Errorf("values = %v, want %v", values, wantValues)
	}
}

func TestResize(t *testing.T) {
	b := New[int](4)
	b.PushBack(1)
	b.PushBack(2)
	b.PushBack(3)
	b.PushBack(4)

	// Resize larger
	b.Resize(6)
	if b.Cap() != 6 {
		t.Errorf("Cap() = %d, want 6", b.Cap())
	}
	want := []int{1, 2, 3, 4}
	if got := b.ToSlice(); !slices.Equal(got, want) {
		t.Errorf("ToSlice() = %v, want %v", got, want)
	}

	// Resize smaller - should discard back elements
	b.Resize(2)
	if b.Cap() != 2 {
		t.Errorf("Cap() = %d, want 2", b.Cap())
	}
	want = []int{1, 2}
	if got := b.ToSlice(); !slices.Equal(got, want) {
		t.Errorf("ToSlice() = %v, want %v", got, want)
	}
}

func TestResize_Panic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Resize(0) should panic")
		}
	}()
	b := New[int](3)
	b.Resize(0)
}

func TestClone(t *testing.T) {
	b := New[int](4)
	b.PushBack(1)
	b.PushBack(2)
	b.PushBack(3)

	clone := b.Clone()

	// Verify clone has same data
	if !slices.Equal(clone.ToSlice(), b.ToSlice()) {
		t.Error("clone should have same elements")
	}
	if clone.Cap() != b.Cap() {
		t.Error("clone should have same capacity")
	}

	// Verify independence
	b.PushBack(4)
	if clone.Len() != 3 {
		t.Error("clone should be independent from original")
	}
}

func TestDo(t *testing.T) {
	b := New[int](4)
	b.PushBack(1)
	b.PushBack(2)
	b.PushBack(3)
	b.PushBack(4)

	var sum int
	b.Do(func(v int) bool {
		sum += v
		return true
	})

	if sum != 10 {
		t.Errorf("sum = %d, want 10", sum)
	}

	// Test early termination
	var collected []int
	b.Do(func(v int) bool {
		collected = append(collected, v)
		return v != 2 // stop at 2
	})

	want := []int{1, 2}
	if !slices.Equal(collected, want) {
		t.Errorf("collected = %v, want %v", collected, want)
	}
}

func TestWrapAround(t *testing.T) {
	b := New[int](4)

	// Fill and overflow multiple times
	for i := 0; i < 10; i++ {
		b.PushBack(i)
	}

	want := []int{6, 7, 8, 9}
	if got := b.ToSlice(); !slices.Equal(got, want) {
		t.Errorf("ToSlice() = %v, want %v", got, want)
	}

	// Pop some and push more
	b.PopFront()
	b.PopFront()
	b.PushBack(10)
	b.PushBack(11)

	want = []int{8, 9, 10, 11}
	if got := b.ToSlice(); !slices.Equal(got, want) {
		t.Errorf("ToSlice() = %v, want %v", got, want)
	}
}

func TestMixedOperations(t *testing.T) {
	b := New[int](4)

	b.PushBack(1)
	b.PushBack(2)
	b.PushFront(0)
	// [0, 1, 2]

	want := []int{0, 1, 2}
	if got := b.ToSlice(); !slices.Equal(got, want) {
		t.Errorf("ToSlice() = %v, want %v", got, want)
	}

	b.PopBack()
	// [0, 1]

	b.PushFront(-1)
	// [-1, 0, 1]

	want = []int{-1, 0, 1}
	if got := b.ToSlice(); !slices.Equal(got, want) {
		t.Errorf("ToSlice() = %v, want %v", got, want)
	}
}

// Benchmarks

func BenchmarkPushBack(b *testing.B) {
	buf := New[int](1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf.PushBack(i)
	}
}

func BenchmarkPushFront(b *testing.B) {
	buf := New[int](1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf.PushFront(i)
	}
}

func BenchmarkPopFront(b *testing.B) {
	buf := New[int](1000)
	for i := 0; i < 1000; i++ {
		buf.PushBack(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if buf.Empty() {
			for j := 0; j < 1000; j++ {
				buf.PushBack(j)
			}
		}
		buf.PopFront()
	}
}

func BenchmarkPopBack(b *testing.B) {
	buf := New[int](1000)
	for i := 0; i < 1000; i++ {
		buf.PushBack(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if buf.Empty() {
			for j := 0; j < 1000; j++ {
				buf.PushBack(j)
			}
		}
		buf.PopBack()
	}
}

func BenchmarkAt(b *testing.B) {
	buf := New[int](1000)
	for i := 0; i < 1000; i++ {
		buf.PushBack(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf.At(i % 1000)
	}
}
