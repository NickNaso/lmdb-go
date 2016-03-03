package bench

import (
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/bmatsuo/lmdb-go/lmdb"
)

func init() {
	log.SetFlags(log.Lshortfile)
}

func TestMain(m *testing.M) {
	initDir(BENCH_DB)
	os.Exit(m.Run())
}

const (
	BENCH_DB  = "bench.db"
	BENCH_DBI = "bench_dbi"
)

func initDir(dir string) error {
	return os.MkdirAll(dir, 0755)
}

func cleanDir(dir string) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Println("[ERR] cleanup err", err)
			return err
		}
		if !info.IsDir() {
			if err := os.Remove(path); err != nil {
				log.Println("[ERR] cannot remove", path)
			}
		}
		return nil
	})
}

func nukeDir(dir string) error {
	return os.RemoveAll(dir)
}

func openEnv(db string, flags uint) *lmdb.Env {
	env, err := lmdb.NewEnv()
	if err != nil {
		log.Fatalln(err)
	}
	if err := env.SetMaxDBs(1024); err != nil {
		log.Fatalln(err)
	}
	if _, err := os.Stat(db); err != nil && os.IsNotExist(err) {
		if err := initDir(db); err != nil {
			log.Fatalln(err)
		}
	}
	if err := env.Open(db, flags, 0644); err != nil {
		log.Fatalln(err)
	}
	return env
}
