all:
	cgogen mdb.yml

clean:
	rm -f mdb/cgo_helpers.go mdb/cgo_helpers.c mdb/cgo_helpers.h mdb/doc.go mdb/types.go mdb/const.go
	rm -f mdb/mdb.go

test:
	cd mdb && go build
	cd lmdb && go test
