package bench

import (
	"encoding/hex"
	"fmt"
	"log"
	"math/rand"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/zenhotels/lmdb-go/lmdb"
)

var (
	SMALL_VAL = []byte(strings.Repeat("ABCD", 2*KB)) // 8KB
	LARGE_VAL = []byte(strings.Repeat("ABCD", 2*MB)) // 8MB
)

const (
	KB = 1024
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

	log.Printf("bench put SMALL (8KB), take %d (n=%d)\n", take, b.N)
	benchPut_LMDB(b, 10*GB, SMALL_VAL)
}

func BenchmarkPutLarge_LMDB(b *testing.B) {
	b.SetParallelism(1)

	take++
	t0 := time.Now()
	defer func() {
		log.Printf("take %d (n=%d) done in %v\n\n", take, b.N, time.Now().Sub(t0))
	}()

	log.Printf("bench put LARGE (8MB), take %d (n=%d)", take, b.N)
	benchPut_LMDB(b, 20*GB, LARGE_VAL)
}

func benchPut_LMDB(b *testing.B, size uint, val []byte) {
	env := openEnv(BENCH_DB, lmdb.DefaultEnvFlags)
	// env := openEnv(fmt.Sprintf("W%d", len(val)/1024)+BENCH_DB, lmdb.DefaultEnvFlags)
	// defer cleanDir(BENCH_DB)
	defer env.Close()

	checkErr(env.SetMapSize(size))
	txn, err := env.BeginTxn(nil, lmdb.DefaultTxnFlags)
	checkErr(err)
	dbi, err := txn.DbiOpen(BENCH_DBI, lmdb.DbiCreate)
	checkErr(err)

	var t0 time.Time

	b.SetBytes(int64(len(val)))
	b.ResetTimer()
	b.ReportAllocs()

	batch := 400

	for i := 0; i < b.N; i++ {
		if i == b.N-1 {
			// record time of the last run
			t0 = time.Now()
		}

		// key := decKey()
		key := randKey()
		checkErr(txn.Put(dbi, key, val, lmdb.DefaultWriteFlags))

		if i%batch == 0 {
			checkErr(txn.Commit())
			txn, err = env.BeginTxn(nil, lmdb.DefaultTxnFlags)
			checkErr(err)
			dbi, err = txn.DbiOpen(BENCH_DBI, lmdb.DefaultDbiFlags)
			checkErr(err)
		}

		if i == b.N-1 {
			checkErr(txn.Commit())
			// print the duration of the last run
			log.Println("last put took:", time.Now().Sub(t0))
			log.Println("last key:", string(key))
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

var src = rand.New(rand.NewSource(time.Now().UnixNano()))

// sortedKey returns key of (9 bytes of clock + 7 (or 14 chars) random bytes).
func sortedKey() []byte {
	clock++
	randPart := make([]byte, 7)
	src.Read(randPart)
	key := fmt.Sprintf("%09d", clock) + hex.EncodeToString(randPart)
	return []byte(key)
}

func randKey() []byte {
	randPart := make([]byte, 16)
	src.Read(randPart)
	return []byte(hex.EncodeToString(randPart))
}

func decKey() []byte {
	clock++
	return []byte(fmt.Sprintf("%09d", clock))
}

func useKey(n int) []byte {
	return []byte(fmt.Sprintf("%09d", n))
}

func list(env lmdb.Env) {
	txn, err := env.BeginTxn(nil, lmdb.TxnReadOnly)
	checkErr(err)
	dbi, err := txn.DbiOpen(BENCH_DBI, lmdb.DefaultDbiFlags)
	checkErr(err)

	stats, err := txn.Stat(dbi)
	checkErr(err)
	log.Printf("depth: %d, branch pg: %d, leaf pg: %d, entries: %d\n",
		stats.Depth(), stats.BranchPages(), stats.LeafPages(), stats.Entries())
}

func mem() {
	t0 := time.Now()
	runtime.GC()
	log.Println("gc:", time.Now().Sub(t0))
	var stats runtime.MemStats
	runtime.ReadMemStats(&stats)
	log.Printf("mem: heap %dMB/%dMB, alloc total %dMB\n", stats.HeapInuse/MB, stats.HeapAlloc/MB, stats.TotalAlloc/MB)
}
