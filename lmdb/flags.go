package lmdb

import "github.com/zenhotels/lmdb-go/mdb"

type EnvFlags uint32

const (
	FixedMap    EnvFlags = mdb.Fixedmap
	NoSubdir    EnvFlags = mdb.Nosubdir
	NoSync      EnvFlags = mdb.Nosync
	ReadOnly    EnvFlags = mdb.Rdonly
	NoMetaSync  EnvFlags = mdb.Nometasync
	WriteMap    EnvFlags = mdb.Writemap
	MapAsync    EnvFlags = mdb.Mapasync
	NoTLS       EnvFlags = mdb.Notls
	NoLock      EnvFlags = mdb.Nolock
	NoReadahead EnvFlags = mdb.Nordahead
	NoMemInit   EnvFlags = mdb.Nomeminit

	DefaultEnvFlags EnvFlags = NoTLS
)

type CpFlags uint32

const (
	CpCompacting CpFlags = mdb.CpCompact
)

type TxnFlags uint32

const (
	TxnReadOnly   TxnFlags = mdb.Rdonly
	TxnNoSync     TxnFlags = mdb.Nosync
	TxnNoMetaSync TxnFlags = mdb.Nometasync

	DefaultTxnFlags TxnFlags = 0
)

type DbiFlags uint32

const (
	DbiReverseKey DbiFlags = mdb.Reversekey
	DbiDupSort    DbiFlags = mdb.Dupsort
	DbiIntegerKey DbiFlags = mdb.Integerkey
	DbiDupFixed   DbiFlags = mdb.Dupfixed
	DbiIntegerDup DbiFlags = mdb.Integerdup
	DbiReverseDup DbiFlags = mdb.Reversedup
	DbiCreate     DbiFlags = mdb.Create

	DefaultDbiFlags DbiFlags = 0
)

type WriteFlags uint32

const (
	NoOverwrite WriteFlags = mdb.Nooverwrite
	NoDupData   WriteFlags = mdb.Nodupdata
	Current     WriteFlags = mdb.Current
	Reserve     WriteFlags = mdb.Reserve
	Append      WriteFlags = mdb.Append
	AppendDup   WriteFlags = mdb.Appenddup
	Multiple    WriteFlags = mdb.Multiple

	DefaultWriteFlags WriteFlags = 0
)

func (x EnvFlags) Has(flags EnvFlags) bool {
	return x&flags == flags
}

func (x CpFlags) Has(flags CpFlags) bool {
	return x&flags == flags
}

func (x TxnFlags) Has(flags TxnFlags) bool {
	return x&flags == flags
}

func (x DbiFlags) Has(flags DbiFlags) bool {
	return x&flags == flags
}

func (x WriteFlags) Has(flags WriteFlags) bool {
	return x&flags == flags
}

type CursorOp uint32

const (
	OpFirst CursorOp = 0
	OpFirstDup
	OpGetBoth
	OpGetBothRange
	OpGetCurrent
	OpGetMultiple
	OpLast
	OpLastDup
	OpNext
	OpNextDup
	OpNextMultiple
	OpNextNoDup
	OpPrev
	OpPrevDup
	OpPrevNoDup
	OpSet
	OpSetKey
	OpSetRange
	OpPrevMultiple
)
