package lmdb

import (
	"unsafe"

	"github.com/zenhotels/lmdb-go/mdb"
)

// TODO(xlab): bench toVal
func toVal(v []byte) *mdb.Val {
	if v == nil {
		return &mdb.Val{}
	}
	hdr := (*sliceHeader)(unsafe.Pointer(&v))
	return &mdb.Val{
		Data: unsafe.Pointer(hdr.Data),
		Size: uint(hdr.Len),
	}
}

// TODO(xlab): bench fromVal
func fromVal(v *mdb.Val) []byte {
	if v == nil {
		return nil
	}
	return *(*[]byte)(unsafe.Pointer(&sliceHeader{
		Data: uintptr(v.Data),
		Len:  int(v.Size),
		Cap:  int(v.Size),
	}))
}

type sliceHeader struct {
	Data uintptr
	Len  int
	Cap  int
}
