package lmdb

import "github.com/zenhotels/lmdb-go/mdb"

type EnvFlags uint32

const (
	FixedMap    EnvFlags = mdb.Fixedmap
	NoSubdir    EnvFlags = mdb.Nosubdir
	NoSync      EnvFlags = mdb.Nosync
	ReadOnly    EnvFlags = mdb.Rdonly
	NoMetaSync  EnvFlags = mdb.Nometasync
	WriteMap    EnvFlags = mdb.Writemap
	MapAsync    EnvFlags = mdb.Mapasync
	NoTLS       EnvFlags = mdb.Notls
	NoLock      EnvFlags = mdb.Nolock
	NoReadahead EnvFlags = mdb.Nordahead
	NoMemInit   EnvFlags = mdb.Nomeminit

	DefaultFlags EnvFlags = NoTLS
)
