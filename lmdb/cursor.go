package lmdb

import (
	"errors"

	"github.com/zenhotels/lmdb-go/mdb"
)

type Cursor struct {
	cur *mdb.Cursor

	envClosed <-chan struct{}
}

func (c Cursor) Close() {
	mdb.CursorClose(c.cur)
	c.cur = nil
}

func (c Cursor) Txn() Txn {
	return Txn{
		txn: mdb.CursorTxn(c.cur),
	}
}

func (c Cursor) Dbi() Dbi {
	return Dbi(mdb.CursorDbi(c.cur))
}

func (c Cursor) Get(key []byte, op CursorOp) (newkey, value []byte, err error) {
	kval := toVal(key)
	var vval mdb.Val
	select {
	case <-c.envClosed:
		err = errors.New("lmdb: env closed")
		return
	default:
		if err = mdbError(mdb.CursorGet(c.cur, kval, &vval, mdb.CursorOp(op))); err != nil {
			return
		}
	}
	vval.Deref()
	kval.Deref()
	value = fromVal(&vval)
	newkey = fromVal(kval)
	return
}

func (c Cursor) Put(key, value []byte, flags WriteFlags) error {
	return mdbError(mdb.CursorPut(c.cur, toVal(key), toVal(value), uint32(flags)))
}

func (c Cursor) Del(flags WriteFlags) error {
	return mdbError(mdb.CursorDel(c.cur, uint32(flags)))
}

func (c Cursor) Count() (uint, error) {
	var count mdb.Size
	if err := mdbError(mdb.CursorCount(c.cur, &count)); err != nil {
		return 0, err
	}
	return uint(count), nil
}
