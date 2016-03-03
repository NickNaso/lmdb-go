package bench

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/zenhotels/lmdb-go/lmdb"
)

func BenchmarkGetSmall_LMDB(b *testing.B) {
	b.SetParallelism(1)

	take++
	t0 := time.Now()
	defer func() {
		log.Printf("take %d (n=%d) done in %v\n\n", take, b.N, time.Now().Sub(t0))
	}()

	log.Printf("bench get SMALL (4kB), take %d (n=%d)\n", take, b.N)
	benchGet_LMDB(b, 10*GB, 510101, SMALL_VAL)
}

func BenchmarkGetLarge_LMDB(b *testing.B) {
	b.SetParallelism(1)

	take++
	t0 := time.Now()
	defer func() {
		log.Printf("take %d (n=%d) done in %v\n\n", take, b.N, time.Now().Sub(t0))
	}()

	log.Printf("bench get LARGE (9MB), take %d (n=%d)", take, b.N)
	benchGet_LMDB(b, 20*GB, 601, LARGE_VAL)
}

func benchGet_LMDB(b *testing.B, size uint, count int, val []byte) {
	env := openEnv(fmt.Sprintf("W%d", len(val)/1024)+BENCH_DB, lmdb.ReadOnly)
	defer env.Close()

	checkErr(env.SetMapSize(size))
	txn, err := env.BeginTxn(nil, lmdb.TxnReadOnly)
	checkErr(err)
	defer txn.Abort()
	dbi, err := txn.DbiOpen(BENCH_DBI, lmdb.DefaultDbiFlags)
	checkErr(err)

	b.ResetTimer()
	b.ReportAllocs()

	var found int
	var missed int
	for i := 0; i < b.N; i++ {
		key := useKey(src.Intn(count + count/20))
		v, err := txn.Get(dbi, key)
		if err != nil {
			if err == lmdb.ErrNotFound {
				missed++
			} else {
				checkErr(err)
			}
		} else {
			found++
			if len(v) != len(val) {
				log.Fatalln("expected", len(val), "but got", len(v))
			}
		}
	}
	log.Println("missed:", missed, "found:", found)
	list(env)
	mem()
}
