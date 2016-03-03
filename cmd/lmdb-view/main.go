package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/zenhotels/lmdb-go/lmdb"
)

// A very simple tool that lists the values. A featureful CLI tool coming soon.

var envPath = flag.String("env", "test.db", "Specify the path of an LMDB environment to read from.")
var dbiName = flag.String("db", "", "Specify the name of DB instance in the environment.")

func init() {
	flag.Parse()
	if len(*dbiName) == 0 {
		flag.Usage()
	}
	log.SetFlags(log.Lshortfile)
}

func main() {
	env, err := lmdb.EnvCreate()
	if err != nil {
		log.Fatalln(err)
	}
	defer env.Close()
	env.SetMaxDBs(1)
	if err := env.Open(*envPath, lmdb.ReadOnly, 0644); err != nil {
		log.Fatalln(err, "at", *envPath)
	}
	txn, err := env.BeginTxn(nil, lmdb.TxnReadOnly)
	if err != nil {
		log.Fatalln(err)
	}
	defer txn.Abort()
	dbi, err := txn.DbiOpen(*dbiName, lmdb.DefaultDbiFlags)
	if err != nil {
		log.Fatalln(err, "at", *dbiName)
	}
	cur, err := txn.CursorOpen(dbi)
	if err != nil {
		log.Fatalln(err)
	}
	defer cur.Close()

	log.Println("reading:", *dbiName, "from", *envPath)
	var count int
	defer func(count *int) {
		log.Println("entries read:", *count)
	}(&count)

	k, v, err := cur.Get(nil, lmdb.OpFirst)
	if err != nil {
		log.Fatalln(err)
	}
	count++
	fmt.Printf(fmtKV, count, len(k), k, len(v), v)
	for err == nil {
		k, v, err = cur.Get(k, lmdb.OpNext)
		if len(k) == 0 && len(v) == 0 {
			break
		}
		fmt.Printf(fmtKV, count, len(k), k, len(v), v)
		count++
	}
}

var fmtKV = "%06d) K(len=%d): %s, V(len=%d): %s\n"
