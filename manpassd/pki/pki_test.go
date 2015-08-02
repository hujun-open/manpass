// pki_test
package pki

import (
	//	"io/ioutil"
	"fmt"
	//	"reflect"
	"testing"
)

func TestECDSA(t *testing.T) {
	ca := new(ECCertKey)
	//	var ca RSACertKey
	err := GenerateROOTCA(ca, "HJROOTCA")
	if err != nil {
		t.Fatal(err)
	}
	PrintCert(ca.cert)
	ee := new(ECCertKey)
	err = CreateCertWithCA(ca, ee)
	if err != nil {
		t.Fatal(err)
	}
	PrintCert(ee.cert)
	enc, err := ee.EncPemKey("disk")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(enc))
	fmt.Println(string(ee.PemCert()))
	enc, err = ee.EncPkg("disk")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(enc))
	newck, err := LoadCertKeyFromEncPkg(enc, "disk")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("loaded--------------------------------")
	PrintCert((*newck).GetCert())

}

func TestRSA(t *testing.T) {
	ca := new(RSACertKey)
	//	var ca RSACertKey
	err := GenerateROOTCA(ca, "HJROOTCA")
	if err != nil {
		t.Fatal(err)
	}
	PrintCert(ca.cert)
	ee := new(RSACertKey)
	err = CreateCertWithCA(ca, ee)
	if err != nil {
		t.Fatal(err)
	}
	PrintCert(ee.cert)
	enc, err := ee.EncPemKey("disk")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(enc))
	fmt.Println(string(ee.PemCert()))
	enc, err = ee.EncPkg("disk")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(enc))
	newck, err := LoadCertKeyFromEncPkg(enc, "disk")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("loaded--------------------------------")
	PrintCert((*newck).GetCert())
}
