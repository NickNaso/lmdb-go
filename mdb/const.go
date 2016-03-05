// THE AUTOGENERATED LICENSE. ALL THE RIGHTS ARE RESERVED BY ROBOTS.

// WARNING: This file has automatically been generated on Sat, 05 Mar 2016 21:42:22 MSK.
// By http://git.io/cgogen. DO NOT EDIT.

package mdb

/*
#include "lmdb.h"
#include <stdlib.h>
#include "cgo_helpers.h"
*/
import "C"

const (
	// VersionMajor as defined in mdb/lmdb.h:207
	VersionMajor = 0
	// VersionMinor as defined in mdb/lmdb.h:209
	VersionMinor = 9
	// VersionPatch as defined in mdb/lmdb.h:211
	VersionPatch = 70
	// VersionFull as defined in mdb/lmdb.h:217
	VersionFull = 589894
	// VersionDate as defined in mdb/lmdb.h:221
	VersionDate = "December 19, 2015"
	// Fixedmap as defined in mdb/lmdb.h:293
	Fixedmap = 0x01
	// Nosubdir as defined in mdb/lmdb.h:295
	Nosubdir = 0x4000
	// Nosync as defined in mdb/lmdb.h:297
	Nosync = 0x10000
	// Rdonly as defined in mdb/lmdb.h:299
	Rdonly = 0x20000
	// Nometasync as defined in mdb/lmdb.h:301
	Nometasync = 0x40000
	// Writemap as defined in mdb/lmdb.h:303
	Writemap = 0x80000
	// Mapasync as defined in mdb/lmdb.h:305
	Mapasync = 0x100000
	// Notls as defined in mdb/lmdb.h:307
	Notls = 0x200000
	// Nolock as defined in mdb/lmdb.h:309
	Nolock = 0x400000
	// Nordahead as defined in mdb/lmdb.h:311
	Nordahead = 0x800000
	// Nomeminit as defined in mdb/lmdb.h:313
	Nomeminit = 0x1000000
	// Reversekey as defined in mdb/lmdb.h:320
	Reversekey = 0x02
	// Dupsort as defined in mdb/lmdb.h:322
	Dupsort = 0x04
	// Integerkey as defined in mdb/lmdb.h:325
	Integerkey = 0x08
	// Dupfixed as defined in mdb/lmdb.h:327
	Dupfixed = 0x10
	// Integerdup as defined in mdb/lmdb.h:329
	Integerdup = 0x20
	// Reversedup as defined in mdb/lmdb.h:331
	Reversedup = 0x40
	// Create as defined in mdb/lmdb.h:333
	Create = 0x40000
	// Nooverwrite as defined in mdb/lmdb.h:340
	Nooverwrite = 0x10
	// Nodupdata as defined in mdb/lmdb.h:345
	Nodupdata = 0x20
	// Current as defined in mdb/lmdb.h:347
	Current = 0x40
	// Reserve as defined in mdb/lmdb.h:351
	Reserve = 0x10000
	// Append as defined in mdb/lmdb.h:353
	Append = 0x20000
	// Appenddup as defined in mdb/lmdb.h:355
	Appenddup = 0x40000
	// Multiple as defined in mdb/lmdb.h:357
	Multiple = 0x80000
	// CpCompact as defined in mdb/lmdb.h:366
	CpCompact = 0x01
	// ErrSuccess as defined in mdb/lmdb.h:411
	ErrSuccess = 0
	// ErrKeyExist as defined in mdb/lmdb.h:413
	ErrKeyExist = (-30799)
	// ErrNotFound as defined in mdb/lmdb.h:415
	ErrNotFound = (-30798)
	// ErrPageNotFound as defined in mdb/lmdb.h:417
	ErrPageNotFound = (-30797)
	// ErrCorrupted as defined in mdb/lmdb.h:419
	ErrCorrupted = (-30796)
	// ErrPanic as defined in mdb/lmdb.h:421
	ErrPanic = (-30795)
	// ErrVersionMismatch as defined in mdb/lmdb.h:423
	ErrVersionMismatch = (-30794)
	// ErrInvalid as defined in mdb/lmdb.h:425
	ErrInvalid = (-30793)
	// ErrMapFull as defined in mdb/lmdb.h:427
	ErrMapFull = (-30792)
	// ErrDbsFull as defined in mdb/lmdb.h:429
	ErrDbsFull = (-30791)
	// ErrReadersFull as defined in mdb/lmdb.h:431
	ErrReadersFull = (-30790)
	// ErrTlsFull as defined in mdb/lmdb.h:433
	ErrTlsFull = (-30789)
	// ErrTxnFull as defined in mdb/lmdb.h:435
	ErrTxnFull = (-30788)
	// ErrCursorFull as defined in mdb/lmdb.h:437
	ErrCursorFull = (-30787)
	// ErrPageFull as defined in mdb/lmdb.h:439
	ErrPageFull = (-30786)
	// ErrMapResized as defined in mdb/lmdb.h:441
	ErrMapResized = (-30785)
	// ErrIncompatible as defined in mdb/lmdb.h:450
	ErrIncompatible = (-30784)
	// ErrBadRslot as defined in mdb/lmdb.h:452
	ErrBadRslot = (-30783)
	// ErrBadTxn as defined in mdb/lmdb.h:454
	ErrBadTxn = (-30782)
	// ErrValueSize as defined in mdb/lmdb.h:456
	ErrValueSize = (-30781)
	// ErrBadDbi as defined in mdb/lmdb.h:458
	ErrBadDbi = (-30780)
	// LastErrcode as defined in mdb/lmdb.h:460
	LastErrcode = ErrBadDbi
)

// CursorOp as declared in mdb/lmdb.h:403
type CursorOp int32

// CursorOp enumeration from mdb/lmdb.h:403
const (
	First        CursorOp = 0
	FirstDup     CursorOp = 1
	GetBoth      CursorOp = 2
	GetBothRange CursorOp = 3
	GetCurrent   CursorOp = 4
	GetMultiple  CursorOp = 5
	Last         CursorOp = 6
	LastDup      CursorOp = 7
	Next         CursorOp = 8
	NextDup      CursorOp = 9
	NextMultiple CursorOp = 10
	NextNodup    CursorOp = 11
	Prev         CursorOp = 12
	PrevDup      CursorOp = 13
	PrevNodup    CursorOp = 14
	Set          CursorOp = 15
	SetKey       CursorOp = 16
	SetRange     CursorOp = 17
	PrevMultiple CursorOp = 18
)
