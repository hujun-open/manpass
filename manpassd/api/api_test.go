// api_test.go
package api

import (
	"manpassd/passsql"
	"manpassd/transport/tlsvr"
	"testing"
)

func TestSvr(t *testing.T) {
	dbfile := "d:\\temp\\1.db"
	tablename := "hujun"
	t.Log("start to test InitDB\n")
	passdb, err := passsql.InitDB(dbfile)
	if err != nil {
		t.Fatal(err)
		t.Fatalf("Failed to init db:%q", dbfile)
	}
	t.Log("start to test IniTable\n")
	passdb.InitTable(tablename)
	svr, err := NewClientAPISVR("127.0.0.1", 9000, "hujun", []byte("zifan123"), *passdb, tablename)
	if err != nil {
		t.Fatal(err)
	}
	err = tlsvr.ListenWithConfig(svr.HttpSvr)
	if err != nil {
		t.Fatal(err)
	}
	var c chan int
	<-c
}
