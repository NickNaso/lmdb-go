package lmdb

import (
	"sync"
	"unsafe"

	"github.com/zenhotels/lmdb-go/mdb"
)

type Dbi uint32

type Txn struct {
	txn *mdb.Txn
}

func (tx Txn) Env() Env {
	return Env{
		env: mdb.TxnEnv(tx.txn),
		m:   new(sync.RWMutex),
	}
}

func (tx Txn) ID() uint {
	return uint(mdb.TxnId(tx.txn))
}

func (tx Txn) Commit() error {
	return mdbError(mdb.TxnCommit(tx.txn))
}

func (tx Txn) Abort() {
	mdb.TxnAbort(tx.txn)
}

func (tx Txn) Reset() {
	mdb.TxnReset(tx.txn)
}

func (tx Txn) Renew() error {
	return mdbError(mdb.TxnRenew(tx.txn))
}

func (tx Txn) DbiOpen(name string, flags DbiFlags) (Dbi, error) {
	var mdbDbi mdb.Dbi
	name = name + "\x00" // get a null-terminated C-string
	if err := mdbError(mdb.DbiOpen(tx.txn, name, uint32(flags), &mdbDbi)); err != nil {
		return 0, err
	}
	return Dbi(mdbDbi), nil
}

func (tx Txn) Stat(dbi Dbi) (Stats, error) {
	var mdbStats mdb.Stats
	if err := mdbError(mdb.Stat(tx.txn, mdb.Dbi(dbi), &mdbStats)); err != nil {
		return nil, err
	}
	mdbStats.Deref()
	return stats{&mdbStats}, nil
}

func (tx Txn) DbiFlags(dbi Dbi) (DbiFlags, error) {
	var flags uint32
	err := mdbError(mdb.DbiFlags(tx.txn, mdb.Dbi(dbi), &flags))
	return DbiFlags(flags), err
}

func (tx Txn) Drop(dbi Dbi, del bool) error {
	if del {
		return mdbError(mdb.Drop(tx.txn, mdb.Dbi(dbi), 1))
	}
	return mdbError(mdb.Drop(tx.txn, mdb.Dbi(dbi), 0))
}

type CmpFunc func(a, b []byte) int

func (tx Txn) SetCompareFunc(dbi Dbi, cmp CmpFunc) error {
	return mdbError(mdb.SetCompare(tx.txn, mdb.Dbi(dbi),
		func(a, b *mdb.Val) int32 {
			a.Deref()
			b.Deref()
			return int32(cmp(fromVal(a), fromVal(b)))
		}))
}

func (tx Txn) SetDupsortFunc(dbi Dbi, cmp CmpFunc) error {
	return mdbError(mdb.SetDupsort(tx.txn, mdb.Dbi(dbi),
		func(a, b *mdb.Val) int32 {
			a.Deref()
			b.Deref()
			return int32(cmp(fromVal(a), fromVal(b)))
		}))
}

type RelFunc func(item []byte, oldptr, newptr, relctx unsafe.Pointer)

func (tx Txn) SetRelFunc(dbi Dbi, rel RelFunc) error {
	return mdbError(mdb.SetRelfunc(tx.txn, mdb.Dbi(dbi),
		func(item *mdb.Val, oldptr unsafe.Pointer, newptr unsafe.Pointer, relctx unsafe.Pointer) {
			item.Deref()
			rel(fromVal(item), oldptr, newptr, relctx)
		}))
}

func (tx Txn) SetRelCtx(dbi Dbi, ctx unsafe.Pointer) error {
	return mdbError(mdb.SetRelctx(tx.txn, mdb.Dbi(dbi), ctx))
}

func (tx Txn) Get(dbi Dbi, key []byte) ([]byte, error) {
	var val mdb.Val
	if err := mdbError(mdb.Get(tx.txn, mdb.Dbi(dbi), toVal(key), &val)); err != nil {
		return nil, err
	}
	val.Deref()
	v := fromVal(&val)
	return v, nil
}

func (tx Txn) Put(dbi Dbi, key, value []byte, flags WriteFlags) error {
	return mdbError(mdb.Put(tx.txn, mdb.Dbi(dbi), toVal(key), toVal(value), uint32(flags)))
}

func (tx Txn) Delete(dbi Dbi, key []byte) error {
	return mdbError(mdb.Del(tx.txn, mdb.Dbi(dbi), toVal(key), nil))
}

func (tx Txn) DeleteDup(dbi Dbi, key, data []byte) error {
	return mdbError(mdb.Del(tx.txn, mdb.Dbi(dbi), toVal(key), toVal(data)))
}

func (tx Txn) Compare(dbi Dbi, a, b []byte) int {
	return int(mdb.Cmp(tx.txn, mdb.Dbi(dbi), toVal(a), toVal(b)))
}

func (tx Txn) DupCompare(dbi Dbi, a, b []byte) int {
	return int(mdb.Dcmp(tx.txn, mdb.Dbi(dbi), toVal(a), toVal(b)))
}

func (tx Txn) CursorOpen(dbi Dbi) (Cursor, error) {
	var mdbCursor *mdb.Cursor
	if err := mdbError(mdb.CursorOpen(tx.txn, mdb.Dbi(dbi), &mdbCursor)); err != nil {
		return Cursor{}, err
	}
	cur := Cursor{
		cur: mdbCursor,
	}
	return cur, nil
}

func (tx Txn) CursorRenew(cursor Cursor) error {
	return mdbError(mdb.CursorRenew(tx.txn, cursor.cur))
}
