package lmdb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmptyTxn(t *testing.T) {
	assert := assert.New(t)
	env, err := EnvCreate()
	if !assert.NoError(err) {
		return
	}
	defer env.Close()

	initDir(TEST_DB)
	defer nukeDir(TEST_DB)
	if !assert.NoError(env.Open(TEST_DB, DefaultEnvFlags, 0644)) {
		return
	}

	txn, err := env.BeginTxn(nil, DefaultTxnFlags)
	assert.NoError(err)
	assert.EqualValues(1, txn.ID())
	assert.NoError(txn.Commit())

	txn, err = env.BeginTxn(nil, DefaultTxnFlags)
	assert.NoError(err)
	assert.EqualValues(1, txn.ID())
	txn.Abort()
}

func TestTxnDbiPut(t *testing.T) {
	assert := assert.New(t)
	env, err := EnvCreate()
	if !assert.NoError(err) {
		return
	}
	defer env.Close()
	assert.NoError(env.SetMaxDBs(10))

	initDir(TEST_DB)
	defer nukeDir(TEST_DB)
	if !assert.NoError(env.Open(TEST_DB, DefaultEnvFlags, 0644)) {
		return
	}

	txn, err := env.BeginTxn(nil, DefaultTxnFlags)
	assert.NoError(err)
	assert.EqualValues(1, txn.ID())
	dbi, err := txn.DbiOpen(TEST_DBI, DbiCreate)
	assert.NoError(err)
	assert.NotEmpty(uint32(dbi))
	assert.NoError(txn.Put(dbi, []byte("hello"), []byte("world"), DefaultWriteFlags))
	assert.NoError(txn.Commit())

	txn, err = env.BeginTxn(nil, DefaultTxnFlags)
	assert.NoError(err)
	assert.EqualValues(2, txn.ID())
	txn.Abort()
}

func TestTxnDbiGet(t *testing.T) {
	assert := assert.New(t)
	env, err := EnvCreate()
	if !assert.NoError(err) {
		return
	}
	defer env.Close()
	if !assert.NoError(env.SetMaxDBs(10)) {
		return
	}

	initDir(TEST_DB)
	defer nukeDir(TEST_DB)
	if !assert.NoError(env.Open(TEST_DB, DefaultEnvFlags, 0644)) {
		return
	}

	txn, err := env.BeginTxn(nil, DefaultTxnFlags)
	if !assert.NoError(err) {
		return
	}
	dbi, err := txn.DbiOpen(TEST_DBI, DbiCreate)
	assert.NoError(err)
	key := []byte("hello")
	val := []byte("world")
	assert.NoError(txn.Put(dbi, key, val, DefaultWriteFlags))
	assert.NoError(txn.Commit())

	txn, err = env.BeginTxn(nil, TxnReadOnly)
	if !assert.NoError(err) {
		return
	}
	defer txn.Abort()
	dbi, err = txn.DbiOpen(TEST_DBI, DefaultDbiFlags)
	assert.NoError(err)
	result, err := txn.Get(dbi, key)
	assert.NoError(err)
	if assert.NotNil(result) {
		assert.Equal(val, result)
	}
}
