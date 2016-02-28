package lmdb

import (
	"log"
	"os"
	"path/filepath"
	"syscall"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zenhotels/lmdb-go/mdb"
)

func init() {
	log.SetFlags(log.Lshortfile)
}

func TestVersion(t *testing.T) {
	assert := assert.New(t)
	major, minor, _, version := Version()
	assert.EqualValues(0, major, "only 0.9.x LMDB versions are currently supported")
	assert.EqualValues(9, minor, "only 0.9.x LMDB versions are currently supported")
	assert.NotEmpty(version)
	log.Printf("Testing %s", version)
}

func TestError(t *testing.T) {
	assert := assert.New(t)

	assert.NoError(mdbError(mdb.ErrSuccess))
	assert.Equal("Cannot allocate memory", mdbError(int32(syscall.ENOMEM)).Error())

	assert.Equal("MDB_KEYEXIST: Key/data pair already exists",
		mdbError(mdb.ErrKeyExist).Error())
	assert.Equal("MDB_NOTFOUND: No matching key/data pair found",
		mdbError(mdb.ErrNotFound).Error())
	assert.Equal("MDB_PAGE_NOTFOUND: Requested page not found",
		mdbError(mdb.ErrPageNotFound).Error())
	assert.Equal("MDB_CORRUPTED: Located page was wrong type",
		mdbError(mdb.ErrCorrupted).Error())
	assert.Equal("MDB_PANIC: Update of meta page failed or environment had fatal error",
		mdbError(mdb.ErrPanic).Error())
	assert.Equal("MDB_VERSION_MISMATCH: Database environment version mismatch",
		mdbError(mdb.ErrVersionMismatch).Error())
	assert.Equal("MDB_INVALID: File is not an LMDB file",
		mdbError(mdb.ErrInvalid).Error())
	assert.Equal("MDB_MAP_FULL: Environment mapsize limit reached",
		mdbError(mdb.ErrMapFull).Error())
	assert.Equal("MDB_DBS_FULL: Environment maxdbs limit reached",
		mdbError(mdb.ErrDbsFull).Error())
	assert.Equal("MDB_READERS_FULL: Environment maxreaders limit reached",
		mdbError(mdb.ErrReadersFull).Error())
	assert.Equal("MDB_TLS_FULL: Thread-local storage keys full - too many environments open",
		mdbError(mdb.ErrTlsFull).Error())
	assert.Equal("MDB_TXN_FULL: Transaction has too many dirty pages - transaction too big",
		mdbError(mdb.ErrTxnFull).Error())
	assert.Equal("MDB_CURSOR_FULL: Internal error - cursor stack limit reached",
		mdbError(mdb.ErrCursorFull).Error())
	assert.Equal("MDB_PAGE_FULL: Internal error - page has no more space",
		mdbError(mdb.ErrPageFull).Error())
	assert.Equal("MDB_MAP_RESIZED: Database contents grew beyond environment mapsize",
		mdbError(mdb.ErrMapResized).Error())
	assert.Equal("MDB_INCOMPATIBLE: Operation and DB incompatible, or DB flags changed",
		mdbError(mdb.ErrIncompatible).Error())
	assert.Equal("MDB_BAD_RSLOT: Invalid reuse of reader locktable slot",
		mdbError(mdb.ErrBadRslot).Error())
	assert.Equal("MDB_BAD_TXN: Transaction must abort, has a child, or is invalid",
		mdbError(mdb.ErrBadTxn).Error())
	assert.Equal("MDB_BAD_VALSIZE: Unsupported size of key/DB name/data, or wrong DUPFIXED size",
		mdbError(mdb.ErrValueSize).Error())
	assert.Equal("MDB_BAD_DBI: The specified DBI handle was closed/changed unexpectedly",
		mdbError(mdb.ErrBadDbi).Error())
}

func TestEnvCreateClose(t *testing.T) {
	assert := assert.New(t)
	env, err := EnvCreate()
	assert.NotNil(env.env)
	if assert.NoError(err) {
		env.Close()
	}
	assert.Nil(env.env)
}

func TestEnvOpen(t *testing.T) {
	assert := assert.New(t)
	env, err := EnvCreate()
	if !assert.NoError(err) {
		return
	}
	defer env.Close()

	initDir(TEST_DB)
	defer nukeDir(TEST_DB)
	assert.NoError(env.Open(TEST_DB, DefaultFlags, 0644))
	_, err = os.Stat(filepath.Join(TEST_DB, "data.mdb"))
	assert.NoError(err, "data.mdb should exist in "+TEST_DB)
	_, err = os.Stat(filepath.Join(TEST_DB, "lock.mdb"))
	assert.NoError(err, "lock.mdb should exist in "+TEST_DB)
}

func TestEnvCopy(t *testing.T) {
	assert := assert.New(t)
	env, err := EnvCreate()
	if !assert.NoError(err) {
		return
	}
	defer env.Close()

	initDir(TEST_DB)
	defer nukeDir(TEST_DB)
	if !assert.NoError(env.Open(TEST_DB, DefaultFlags, 0644)) {
		return
	}

	initDir(TEST_DB2)
	defer nukeDir(TEST_DB2)
	assert.NoError(env.Copy(TEST_DB2))
	_, err = os.Stat(filepath.Join(TEST_DB2, "data.mdb"))
	assert.NoError(err, "data.mdb should exist in "+TEST_DB2)
	_, err = os.Stat(filepath.Join(TEST_DB2, "lock.mdb"))
	assert.True(os.IsNotExist(err), "lock.mdb should not exist in "+TEST_DB2)
}

func TestEnvCopyWithOptions(t *testing.T) {
	assert := assert.New(t)
	env, err := EnvCreate()
	if !assert.NoError(err) {
		return
	}
	defer env.Close()

	initDir(TEST_DB)
	defer nukeDir(TEST_DB)
	if !assert.NoError(env.Open(TEST_DB, DefaultFlags, 0644)) {
		return
	}

	initDir(TEST_DB2)
	defer nukeDir(TEST_DB2)
	assert.NoError(env.CopyWithOptions(TEST_DB2, CompactingCopy))
	_, err = os.Stat(filepath.Join(TEST_DB2, "data.mdb"))
	assert.NoError(err, "data.mdb should exist in "+TEST_DB2)
	_, err = os.Stat(filepath.Join(TEST_DB2, "lock.mdb"))
	assert.True(os.IsNotExist(err), "lock.mdb should not exist in "+TEST_DB2)
}

func TestEnvStat(t *testing.T) {
	assert := assert.New(t)
	env, err := EnvCreate()
	if !assert.NoError(err) {
		return
	}
	defer env.Close()

	initDir(TEST_DB)
	defer nukeDir(TEST_DB)
	if !assert.NoError(env.Open(TEST_DB, DefaultFlags, 0644)) {
		return
	}

	stats, err := env.Stats()
	assert.NoError(err)
	assert.EqualValues(stats.PageSize(), 4096)
	// TODO(xlab): bench the stat call
}

const (
	TEST_DB  = "test.db"
	TEST_DB2 = "test_copy.db"
)

func initDir(dir string) error {
	return os.MkdirAll(dir, 0755)
}

func nukeDir(dir string) error {
	return os.RemoveAll(dir)
}
