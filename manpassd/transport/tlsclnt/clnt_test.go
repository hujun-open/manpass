// clnt_test
package tlsclnt

import (
	//"crypto/tls"
	//"fmt"
	//"net/http"
	"testing"
)

func TestClient(t *testing.T) {
	c, err := GetClient()
	if err != nil {
		t.Fatal(err)
		t.Fatal("failed to create client")
	}
	//fmt.Println(c.Transport.(*http.Transport).TLSClientConfig.ClientCAs.Subjects())
	resp, err := c.Get("https://127.0.0.1/")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(resp)

}
