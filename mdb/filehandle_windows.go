package mdb

// +build windows

// #include "lmdb.h"
import "C"

// Filehandle type as declared in mdb/lmdb.h:196
type Filehandle C.mdb_filehandle_t
