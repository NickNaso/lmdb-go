package lmdb

import (
	"os"
	"sync"
	"unsafe"

	"github.com/zenhotels/lmdb-go/mdb"
)

type Env struct {
	m   *sync.RWMutex
	env *mdb.Env
}

func EnvCreate() (Env, error) {
	env := Env{
		m: new(sync.RWMutex),
	}
	var mdbEnv *mdb.Env
	if err := mdbError(mdb.EnvCreate(&mdbEnv)); err != nil {
		return env, nil
	}
	env.env = mdbEnv
	return env, nil
}

func (e *Env) Close() {
	e.m.Lock()
	if e.env != nil {
		mdb.EnvClose(e.env)
		e.env = nil
	}
	e.m.Unlock()
}

func (e Env) Open(path string, flags EnvFlags, mode os.FileMode) error {
	e.m.Lock()
	path = path + "\x00" // get a null-terminated C-string
	// WARN: the NoTLS should be always enabled to prevent the driver from using the thread-tied memory.
	err := mdbError(mdb.EnvOpen(e.env, path, uint32(NoTLS|flags), mdb.Mode(mode)))
	e.m.Unlock()
	return err
}

func (e Env) Copy(newpath string) error {
	e.m.RLock()
	newpath = newpath + "\x00" // get a null-terminated C-string
	err := mdbError(mdb.EnvCopy(e.env, newpath))
	e.m.RUnlock()
	return err
}

func (e Env) CopyWithOptions(newpath string, flags CpFlags) error {
	e.m.RLock()
	newpath = newpath + "\x00" // get a null-terminated C-string
	err := mdbError(mdb.EnvCopy2(e.env, newpath, uint32(flags)))
	e.m.RUnlock()
	return err
}

func (e Env) Stat() (Stats, error) {
	e.m.RLock()
	var mdbStats mdb.Stats
	err := mdbError(mdb.EnvStat(e.env, &mdbStats))
	e.m.RUnlock()
	if err != nil {
		return nil, err
	}
	mdbStats.Deref()
	return stats{&mdbStats}, nil
}

func (e Env) Info() (EnvInfo, error) {
	e.m.RLock()
	var mdbEnvInfo mdb.Envinfo
	err := mdbError(mdb.EnvInfo(e.env, &mdbEnvInfo))
	e.m.RUnlock()
	if err != nil {
		return nil, err
	}
	mdbEnvInfo.Deref()
	return envInfo{&mdbEnvInfo}, nil
}

func (e Env) Sync(force bool) {
	e.m.Lock()
	if force {
		mdb.EnvSync(e.env, 1)
	} else {
		mdb.EnvSync(e.env, 0)
	}
	e.m.Unlock()
}

func (e Env) SetFlags(flags EnvFlags, enable bool) (err error) {
	e.m.Lock()
	if enable {
		err = mdbError(mdb.EnvSetFlags(e.env, uint32(flags), 1))
	} else {
		err = mdbError(mdb.EnvSetFlags(e.env, uint32(flags), 0))
	}
	e.m.Unlock()
	return
}

func (e Env) GetFlags() (EnvFlags, error) {
	e.m.RLock()
	var flags uint32
	err := mdbError(mdb.EnvGetFlags(e.env, &flags))
	e.m.RUnlock()
	return EnvFlags(flags), err
}

func (e Env) GetPath() (path string, err error) {
	e.m.RLock()
	paths := make([]string, 1)
	err = mdbError(mdb.EnvGetPath(e.env, paths))
	e.m.RUnlock()
	return paths[0], err
}

func (e Env) SetMapSize(size uint) error {
	e.m.Lock()
	err := mdbError(mdb.EnvSetMapsize(e.env, mdb.Size(size)))
	e.m.Unlock()
	return err
}

func (e Env) SetMaxReaders(readers uint32) error {
	e.m.Lock()
	err := mdbError(mdb.EnvSetMaxreaders(e.env, readers))
	e.m.Unlock()
	return err
}

func (e Env) GetMaxReaders() (uint32, error) {
	e.m.RLock()
	var size uint32
	err := mdbError(mdb.EnvGetMaxreaders(e.env, &size))
	e.m.RUnlock()
	if err != nil {
		return 0, err
	}
	return size, nil
}

func (e Env) SetMaxDBs(dbs uint) error {
	e.m.Lock()
	err := mdbError(mdb.EnvSetMaxdbs(e.env, mdb.Dbi(dbs)))
	e.m.Unlock()
	return err
}

func (e Env) GetMaxKeySize() int32 {
	e.m.RLock()
	size := mdb.EnvGetMaxkeysize(e.env)
	e.m.RUnlock()
	return size
}

func (e Env) SetUserContext(v unsafe.Pointer) error {
	e.m.Lock()
	err := mdbError(mdb.EnvSetUserctx(e.env, v))
	e.m.Unlock()
	return err
}

func (e Env) GetUserContext() unsafe.Pointer {
	e.m.RLock()
	v := mdb.EnvGetUserctx(e.env)
	e.m.RUnlock()
	return v
}

type AssertFunc func(msg string)

func (e Env) SetAssertFunc(assert AssertFunc) error {
	e.m.Lock()
	err := mdbError(mdb.EnvSetAssert(e.env,
		func(env *mdb.Env, msg string) {
			assert(msg)
		}))
	e.m.Unlock()
	return err
}

func (e Env) BeginTxn(parent *Txn, flags TxnFlags) (txn Txn, err error) {
	e.m.Lock()
	if parent != nil {
		err = mdbError(mdb.TxnBegin(e.env, parent.txn, uint32(flags), &txn.txn))
	} else {
		err = mdbError(mdb.TxnBegin(e.env, nil, uint32(flags), &txn.txn))
	}
	e.m.Unlock()
	if err != nil {
		txn.txn = nil
	}
	return
}

func (e Env) DbiClose(dbi Dbi) {
	e.m.Lock()
	mdb.DbiClose(e.env, mdb.Dbi(dbi))
	e.m.Unlock()
}

type MsgFunc func(msg string, ctx unsafe.Pointer) error

func (e Env) ReaderList(msg MsgFunc, ctx unsafe.Pointer) error {
	e.m.RLock()
	msgFunc := func(text string, ctx unsafe.Pointer) int32 {
		if err := msg(text, ctx); err != nil {
			return -1
		}
		return 0
	}
	err := mdbError(mdb.ReaderList(e.env, msgFunc, ctx))
	e.m.RUnlock()
	return err
}

func (e Env) ReaderCheck() (int32, error) {
	e.m.RLock()
	var dead int32
	err := mdbError(mdb.ReaderCheck(e.env, &dead))
	e.m.RUnlock()
	if err != nil {
		return 0, err
	}
	return dead, nil
}
