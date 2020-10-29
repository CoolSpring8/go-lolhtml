package lolhtml

// Credit to https://github.com/mattn/go-pointer.

// #include <stdlib.h>
import "C"
import (
	"sync"
	"unsafe"
)

// sync.Map documentation states that it is optimized for "when the entry for a given key is only
// ever written once but read many times, as in caches that only grow". My benchmarks show that sync.Map
// version rewriter is slower in single-goroutine calls, but faster when used in multiple goroutines
// (and personally I think the latter is more important).
var store sync.Map

func savePointer(v interface{}) unsafe.Pointer {
	if v == nil {
		return nil
	}

	ptr := C.malloc(C.size_t(1))
	if ptr == nil {
		panic(`can't allocate "cgo-pointer hack index pointer": ptr == nil`)
	}

	store.Store(ptr, v)

	return ptr
}

func restorePointer(ptr unsafe.Pointer) (v interface{}) {
	if ptr == nil {
		return nil
	}

	if v, ok := store.Load(ptr); ok {
		return v
	}
	return nil
}

func unrefPointer(ptr unsafe.Pointer) {
	if ptr == nil {
		return
	}

	store.Delete(ptr)

	C.free(ptr)
}

func unrefPointers(ptrs []unsafe.Pointer) {
	for _, ptr := range ptrs {
		unrefPointer(ptr)
	}
}
