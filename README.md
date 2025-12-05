# gocircular

A fast, generic circular buffer (ring buffer) implementation for Go 1.23+.

## Features

- **Generics**: Type-safe buffer for any type
- **O(1) operations**: Push/pop at both ends
- **Iterator support**: Go 1.23 `iter.Seq` and `iter.Seq2` iterators
- **Auto-overwrite**: When full, new elements overwrite oldest ones
- **Zero dependencies**: Standard library only

## Installation

```bash
go get github.com/rvncerr/gocircular
```

## Usage

```go
package main

import (
    "fmt"
    "github.com/rvncerr/gocircular"
)

func main() {
    // Create a buffer with capacity 3
    buf := gocircular.New[int](3)

    // Add elements
    buf.PushBack(1)
    buf.PushBack(2)
    buf.PushBack(3)
    buf.PushBack(4) // Overwrites 1

    // Access elements
    front, _ := buf.Front() // 2
    back, _ := buf.Back()   // 4
    
    // Iterate using Go 1.23 iterators
    for i, v := range buf.All() {
        fmt.Printf("buf[%d] = %d\n", i, v)
    }
    
    // Convert to slice
    slice := buf.ToSlice() // [2, 3, 4]
}
```

## API Reference

### Creation

| Function | Description |
|----------|-------------|
| `New[T](capacity)` | Create a new buffer with given capacity |

### Push/Pop Operations

| Method | Description |
|--------|-------------|
| `PushBack(value)` | Add to back, returns true if element was overwritten |
| `PushFront(value)` | Add to front, returns true if element was overwritten |
| `PopBack()` | Remove and return back element |
| `PopFront()` | Remove and return front element |

### Access Operations

| Method | Description |
|--------|-------------|
| `Front()` | Return front element without removing |
| `Back()` | Return back element without removing |
| `At(index)` | Return element at index (0-based from front) |
| `Set(index, value)` | Update element at index |

### State Operations

| Method | Description |
|--------|-------------|
| `Len()` | Number of elements |
| `Cap()` | Buffer capacity |
| `Empty()` | True if buffer has no elements |
| `Full()` | True if buffer is at capacity |
| `Clear()` | Remove all elements |

### Iteration

| Method | Description |
|--------|-------------|
| `All()` | Iterator yielding (index, value) pairs front-to-back |
| `Values()` | Iterator yielding values front-to-back |
| `Backward()` | Iterator yielding (index, value) pairs back-to-front |
| `Do(func(T) bool)` | Call function on each element, stop if returns false |

### Conversion & Manipulation

| Method | Description |
|--------|-------------|
| `ToSlice()` | Return elements as a slice |
| `Clone()` | Create a deep copy |
| `Resize(newCap)` | Change capacity (discards back elements if shrinking) |

## License

MIT License - see [LICENSE](LICENSE) file.
