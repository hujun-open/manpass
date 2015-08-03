// passsql
package passsql

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"manpassd/common"
	"os"
	"path/filepath"
	"time"
)

const (
	CREATEDBSQL = `
	drop table if exists %[1]s;
	create table %[1]s (
		meta_id blob not null ,/* sha256 of key and uname*/
		pass_rev INTEGER not null,
		meta blob not null, /* key is the meta data associated with username/password, could be a URL */
		uname blob not null, /* username, encrypted*/
		pass blob not null, /* current password, encrypted */
		pass_time timestamp default CURRENT_TIMESTAMP,
		primary key (meta_id,pass_rev)
	);
	
	`
	INSERTSQL = `
	insert into %[1]s (uname,pass,meta,meta_id,pass_rev) select ?,?,?,?,
	case 
		when exists
			(select pass_rev from %[1]s where meta_id=?)
		then 
			max(pass_rev)+1
		else 0
	end
	from %[1]s where meta_id=?
	`
	UPDATESQL = `
	update or fail %[1]s set meta=?, uname=?, pass=? where meta_id=? and pass_rev=?
	`
)

type PassDB struct {
	PDB *sql.DB
}

type PassRecord struct {
	Meta_id   string
	Pass_rev  int
	Meta      string
	Uname     string
	Pass      string
	Pass_time time.Time
}

func InitDB(filename string) (*PassDB, error) {
	//remove exising file and create a sqlite3 file
	os.Remove(filename)
	db, err := sql.Open("sqlite3", filename)
	if err != nil {
		return nil, err
	}
	pdb := new(PassDB)
	pdb.PDB = db
	return pdb, nil
}

func LoadDB(filename string) (*PassDB, error) {
	//load an existing sqlite3 file
	db, err := sql.Open("sqlite3", filename)
	if err != nil {
		return nil, err
	}
	pdb := new(PassDB)
	pdb.PDB = db
	return pdb, nil
}

func (pdb PassDB) InitTable(tablename string) error {
	//create a new table in the specified file, drop the existing table
	//this could also be used to remove all records for a given table
	tx, err := pdb.PDB.Begin()
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = pdb.PDB.Exec(fmt.Sprintf(CREATEDBSQL, tablename))
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (pdb PassDB) Insert(tablename string, pr PassRecord) error {
	tx, err := pdb.PDB.Begin()
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = tx.Exec(fmt.Sprintf(INSERTSQL, tablename), pr.Uname, pr.Pass, pr.Meta, pr.Meta_id, pr.Meta_id, pr.Meta_id)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil

}

func (pdb PassDB) ReplaceAll(tablename string, prlist []PassRecord) error {
	tx, err := pdb.PDB.Begin()
	if err != nil {
		tx.Rollback()
		return err
	}
	for _, pr := range prlist {
		_, err = tx.Exec(fmt.Sprintf(UPDATESQL, tablename), pr.Meta, pr.Uname, pr.Pass, pr.Meta_id, pr.Pass_rev)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (pdb PassDB) GetRecord(tablename string, meta_id string, pass_rev int) (*PassRecord, error) {
	//return one specified record based on meta-id, return nil if not found
	//if pass_rev<0, then return the latest one
	r := new(PassRecord)
	var sql_template string
	var err error
	var rows *sql.Rows
	if pass_rev < 0 {
		sql_template = fmt.Sprintf(`select * from %[1]s where meta_id=? and pass_rev in (select max(pass_rev) from %[1]s where meta_id=?)`, tablename)
		rows, err = pdb.PDB.Query(sql_template, meta_id, meta_id)
	} else {
		sql_template = fmt.Sprintf(`select * from %[1]s where meta_id=? and pass_rev=?`, tablename)
		rows, err = pdb.PDB.Query(sql_template, meta_id, pass_rev)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	i := 0
	for rows.Next() {
		err := rows.Scan(&r.Meta_id, &r.Pass_rev, &r.Meta, &r.Uname, &r.Pass, &r.Pass_time)
		if err != nil {
			return nil, err
		}
		i += 1
	}
	if i == 0 {
		return nil, err
	}
	return r, nil
}

func (pdb PassDB) GetAll(tablename string) ([]PassRecord, error) {
	// return all records in the table
	rlist := []PassRecord{}
	r := new(PassRecord)
	rows, err := pdb.PDB.Query(fmt.Sprintf("select * from %s;", tablename))
	if err != nil {
		return nil, err

	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&r.Meta_id, &r.Pass_rev, &r.Meta, &r.Uname, &r.Pass, &r.Pass_time)
		if err == nil {
			rlist = append(rlist, *r)
		}
		r = new(PassRecord)
	}
	return rlist, nil
}

func (pdb PassDB) GetAllLatest(tablename string) ([]PassRecord, error) {
	rlist := []PassRecord{}
	rows, err := pdb.PDB.Query(fmt.Sprintf("select distinct(meta_id) from %s", tablename))
	if err != nil {
		return nil, err

	}
	defer rows.Close()
	var mid string
	for rows.Next() {
		err := rows.Scan(&mid)
		if err == nil {
			r, err := pdb.GetRecord(tablename, mid, -1)
			if err == nil {
				rlist = append(rlist, *r)
			}

		}
	}
	return rlist, nil

}

func (pdb PassDB) GetAllMetaId(tablename string) ([]string, error) {
	var metaid_list []string
	rows, err := pdb.PDB.Query(fmt.Sprintf("select distinct(meta_id) from %s", tablename))
	if err != nil {
		return nil, err

	}
	defer rows.Close()
	var mid string
	for rows.Next() {
		err := rows.Scan(&mid)
		if err == nil {
			metaid_list = append(metaid_list, mid)
		}
	}
	return metaid_list, nil

}

func (pdb PassDB) GetAllRevForMetaId(tablename string, meta_id string) ([]PassRecord, error) {
	//return all records for a given meta-id
	rlist := []PassRecord{}
	r := new(PassRecord)
	rows, err := pdb.PDB.Query(fmt.Sprintf("select * from %s where meta_id=?;", tablename), meta_id)
	if err != nil {
		return nil, err

	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&r.Meta_id, &r.Pass_rev, &r.Meta, &r.Uname, &r.Pass, &r.Pass_time)
		if err == nil {
			rlist = append(rlist, *r)
		}
		r = new(PassRecord)
	}
	return rlist, nil
}

func (pdb PassDB) RemovePassForRev(tablename string, meta_id string, pass_rev int) error {
	//remove a specfic record with speicifed meta-id and pass_rev
	tx, err := pdb.PDB.Begin()
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = tx.Exec(fmt.Sprintf("delete from %[1]s where meta_id=? and pass_rev=?", tablename), meta_id, pass_rev)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (pdb PassDB) RemovePass(tablename string, meta_id string) error {
	//remove all records of specified meta-id
	tx, err := pdb.PDB.Begin()
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = tx.Exec(fmt.Sprintf("delete from %[1]s where meta_id=? ", tablename), meta_id)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

//func (pdb PassDB) update(tablename string, meta_id []byte, new_pr PassRecord, update_pass bool) (sql.Result, error) {
//	tx, err := pdb.PDB.Begin()
//	var r sql.Result
//	if err != nil {
//		tx.Rollback()
//		return r, err
//	}
//	if update_pass == false {
//		r, err = tx.Exec(fmt.Sprintf("update or abort %[1]s set meta=?,uname=? where meta_id=?", tablename), new_pr.meta, new_pr.uname, new_pr.meta_id)
//	} else {
//		r, err = tx.Exec(fmt.Sprintf("update or abort %[1]s set meta=?,uname=?,pass=?,pass_time=?,old_pass=pass where meta_id=?", tablename), new_pr.meta, new_pr.uname, new_pr.pass, time.Now(), meta_id)
//	}

//	if err != nil {
//		tx.Rollback()
//		return r, err
//	}
//	err = tx.Commit()
//	if err != nil {
//		return r, err
//	}
//	return r, nil

//}

func (pdb PassDB) PrintAll() {
	var tablename string
	rows, err := pdb.PDB.Query("select name from sqlite_master where type='table' order by name;")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer rows.Close()
	i := 0
	for rows.Next() {
		err := rows.Scan(&tablename)
		if err != nil {
			log.Fatal(err)
		} else {
			i += 1
			pdb.PrintTable(tablename)
		}
	}
	log.Printf("there are total %d tables", i)
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
}

func (pdb PassDB) PrintTable(tablename string) {
	rlist, err := pdb.GetAll(tablename)
	if err != nil {
		log.Fatal(err)
	}
	for _, r := range rlist {
		log.Println(r)
	}
}

func InitPassDB(uname string) error {
	confdir := common.GetConfDir(uname)
	fi, err := os.Stat(confdir)
	if err != nil {
		return err
	}
	if !fi.IsDir() {
		return fmt.Errorf("%s is not a directory", confdir)
	}
	dbfname := filepath.Join(confdir, uname+".db")
	pdb, err := InitDB(dbfname)
	if err != nil {
		return err
	}
	pdb.InitTable(uname)
	return nil

}
