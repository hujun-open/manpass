// svr_test
package tlsvr

import (
	//"fmt"
	"testing"
)

//func cSvr(t *testing.T) {

//}

func TestSvr(t *testing.T) {
	svr, err := GetServer("0.0.0.0", 9000, "hujun", []byte("zifan123"))
	if err != nil {
		t.Fatal(err)
		t.Fatal("failed to create server")
	}
	t.Log("Svr started")
	//err = svr.ListenAndServeTLS("d:\\temp\\tls\\a1.cer", "d:\\temp\\tls\\a1_key.der")
	err = ListenWithConfig(svr)
	if err != nil {
		t.Fatal(err)
	}
	var c chan int
	<-c
	t.Log("done.")
}
