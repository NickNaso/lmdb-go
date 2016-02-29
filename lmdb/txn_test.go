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
