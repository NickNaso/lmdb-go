package lmdb

import (
	"log"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

func TestEmptyTxn(t *testing.T) {
	assert := assert.New(t)
	env := openEnv()
	defer closeEnv(env)

	txn, err := env.BeginTxn(nil, DefaultTxnFlags)
	assert.NoError(err)
	assert.EqualValues(1, txn.ID())
	assert.NoError(txn.Commit())

	txn, err = env.BeginTxn(nil, DefaultTxnFlags)
	assert.NoError(err)
	assert.EqualValues(1, txn.ID())
	path1, err := env.GetPath()
	assert.NoError(err)
	path2, err := txn.Env().GetPath()
	assert.NoError(err)
	assert.Equal(path1, path2)
	txn.Abort()
}

func TestDbiFlags(t *testing.T) {
	assert := assert.New(t)
	env := openEnv()
	defer closeEnv(env)

	txn, err := env.BeginTxn(nil, DefaultTxnFlags)
	if err != nil {
		log.Fatalln(err)
	}
	dbi, err := txn.DbiOpen(TEST_DBI, DbiCreate|DbiReverseKey)
	assert.NoError(err)
	flags, err := txn.DbiFlags(dbi)
	assert.NoError(err)
	// DbiCreate doesn't get into Dbi flag usually
	assert.EqualValues(DbiReverseKey, flags)
}

func TestDbiDeleteDrop(t *testing.T) {
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
	assert.NoError(err)

	stats, err := txn.Stat(dbi)
	assert.NoError(err)
	assert.EqualValues(1000, stats.Entries())

	assert.NoError(txn.Delete(dbi, []byte("key0")))
	assert.NoError(txn.Delete(dbi, []byte("key1")))
	assert.NoError(txn.Delete(dbi, []byte("key2")))
	assert.Error(txn.Delete(dbi, []byte("key0")))
	stats, err = txn.Stat(dbi)
	assert.NoError(err)
	assert.EqualValues(997, stats.Entries())

	txn.Drop(dbi, false)
	stats, err = txn.Stat(dbi)
	assert.NoError(err)
	assert.EqualValues(0, stats.Entries())

	assert.NoError(txn.Drop(dbi, true))
	_, err = txn.Stat(dbi)
	assert.Error(err)
}

func TestTxnFuncs(t *testing.T) {
	assert := assert.New(t)
	env := openEnv()
	defer closeEnv(env)

	txn, err := env.BeginTxn(nil, DefaultTxnFlags)
	if err != nil {
		log.Fatalln(err)
	}
	dbi, err := txn.DbiOpen(TEST_DBI, DbiCreate|DbiDupSort)
	assert.NoError(err)

	reverseCmp := func(a, b []byte) int {
		// i.e. reverse cmp
		if len(a) < len(b) {
			return 1
		} else if len(a) > len(b) {
			return -1
		}
		return 0
	}

	assert.NoError(txn.SetCompareFunc(dbi, reverseCmp))
	assert.NoError(txn.SetDupsortFunc(dbi, reverseCmp))

	assert.EqualValues(1, txn.Compare(dbi, []byte("ab"), []byte("abc")))
	assert.EqualValues(-1, txn.Compare(dbi, []byte("abc"), []byte("ab")))
	assert.EqualValues(0, txn.Compare(dbi, []byte("ab"), []byte("ab")))
	// this will panic if the Dbi hasn't been opened with DbiDupSort (as per LMDB docs)
	assert.EqualValues(-1, txn.DupCompare(dbi, []byte("abc"), []byte("ab")))

	ctx := 1000
	assert.NoError(txn.SetRelCtx(dbi, unsafe.Pointer(&ctx)))
	assert.NoError(txn.SetRelFunc(dbi, func(item []byte, oldptr, newptr, relctx unsafe.Pointer) {
		log.Println("WARN relocation called for", string(item))
		if assert.NotNil(relctx) {
			log.Println("the context was:", *(*int)(relctx))
		}
	}))
}

func TestTxnDbiPut(t *testing.T) {
	assert := assert.New(t)
	env := openEnv()
	defer closeEnv(env)

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
	env := openEnv()
	defer closeEnv(env)

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

func TestDbiStat(t *testing.T) {
	assert := assert.New(t)
	env := openEnv()
	defer closeEnv(env)

	populateDbi(env, TEST_DBI, 1000)
	tx, err := env.BeginTxn(nil, TxnReadOnly)
	if err != nil {
		log.Fatalln(err)
	}
	defer tx.Abort()

	dbi, err := tx.DbiOpen(TEST_DBI, DefaultDbiFlags)
	if err != nil {
		log.Fatalln(err)
	}

	env.DbiClose(dbi)
	tx.Reset()
	tx.Renew()

	dbi2, err := tx.DbiOpen(TEST_DBI, DefaultDbiFlags)
	if err != nil {
		log.Fatalln(err)
	}
	if !assert.EqualValues(dbi, dbi2) {
		return
	}

	stats, err := tx.Stat(dbi)
	assert.NoError(err)
	assert.EqualValues(2, stats.Depth())
	assert.EqualValues(1, stats.BranchPages())
	assert.EqualValues(1000, stats.Entries())
	assert.EqualValues(11, stats.LeafPages())
	assert.EqualValues(0, stats.OverflowPages())
}
