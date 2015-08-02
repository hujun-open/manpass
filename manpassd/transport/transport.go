// loadcerts
package transport

import (
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
	//	"log"
	"manpassd/common"
	"manpassd/passcrypto"
	"path/filepath"
)

func LoadCertsKeys(uname string, upass []byte) (ca []byte, ee []byte, pkey []byte, err error) {
	user_dir := common.GetConfDir(uname)
	cafs, err := ioutil.ReadFile(filepath.Join(user_dir, "ca.cert"))
	if err != nil {
		return nil, nil, nil, err
	}
	ca, err = passcrypto.DecryptMeBase32(string(cafs), upass)
	if err != nil {
		return nil, nil, nil, err
	}
	ee_pem, err := ioutil.ReadFile(filepath.Join(user_dir, "ee.cert"))
	if err != nil {
		return nil, nil, nil, err
	}
	ee_blk, _ := pem.Decode(ee_pem)
	ee = ee_blk.Bytes
	pkeyfs, err := ioutil.ReadFile(filepath.Join(user_dir, "ee.key"))
	if err != nil {
		return nil, nil, nil, err
	}
	pkey_pem, _ := pem.Decode(pkeyfs)
	pkey, err = x509.DecryptPEMBlock(pkey_pem, upass)
	if err != nil {
		return nil, nil, nil, err
	}
	return
}

func GetConfig(ca_f []byte, ee_f []byte, key_f []byte) (*tls.Config, error) {
	ca, err := x509.ParseCertificate(ca_f)
	if err != nil {
		return nil, err
	}
	pkey, err := x509.ParsePKCS1PrivateKey(key_f)
	if err != nil {
		return nil, err
	}
	ca_pool := x509.NewCertPool()
	ca_pool.AddCert(ca)
	ee_cert := tls.Certificate{
		Certificate: [][]byte{ee_f},
		PrivateKey:  pkey,
	}
	config := new(tls.Config)
	config.ClientAuth = tls.RequireAndVerifyClientCert
	config.Certificates = []tls.Certificate{ee_cert}
	config.ClientCAs = ca_pool
	config.RootCAs = ca_pool
	config.Rand = rand.Reader
	config.BuildNameToCertificate()
	return config, nil
}

//func L2() (ca []byte, ee []byte, pkey []byte, err error) {
//	ca, err = ioutil.ReadFile("d:\\temp\\tls\\root.cer")
//	if err != nil {
//		return nil, nil, nil, err
//	}
//	ee, err = ioutil.ReadFile("d:\\temp\\tls\\a2.cer")
//	if err != nil {
//		return nil, nil, nil, err
//	}
//	pkey, err = ioutil.ReadFile("d:\\temp\\tls\\a2_key.der")
//	if err != nil {
//		return nil, nil, nil, err
//	}
//	return
//}
