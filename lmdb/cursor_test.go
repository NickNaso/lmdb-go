package lmdb

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCursorOpen(t *testing.T) {
	assert := assert.New(t)
	env := openEnv()
	defer closeEnv(env)

	populateDbi(env, TEST_DBI, 1000)

	txn, err := env.BeginTxn(nil, TxnReadOnly)
	if err != nil {
		log.Fatalln(err)
	}
	defer txn.Abort()
	dbi, err := txn.DbiOpen(TEST_DBI, DefaultDbiFlags)
	if err != nil {
		log.Fatalln(err)
	}

	cur, err := txn.CursorOpen(dbi)
	if !assert.NoError(err) {
		return
	}
	defer cur.Close()
	assert.Equal(txn.txn, cur.Txn().txn)
	assert.EqualValues(dbi, cur.Dbi())
}

func TestCursorOps(t *testing.T) {
	assert := assert.New(t)
	env := openEnv()
	defer closeEnv(env)

	populateDbi(env, TEST_DBI, 1000)

	txn, err := env.BeginTxn(nil, TxnReadOnly)
	if err != nil {
		log.Fatalln(err)
	}
	defer txn.Abort()
	dbi, err := txn.DbiOpen(TEST_DBI, DefaultDbiFlags)
	if err != nil {
		log.Fatalln(err)
	}

	cur, err := txn.CursorOpen(dbi)
	if err != nil {
		log.Fatalln(err)
	}
	assert.NoError(txn.CursorRenew(cur))
	defer cur.Close()

	var count int
	k, v, err := cur.Get(nil, OpFirst)
	assert.NoError(err)
	assert.Equal([]byte("key0"), k)
	assert.Equal([]byte("value0"), v)
	count++
	for err == nil {
		k, v, err = cur.Get(k, OpNext)
		if count < 1000 {
			assert.NotEmpty(k)
			assert.NotEmpty(v)
		} else {
			assert.Empty(k)
			assert.Empty(v)
			break
		}
		count++
	}
}

func TestCursorDel(t *testing.T) {
	assert := assert.New(t)
	env := openEnv()
	defer closeEnv(env)

	populateDbi(env, TEST_DBI, 1000)

	txn, err := env.BeginTxn(nil, DefaultTxnFlags)
	if err != nil {
		log.Fatalln(err)
	}
	defer txn.Abort()
	dbi, err := txn.DbiOpen(TEST_DBI, DefaultDbiFlags)
	if err != nil {
		log.Fatalln(err)
	}

	cur, err := txn.CursorOpen(dbi)
	if err != nil {
		log.Fatalln(err)
	}
	defer cur.Close()

	k, _, err := cur.Get(nil, OpFirst)
	if !assert.NoError(err) {
		return
	}
	for i := 0; i < 500; i++ {
		assert.NoError(cur.Del(DefaultWriteFlags))
		k, _, err = cur.Get(k, OpNext)
		assert.NoError(err)
	}

	k, v, err := cur.Get(nil, OpFirst)
	assert.NoError(err)
	assert.Equal([]byte("key549"), k)
	assert.Equal([]byte("value549"), v)
}
