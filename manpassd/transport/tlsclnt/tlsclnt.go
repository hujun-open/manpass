// tlsclnt
package tlsclnt

import (
	"manpassd/transport"
	"net/http"
)

func GetClient() (client *http.Client, err error) {
	ca_f, ee_f, key_f, err := transport.L2()
	if err != nil {
		return nil, err
	}
	config, err := transport.GetConfig(ca_f, ee_f, key_f)
	if err != nil {
		return nil, err
	}
	config.InsecureSkipVerify = false
	config.SkipVerifyHostname = true
	trans := http.Transport{
		TLSClientConfig: config,
	}
	client = new(http.Client)
	client.Transport = &trans
	return client, nil

}
