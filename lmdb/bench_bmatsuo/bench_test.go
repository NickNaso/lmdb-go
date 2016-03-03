package bench

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"runtime"
	"strings"
	"testing"
	"time"
	"unsafe"

	"github.com/bmatsuo/lmdb-go/lmdb"
)

var (
	SMALL_VAL = []byte(strings.Repeat("ABC", 10))   // 30B
	LARGE_VAL = []byte(strings.Repeat("ABC", 3*MB)) // 9MB
)

const (
	MB = 1024 * 1024
	GB = 1024 * 1024 * 1024
)

var take int

func BenchmarkPutSmall_LMDB(b *testing.B) {
	b.SetParallelism(1)

	take++
	t0 := time.Now()
	defer func() {
		log.Printf("take %d (n=%d) done in %v\n\n", take, b.N, time.Now().Sub(t0))
	}()

	log.Printf("bench put SMALL (10B), take %d (n=%d)\n", take, b.N)
	benchPut_LMDB(b, 1024*MB, SMALL_VAL)
}

func BenchmarkPutLarge_LMDB(b *testing.B) {
	b.SetParallelism(1)

	take++
	t0 := time.Now()
	defer func() {
		log.Printf("take %d (n=%d) done in %v\n\n", take, b.N, time.Now().Sub(t0))
	}()

	log.Printf("bench put LARGE (9MB), take %d (n=%d)", take, b.N)
	benchPut_LMDB(b, 20*GB, LARGE_VAL)
}

func benchPut_LMDB(b *testing.B, size int64, val []byte) {
	env := openEnv(BENCH_DB, 0)
	defer cleanDir(BENCH_DB)
	defer env.Close()

	checkErr(env.SetMapSize(size))
	txn, err := env.BeginTxn(nil, 0)
	checkErr(err)
	dbi, err := txn.OpenDBI(BENCH_DBI, lmdb.Create)
	checkErr(err)

	var t0 time.Time

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		if i == b.N-1 {
			// record time of the last run
			t0 = time.Now()
		}

		key := sortedKey()
		checkErr(txn.Put(dbi, key, val, 0))
		checkErr(txn.Commit())
		b.SetBytes(int64(len(val)))

		txn, err = env.BeginTxn(nil, 0)
		checkErr(err)
		dbi, err = txn.OpenDBI(BENCH_DBI, 0)
		checkErr(err)

		if i == b.N-1 {
			// print the duration of the last run
			log.Println("last put took:", time.Now().Sub(t0))
			txn.Abort()
		}
	}
	list(env)
	mem()
}

func checkErr(err error) {
	if err != nil {
		log.Fatalln("[ERR]", err)
	}
}

var clock int

// sortedKey returns key of (9 bytes of clock + 7 (or 14 chars) random bytes).
func sortedKey() []byte {
	clock++
	randPart := make([]byte, 7)
	if _, err := rand.Read(randPart); err != nil {
		copy(randPart, (*(*[8]byte)(unsafe.Pointer(&clock)))[:7])
	}
	key := fmt.Sprintf("%09d", clock) + hex.EncodeToString(randPart)
	return []byte(key)
}

func list(env *lmdb.Env) {
	txn, err := env.BeginTxn(nil, lmdb.Readonly)
	checkErr(err)
	defer txn.Abort()
	dbi, err := txn.OpenDBI(BENCH_DBI, 0)
	checkErr(err)

	stats, err := txn.Stat(dbi)
	checkErr(err)
	log.Printf("depth: %d, branch pg: %d, leaf pg: %d, entries: %d\n",
		stats.Depth, stats.BranchPages, stats.LeafPages, stats.Entries)
}

func mem() {
	t0 := time.Now()
	runtime.GC()
	log.Println("gc:", time.Now().Sub(t0))
	var stats runtime.MemStats
	runtime.ReadMemStats(&stats)
	log.Printf("mem: heap %dMB/%dMB, alloc %dMB\n", stats.HeapInuse/MB, stats.HeapAlloc/MB, stats.Alloc/MB)
}
