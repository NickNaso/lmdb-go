package lmdb

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
	assert.NoError(env.Open(TEST_DB, DefaultEnvFlags, 0644))
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
	if !assert.NoError(env.Open(TEST_DB, DefaultEnvFlags, 0644)) {
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
	if !assert.NoError(env.Open(TEST_DB, DefaultEnvFlags, 0644)) {
		return
	}

	initDir(TEST_DB2)
	defer nukeDir(TEST_DB2)
	assert.NoError(env.CopyWithOptions(TEST_DB2, CpCompacting))
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
	if !assert.NoError(env.Open(TEST_DB, DefaultEnvFlags, 0644)) {
		return
	}

	stat, err := env.Stat()
	// TODO(xlab): bench the stat call

	assert.NoError(err)
	assert.EqualValues(0x1000, stat.PageSize())
	assert.EqualValues(0, stat.Depth())
	assert.EqualValues(0, stat.BranchPages())
	assert.EqualValues(0, stat.Entries())
	assert.EqualValues(0, stat.LeafPages())
	assert.EqualValues(0, stat.OverflowPages())
}

func TestEnvInfo(t *testing.T) {
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

	info, err := env.Info()
	// TODO(xlab): bench the info call

	assert.NoError(err)
	assert.True(info.MapAddr() == nil)
	assert.EqualValues(0x100000, info.MapSize())
	assert.EqualValues(1, info.LastPageNo())
	assert.EqualValues(0, info.LastTxnID())
	assert.EqualValues(126, info.MaxReaders())
	assert.EqualValues(0, info.NumReaders())
}

func TestEnvFlagsSetGet(t *testing.T) {
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

	flags, err := env.GetFlags()
	assert.NoError(err)
	assert.True(flags.Has(NoTLS))
	assert.False(flags.Has(NoMetaSync))
	assert.False(flags.Has(NoSync))

	assert.NoError(env.SetFlags(NoMetaSync|NoSync, true))
	flags, err = env.GetFlags()
	assert.NoError(err)
	assert.True(flags.Has(NoMetaSync | NoSync | NoTLS))
	assert.True(flags.Has(NoMetaSync | NoSync))
	assert.True(flags.Has(NoMetaSync | NoTLS))
	assert.True(flags.Has(NoMetaSync))
	assert.True(flags.Has(NoSync | NoTLS))
	assert.True(flags.Has(NoTLS))

	assert.NoError(env.SetFlags(NoSync, false))
	flags, err = env.GetFlags()
	assert.NoError(err)
	assert.False(flags.Has(NoMetaSync | NoSync | NoTLS))
	assert.False(flags.Has(NoMetaSync | NoSync))
	assert.True(flags.Has(NoMetaSync | NoTLS))
	assert.True(flags.Has(NoMetaSync))
	assert.False(flags.Has(NoSync))
	assert.True(flags.Has(NoTLS))
}

func TestEnvMisc(t *testing.T) {
	assert := assert.New(t)
	env, err := EnvCreate()
	if !assert.NoError(err) {
		return
	}
	defer env.Close()

	path, err := env.GetPath()
	assert.NoError(err)
	assert.Empty(path)

	assert.NoError(env.SetMapSize(0x200000))
	// should be called before Open only
	assert.NoError(env.SetMaxReaders(256))
	// should be called before Open only
	assert.NoError(env.SetMaxDBs(100))

	initDir(TEST_DB)
	defer nukeDir(TEST_DB)
	if !assert.NoError(env.Open(TEST_DB, DefaultEnvFlags, 0644)) {
		return
	}

	path, err = env.GetPath()
	assert.NoError(err)
	assert.Equal(TEST_DB, path)

	info, err := env.Info()
	assert.NoError(err)
	assert.True(info.MapAddr() == nil)
	assert.EqualValues(0x200000, info.MapSize())
	assert.EqualValues(256, info.MaxReaders())
	assert.EqualValues(511, env.GetMaxKeySize())
}
