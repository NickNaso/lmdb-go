package lmdb

import "github.com/zenhotels/lmdb-go/mdb"

type Cursor struct {
	cur *mdb.Cursor
}

func (c Cursor) Close() {
	mdb.CursorClose(c.cur)
}

func (c Cursor) Txn() Txn {
	return Txn{
		txn: mdb.CursorTxn(c.cur),
	}
}

func (c Cursor) Dbi() Dbi {
	return Dbi(mdb.CursorDbi(c.cur))
}

func (c Cursor) Get(key []byte, op CursorOp) ([]byte, error) {
	var val mdb.Val
	if err := mdbError(mdb.CursorGet(c.cur, toVal(key), &val, mdb.CursorOp(op))); err != nil {
		return nil, err
	}
	v := fromVal(&val)
	return v, nil
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
