package lmdb

import (
	"log"
	"os"
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

const (
	TEST_DB  = "test.db"
	TEST_DB2 = "test_copy.db"
	TEST_DBI = "test_dbi"
)

func initDir(dir string) error {
	return os.MkdirAll(dir, 0755)
}

func nukeDir(dir string) error {
	return os.RemoveAll(dir)
}
