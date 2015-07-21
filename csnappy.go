// Package csnappy are cgo bindings for the snappy compression library.
/*
This package exposes the same interface as `code.google.com/p/snappy-go`, so
using it in your code should be as easy as changing the import line
appropriately.  Note that this packages name is `csnappy`, so you'll need to
import it with the `snappy` identifier if you want true drop-in compatibility.

*/
package csnappy

/*
#cgo LDFLAGS: -lsnappy

#include <snappy-c.h>
*/
import "C"

import (
	"strconv"
	"unsafe"
)

// Errno is a snappy error
type Errno int

var errText = map[Errno]string{
	0: "ok",
	1: "invalid input",
	2: "buffer too small",
}

func (e Errno) Error() string {
	s := errText[e]
	if s == "" {
		return "snappy errno " + strconv.Itoa(int(e))
	}
	return s
}

// Errors returned from the package
var (
	ErrOk             = Errno(0)
	ErrInvalidInput   = Errno(1)
	ErrBufferTooSmall = Errno(2)
)

// Encode compresses the byte array src [and returns the compressed data.  The
// returned array may be a subslice of dst if it was large enough.
func Encode(dst, src []byte) []byte {

	if src == nil {
		return nil
	}

	dLen := C.snappy_max_compressed_length(C.size_t(len(src)))

	if C.size_t(len(dst)) < dLen {
		dst = make([]byte, dLen)
	}

	C.snappy_compress((*C.char)(unsafe.Pointer(&src[0])), C.size_t(len(src)),
		(*C.char)(unsafe.Pointer(&dst[0])), (*C.size_t)(unsafe.Pointer(&dLen)))

	return dst[:dLen]
}

// Decode uncompresses the byte array src and returns the uncompressed data.
// The returned slice may be a sub-slice of dst if it was large enough.
func Decode(dst, src []byte) ([]byte, error) {

	if src == nil {
		return nil, nil
	}

	dLen, err := DecodedLen(src)
	if err != nil {
		return nil, err
	}

	if len(dst) < dLen {
		dst = make([]byte, dLen)
	}

	cerr := C.snappy_uncompress((*C.char)(unsafe.Pointer(&src[0])), C.size_t(len(src)),
		(*C.char)(unsafe.Pointer(&dst[0])), (*C.size_t)(unsafe.Pointer(&dLen)))

	// decompression failed :(
	if cerr != 0 {
		return nil, Errno(cerr)
	}

	return dst[:dLen], nil
}

// MaxEncodedLen returns the maximum length of a snappy block, given its
// uncompressed length.
func MaxEncodedLen(srcLen int) int {
	return int(C.snappy_max_compressed_length(C.size_t(srcLen)))
}

// DecodedLen returns the length of the decoded block.
func DecodedLen(src []byte) (int, error) {

	result := 0

	err := C.snappy_uncompressed_length((*C.char)(unsafe.Pointer(&src[0])), C.size_t(len(src)), (*C.size_t)(unsafe.Pointer(&result)))

	if err != C.SNAPPY_OK {
		return 0, Errno(err)
	}

	return result, nil
}
