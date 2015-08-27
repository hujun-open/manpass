// sql_test

package passsql

import (
	//	"io/ioutil"
	"testing"
)

//func TestSqldump(t *testing.T) {
//	dbfile := "D:\\hujun\\Dropbox\\manpass_hujun_conf\\hujun.db"
//	fdb, err := LoadDB(dbfile)
//	if err != nil {
//		t.Fatal(err)
//	}
//	js, err := fdb.Dump("hujun", "json")
//	if err != nil {
//		t.Fatal(err)
//	}
//	ioutil.WriteFile("D:\\hujun\\Dropbox\\manpass_hujun_conf\\hujun.db.dump", []byte(js), 0777)
//	newdbfile := "new.db"
//	newdb, err := InitDB(newdbfile)
//	if err != nil {
//		t.Fatal(err)
//	}
//	newdb.InitTable("hujun")
//	err = newdb.Import("hujun", js, "json")
//	if err != nil {
//		t.Fatal(err)
//	}

//}

func TestSql(t *testing.T) {
	dbfile := "test.db"
	tablename := "testuser"
	t.Log("start to test InitDB\n")
	passdb, err := InitDB(dbfile)
	if err != nil {
		t.Fatal(err)
		t.Fatalf("Failed to init db:%q", dbfile)
	}
	t.Log("start to test IniTable\n")
	err = passdb.InitTable(tablename)
	if err != nil {
		t.Fatal(err)
		t.Fatalf("Failed to init table:%q", tablename)
	}
	t.Log(passdb)
	t.Log("start to test insert")
	pr := PassRecord{
		Meta_id: "metaid-1",
		Meta:    "gmail.com",
		Uname:   "abc1@gmail.com",
		Pass:    "abc123",
		Kgroup:  "normal",
		Remark:  "remark-1",
	}
	pr2 := PassRecord{
		Meta_id: "metaid-2",
		Meta:    "gmail.com",
		Uname:   "abc2@gmail.com",
		Pass:    "abc256",
		Kgroup:  "normal",
		Remark:  "remark-2",
	}
	pr3 := PassRecord{
		Meta_id: "metaid-2",
		Meta:    "gmail.com",
		Uname:   "abc2@gmail.com",
		Pass:    "abc512",
		Kgroup:  "hidden",
		Remark:  "remark-3",
	}
	err = passdb.Insert(tablename, pr)
	if err != nil {
		t.Fatal(err)
	}
	err = passdb.Insert(tablename, pr2)
	if err != nil {
		t.Fatal(err)
	}
	err = passdb.Insert(tablename, pr3)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("start to test getAllLatest\n")
	rlist, err := passdb.GetAllLatest(tablename)
	if err != nil {
		t.Fatal(err)
		t.Fatalf("Failed to get all records from table %q", tablename)
	}
	t.Log(rlist)
	t.Logf("there are %d latest records ", len(rlist))
	t.Log("start to test getAll\n")
	rlist, err = passdb.GetAll(tablename)
	if err != nil {
		t.Fatal(err)
		t.Fatalf("Failed to get all records from table %q", tablename)
	}
	t.Logf("There are total %d records in table %q", len(rlist), tablename)
	t.Log("start to test getRecord for latest pass\n")
	qpr, err := passdb.GetRecord(tablename, pr2.Meta_id, -1)
	if err != nil {
		t.Fatal(err)
		t.Fatalf("Failed to get record for %q", pr2.Meta_id)
	}
	t.Log(qpr)

	t.Log("start to test getRecord for specfic record\n")
	qpr, err = passdb.GetRecord(tablename, pr2.Meta_id, 0)
	if err != nil {
		t.Fatal(err)
		t.Fatalf("Failed to get record for %q", pr2.Meta_id)
	}
	t.Log(qpr)

	t.Logf("start negative test for getRecord")
	qpr, err = passdb.GetRecord(tablename, "xix", -1)
	if err != nil {
		t.Fatal(err)
	}
	if qpr != nil {
		t.Fatalf("%q shouldn't exisit", "xix")
	}
	t.Log("start to test getRecord for getAllRevForMetaId\n")
	rlist, err = passdb.GetAllRevForMetaId(tablename, pr2.Meta_id)
	if err != nil {
		t.Fatal(err)
	}
	if len(rlist) != 2 {
		t.Fatal("failed to return two results")
	}
	passdb.PrintAll()
	t.Log("start to test getRecord for removePassForRev")
	err = passdb.RemovePassForRev(tablename, pr2.Meta_id, 1)
	if err != nil {
		t.Fatal(err)
	}
	r, err := passdb.GetRecord(tablename, pr2.Meta_id, 1)
	if err != nil {
		t.Fatal(err)
	}
	if r != nil {
		t.Fatal(err)
	}
	err = passdb.Insert(tablename, pr3)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("start to test getRecord for removePass")
	err = passdb.RemovePass(tablename, pr2.Meta_id)
	if err != nil {
		t.Fatal(err)
	}
	rlist, err = passdb.GetAllRevForMetaId(tablename, pr2.Meta_id)
	if err != nil || len(rlist) > 0 {
		t.Fatal(err)
	}

}
