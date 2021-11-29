package fuzzymatcher

import "unsafe"

// b2s converts a byte array into a string without allocating new memory
// Note that any changes to a will result in a different string
func b2s(a []byte) string {
	return *(*string)(unsafe.Pointer(&a))
}

// s2b converts a string into a byte array without allocating new memory
// Note that any changes to a will result in a different string
func s2b(a string) []byte {
	return *(*[]byte)(unsafe.Pointer(&a))
}
