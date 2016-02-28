package lmdb

import (
	"errors"
	"os"
	"sync"

	"github.com/zenhotels/lmdb-go/mdb"
)

func Version() (major, minor, patch int32, version string) {
	v := mdb.Version(&major, &minor, &patch)
	version = toString(v, 64)
	return
}

func mdbError(err int32) error {
	switch err {
	case mdb.ErrSuccess:
		return nil
	case mdb.ErrKeyExist:
		return ErrKeyExist
	case mdb.ErrNotFound:
		return ErrNotFound
	case mdb.ErrPageNotFound:
		return ErrPageNotFound
	case mdb.ErrCorrupted:
		return ErrCorrupted
	case mdb.ErrPanic:
		return ErrPanic
	case mdb.ErrVersionMismatch:
		return ErrVersionMismatch
	case mdb.ErrInvalid:
		return ErrInvalid
	case mdb.ErrMapFull:
		return ErrMapFull
	case mdb.ErrDbsFull:
		return ErrDbsFull
	case mdb.ErrReadersFull:
		return ErrReadersFull
	case mdb.ErrTlsFull:
		return ErrTlsFull
	case mdb.ErrTxnFull:
		return ErrTxnFull
	case mdb.ErrCursorFull:
		return ErrCursorFull
	case mdb.ErrPageFull:
		return ErrPageFull
	case mdb.ErrMapResized:
		return ErrMapResized
	case mdb.ErrIncompatible:
		return ErrIncompatible
	case mdb.ErrBadRslot:
		return ErrBadRslot
	case mdb.ErrBadTxn:
		return ErrBadTxn
	case mdb.ErrValueSize:
		return ErrValueSize
	case mdb.ErrBadDbi:
		return ErrBadDbi
	default:
		return errors.New(strError(err))
	}
}

func strError(err int32) string {
	return toString(mdb.StrError(err), 255)
}

var (
	ErrSuccess         = errors.New(strError(mdb.ErrSuccess))
	ErrKeyExist        = errors.New(strError(mdb.ErrKeyExist))
	ErrNotFound        = errors.New(strError(mdb.ErrNotFound))
	ErrPageNotFound    = errors.New(strError(mdb.ErrPageNotFound))
	ErrCorrupted       = errors.New(strError(mdb.ErrCorrupted))
	ErrPanic           = errors.New(strError(mdb.ErrPanic))
	ErrVersionMismatch = errors.New(strError(mdb.ErrVersionMismatch))
	ErrInvalid         = errors.New(strError(mdb.ErrInvalid))
	ErrMapFull         = errors.New(strError(mdb.ErrMapFull))
	ErrDbsFull         = errors.New(strError(mdb.ErrDbsFull))
	ErrReadersFull     = errors.New(strError(mdb.ErrReadersFull))
	ErrTlsFull         = errors.New(strError(mdb.ErrTlsFull))
	ErrTxnFull         = errors.New(strError(mdb.ErrTxnFull))
	ErrCursorFull      = errors.New(strError(mdb.ErrCursorFull))
	ErrPageFull        = errors.New(strError(mdb.ErrPageFull))
	ErrMapResized      = errors.New(strError(mdb.ErrMapResized))
	ErrIncompatible    = errors.New(strError(mdb.ErrIncompatible))
	ErrBadRslot        = errors.New(strError(mdb.ErrBadRslot))
	ErrBadTxn          = errors.New(strError(mdb.ErrBadTxn))
	ErrValueSize       = errors.New(strError(mdb.ErrValueSize))
	ErrBadDbi          = errors.New(strError(mdb.ErrBadDbi))
)

type Env struct {
	m   *sync.RWMutex
	env *mdb.Env
}

func EnvCreate() (*Env, error) {
	var env Env
	if err := mdbError(mdb.EnvCreate(&env.env)); err != nil {
		return nil, err
	}
	env.m = new(sync.RWMutex)
	return &env, nil
}

func (e *Env) Close() {
	e.m.Lock()
	if e.env != nil {
		mdb.EnvClose(e.env)
		e.env = nil
	}
	e.m.Unlock()
}

func (e *Env) Open(path string, flags EnvFlags, mode os.FileMode) error {
	e.m.Lock()
	// WARN: the NoTLS should be always enabled to prevent the driver from using the thread-tied memory.
	err := mdbError(mdb.EnvOpen(e.env, path, uint32(NoTLS|flags), mdb.Mode(mode)))
	e.m.Unlock()
	return err
}

func (e *Env) Copy(newpath string) error {
	e.m.RLock()
	err := mdbError(mdb.EnvCopy(e.env, newpath))
	e.m.RUnlock()
	return err
}

func (e *Env) CopyWithOptions(newpath string, flags CopyFlags) error {
	e.m.RLock()
	err := mdbError(mdb.EnvCopy2(e.env, newpath, uint32(flags)))
	e.m.RUnlock()
	return err
}

func (e *Env) Stats() (Stats, error) {
	e.m.RLock()
	var mdbStats mdb.Stats
	if err := mdbError(mdb.EnvStat(e.env, &mdbStats)); err != nil {
		return nil, err
	}
	mdbStats.Deref()
	e.m.RUnlock()
	return stats{&mdbStats}, nil
}
