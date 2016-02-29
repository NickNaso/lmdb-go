package lmdb

import (
	"unsafe"

	"github.com/zenhotels/lmdb-go/mdb"
)

type Stats interface {
	PageSize() uint32
	Depth() uint32
	BranchPages() uint
	LeafPages() uint
	OverflowPages() uint
	Entries() uint
}

type stats struct {
	*mdb.Stats
}

func (s stats) PageSize() uint32 {
	return s.Stats.Psize
}

func (s stats) Depth() uint32 {
	return s.Stats.Depth
}

func (s stats) BranchPages() uint {
	return uint(s.Stats.BranchPages)
}

func (s stats) LeafPages() uint {
	return uint(s.Stats.LeafPages)
}

func (s stats) OverflowPages() uint {
	return uint(s.Stats.OverflowPages)
}

func (s stats) Entries() uint {
	return uint(s.Stats.Entries)
}

type EnvInfo interface {
	MapAddr() unsafe.Pointer
	MapSize() uint
	LastPageNo() uint
	LastTxnID() uint
	MaxReaders() uint32
	NumReaders() uint32
}

type envInfo struct {
	*mdb.Envinfo
}

func (e envInfo) MapAddr() unsafe.Pointer {
	return e.Envinfo.Mapaddr
}
func (e envInfo) MapSize() uint {
	return uint(e.Envinfo.Mapsize)
}
func (e envInfo) LastPageNo() uint {
	return uint(e.Envinfo.LastPgno)
}
func (e envInfo) LastTxnID() uint {
	return uint(e.Envinfo.LastTxnid)
}
func (e envInfo) MaxReaders() uint32 {
	return e.Envinfo.Maxreaders
}
func (e envInfo) NumReaders() uint32 {
	return e.Envinfo.Numreaders
}
