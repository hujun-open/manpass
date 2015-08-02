// pki
package pki

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base32"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"manpassd/common"
	"manpassd/passcrypto"
	"math/big"
	"os"
	"path/filepath"
	"time"
)

const (
	validYears = 10
	pkgTypeStr = "KEY AND CERT"
)

var supportedKeyTypes = map[string]bool{
	"EC PRIVATE KEY":  true,
	"RSA PRIVATE KEY": true,
}

type CertKey interface {
	RawPubKey() []byte
	CreateCAKey() error
	CreateEEKey() error
	SetCert(c x509.Certificate)
	PrivKey() interface{}
	PubKey() interface{}
	GetCert() x509.Certificate
	PemKey() []byte
	EncPemKey(passwd []byte) ([]byte, error)
	PemCert() []byte
	EncPkg(passwd string) ([]byte, error) //an enceypted PEM encoded cert+key
	EncBaseCert(passwd []byte) ([]byte, error)
}

type RSACertKey struct {
	cert x509.Certificate
	key  *rsa.PrivateKey
}

func (ck *RSACertKey) PrivKey() interface{} {

	return ck.key
}
func (ck *RSACertKey) PubKey() interface{} {
	return &ck.key.PublicKey
}

func (ck *RSACertKey) RawPubKey() []byte {
	var bs []byte
	k := ck.key
	bs = k.PublicKey.N.Bytes()
	return bs
}
func (ck *RSACertKey) CreateCAKey() error {
	privkey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return err
	}
	ck.key = privkey
	return nil
}

func (ck *RSACertKey) CreateEEKey() error {
	privkey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}
	ck.key = privkey
	return nil
}

func (ck *RSACertKey) SetCert(c x509.Certificate) {
	ck.cert = c
}
func (ck *RSACertKey) GetCert() x509.Certificate {
	return ck.cert
}

func (ck *RSACertKey) PemKey() []byte {
	return pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(ck.key),
		},
	)
}

func (ck *RSACertKey) PemCert() []byte {
	return pem.EncodeToMemory(
		&pem.Block{
			Type:  "CERTIFICATE",
			Bytes: ck.cert.Raw,
		},
	)
}

func (ck *RSACertKey) EncPemKey(passwd []byte) ([]byte, error) {
	//kpem := ck.PemKey()
	kpem := x509.MarshalPKCS1PrivateKey(ck.key)
	encblock, err := x509.EncryptPEMBlock(rand.Reader, "RSA PRIVATE KEY", kpem, passwd, x509.PEMCipherAES128)
	if err != nil {
		return nil, err
	}
	return pem.EncodeToMemory(encblock), nil
}

func (ck *RSACertKey) EncBaseCert(passwd []byte) ([]byte, error) {
	//encrypt DER encoded cert in base32
	r, err := passcrypto.EncryptMeBase32(ck.cert.Raw, passwd)
	if err != nil {
		return nil, err
	}
	return []byte(r), nil

}

func (ck *RSACertKey) EncPkg(passwd string) ([]byte, error) {
	var pkgpem []byte
	pkgpem = append(pkgpem, ck.PemKey()...)
	pkgpem = append(pkgpem, ck.PemCert()...)
	encblock, err := x509.EncryptPEMBlock(rand.Reader, pkgTypeStr, pkgpem, []byte(passwd), x509.PEMCipherAES128)
	if err != nil {
		return nil, err
	}
	return pem.EncodeToMemory(encblock), nil
}

type ECCertKey struct {
	cert x509.Certificate
	key  *ecdsa.PrivateKey
}

func (ck *ECCertKey) EncBaseCert(passwd []byte) ([]byte, error) {
	//encrypt DER encoded cert in base32
	r, err := passcrypto.EncryptMeBase32(ck.cert.Raw, passwd)
	if err != nil {
		return nil, err
	}
	return []byte(r), nil

}

func (ck *ECCertKey) PrivKey() interface{} {
	return ck.key
}

func (ck *ECCertKey) PubKey() interface{} {
	return &ck.key.PublicKey
}

func (ck ECCertKey) RawPubKey() []byte {
	buf := []byte{}
	buf = append(buf, ck.key.PublicKey.X.Bytes()...)
	buf = append(buf, ck.key.PublicKey.Y.Bytes()...)
	buf = append(buf, ck.key.PublicKey.Curve.Params().B.Bytes()...)
	buf = append(buf, ck.key.PublicKey.Curve.Params().P.Bytes()...)
	buf = append(buf, ck.key.PublicKey.Curve.Params().N.Bytes()...)
	buf = append(buf, ck.key.PublicKey.Curve.Params().Gx.Bytes()...)
	buf = append(buf, ck.key.PublicKey.Curve.Params().Gy.Bytes()...)
	return buf
}

func (ck *ECCertKey) CreateCAKey() error {
	c := elliptic.P521()
	privkey, err := ecdsa.GenerateKey(c, rand.Reader)
	if err != nil {
		return err
	}
	ck.key = privkey
	return nil
}

func (ck *ECCertKey) CreateEEKey() error {
	c := elliptic.P384()
	privkey, err := ecdsa.GenerateKey(c, rand.Reader)
	if err != nil {
		return err
	}
	ck.key = privkey
	return nil
}

func (ck *ECCertKey) SetCert(c x509.Certificate) {
	ck.cert = c
}

func (ck *ECCertKey) GetCert() x509.Certificate {
	return ck.cert
}

func (ck *ECCertKey) PemKey() []byte {
	derk, _ := x509.MarshalECPrivateKey(ck.key)
	return pem.EncodeToMemory(
		&pem.Block{
			Type:  "EC PRIVATE KEY",
			Bytes: derk,
		},
	)
}

func (ck *ECCertKey) PemCert() []byte {
	return pem.EncodeToMemory(
		&pem.Block{
			Type:  "CERTIFICATE",
			Bytes: ck.cert.Raw,
		},
	)
}

func (ck *ECCertKey) EncPemKey(passwd []byte) ([]byte, error) {
	//kpem := ck.PemKey()

	kpem, err := x509.MarshalECPrivateKey(ck.key)
	if err != nil {
		return nil, err
	}
	encblock, err := x509.EncryptPEMBlock(rand.Reader, "EC PRIVATE KEY", kpem, passwd, x509.PEMCipherAES128)
	if err != nil {
		return nil, err
	}
	return pem.EncodeToMemory(encblock), nil
}

func (ck *ECCertKey) EncPkg(passwd string) ([]byte, error) {
	var pkgpem []byte
	pkgpem = append(pkgpem, ck.PemKey()...)
	pkgpem = append(pkgpem, ck.PemCert()...)
	encblock, err := x509.EncryptPEMBlock(rand.Reader, pkgTypeStr, pkgpem, []byte(passwd), x509.PEMCipherAES128)
	if err != nil {
		return nil, err
	}
	return pem.EncodeToMemory(encblock), nil
}

func GenerateROOTCA(ck CertKey, CN string) error {
	//the CommonName in result cert is b32 encoding of sha1 of pub key
	err := ck.CreateCAKey()
	if err != nil {
		return err
	}
	subkeyid := sha1.Sum(ck.RawPubKey())
	template := &x509.Certificate{
		IsCA: true,
		BasicConstraintsValid: true,
		SubjectKeyId:          subkeyid[:],
		SerialNumber:          big.NewInt(1),
		Subject: pkix.Name{
			CommonName: CN,
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().AddDate(validYears, 0, 0),
		// see http://golang.org/pkg/crypto/x509/#KeyUsage
		KeyUsage: x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
	}
	cert, err := x509.CreateCertificate(rand.Reader, template, template, ck.PubKey(), ck.PrivKey())
	if err != nil {
		return err
	}
	parsed_cert, err := x509.ParseCertificate(cert)
	if err != nil {
		return err
	}

	ck.SetCert(*parsed_cert)

	return nil
}

func sdbm(s []byte) int64 {
	var hash int64 = 0
	for _, v := range s {
		hash = int64(v) + (hash << 6) + (hash << 16) - hash
	}
	if hash <= 0 {
		hash = hash * -1
	}
	return hash
}

func CreateCertWithCA(ca CertKey, ee CertKey) error {
	//ee will be signed by ca
	ee.CreateEEKey()
	subkeyid := sha1.Sum(ee.RawPubKey())
	cname := base32.StdEncoding.EncodeToString(subkeyid[:])
	sn := sdbm(subkeyid[:])
	template := &x509.Certificate{
		IsCA: false,
		BasicConstraintsValid: true,
		SubjectKeyId:          subkeyid[:],
		SerialNumber:          big.NewInt(sn),
		Subject: pkix.Name{
			CommonName: cname,
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().AddDate(validYears, 0, 0),
		// see http://golang.org/pkg/crypto/x509/#KeyUsage
		KeyUsage:    x509.KeyUsageDigitalSignature,
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
	}
	ca_cert := ca.GetCert()
	cert, err := x509.CreateCertificate(rand.Reader, template, &ca_cert, ee.PubKey(), ca.PrivKey())
	if err != nil {
		return err
	}
	parsed_cert, err := x509.ParseCertificate(cert)
	if err != nil {
		return err
	}
	ee.SetCert(*parsed_cert)
	return nil
}

func PrintCert(cert x509.Certificate) {
	//print content of cert
	fmt.Println("----------------")
	fmt.Printf("version:%d\n", cert.Version)
	fmt.Printf("Serial Number:%d\n", cert.SerialNumber)
	fmt.Printf("Signature Alg:%s\n", cert.SignatureAlgorithm)
	fmt.Printf("Issuer:%s\n", cert.Issuer)
	fmt.Printf("Valid From:%s\n", cert.NotBefore)
	fmt.Printf("Valid To:%s\n", cert.NotAfter)
	fmt.Printf("Subject:%s\n", cert.Subject)
	fmt.Printf("CA?:%s\n", cert.IsCA)
	fmt.Printf("KeyUsage:%s\n", cert.KeyUsage)
	fmt.Printf("ExtKeyUsage:%s\n", cert.ExtKeyUsage)
	fmt.Printf("SubjectAltNames(DNS):%s\n", cert.DNSNames)
	fmt.Printf("SubjectAltNames(IP Addres):%s\n", cert.IPAddresses)
	fmt.Printf("SubjectAltNames(Email):%s\n", cert.EmailAddresses)

}

func LoadCertKeyFromEncPkg(encpkg []byte, passwd string) (*CertKey, error) {
	blk, _ := pem.Decode(encpkg)
	if blk == nil {
		return nil, errors.New("Invalid PEM data")
	}
	if blk.Type != pkgTypeStr {
		return nil, errors.New("PEM type is not " + pkgTypeStr)
	}
	decrypted_pem, err := x509.DecryptPEMBlock(blk, []byte(passwd))
	if err != nil {
		return nil, err
	}
	key_pem, rest := pem.Decode(decrypted_pem)
	if blk == nil {
		return nil, errors.New("decrypted content is not PEM")
	}
	cert_pem, _ := pem.Decode(rest)
	if blk == nil {
		return nil, errors.New("Can't find the cert PEM")
	}
	if _, ok := supportedKeyTypes[key_pem.Type]; !ok {
		return nil, errors.New("Unsupported Key types")
	}
	if cert_pem.Type != "CERTIFICATE" {
		return nil, errors.New("Can't find certificate in decrypted PEM data")
	}
	var ck CertKey
	switch key_pem.Type {
	case "RSA PRIVATE KEY":
		rsack := new(RSACertKey)
		priv_key, err := x509.ParsePKCS1PrivateKey(key_pem.Bytes)
		if err != nil {
			return nil, err
		}
		rsack.key = priv_key
		cert, err := x509.ParseCertificates(cert_pem.Bytes)
		if err != nil {
			return nil, err
		}
		rsack.cert = *cert[0]
		ck = rsack
		return &ck, nil
	case "EC PRIVATE KEY":
		ecck := new(ECCertKey)
		priv_key, err := x509.ParseECPrivateKey(key_pem.Bytes)
		if err != nil {
			return nil, err
		}
		ecck.key = priv_key
		cert, err := x509.ParseCertificates(cert_pem.Bytes)
		if err != nil {
			return nil, err
		}
		ecck.cert = *cert[0]
		ck = ecck
		return &ck, nil
	}
	return nil, errors.New("Unussal error, you shouldn't see this")

}

func GenerateCAandEEFiles(username string, passwd []byte) error {
	//create root CA cert/key and a EE cert/key in the confdir of the username
	// CA cert are encrypted, CA and EE key are encrypted PEM
	confdir := common.GetConfDir(username)

	fi, err := os.Stat(confdir)
	if err != nil {
		return err
	}
	if !fi.IsDir() {
		return fmt.Errorf("%s is not a directory", confdir)
	}
	ca := new(RSACertKey)
	err = GenerateROOTCA(ca, username+"-ROOTCA")
	if err != nil {
		return err
	}
	ee := new(RSACertKey)
	err = CreateCertWithCA(ca, ee)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filepath.Join(confdir, "ee.cert"), ee.PemCert(), 400)
	if err != nil {
		return err
	}
	ekeys, err := ee.EncPemKey(passwd)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filepath.Join(confdir, "ee.key"), ekeys, 400)
	if err != nil {
		return err
	}
	ekeys, err = ca.EncPemKey(passwd)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filepath.Join(confdir, "ca.key"), ekeys, 400)
	if err != nil {
		return err
	}
	ecerts, err := ca.EncBaseCert(passwd)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filepath.Join(confdir, "ca.cert"), ecerts, 400)
	if err != nil {
		return err
	}
	return nil
}

func LoadManpassEE(uname string, passwd []byte) (*CertKey, error) {
	//generate a new EE cert/key with specified uname's CA
	//return encrypted CA cert, EE cert and encrypted EE key in a string
	confdir := common.GetConfDir(uname)
	fi, err := os.Stat(confdir)
	if err != nil {
		return nil, err
	}
	if !fi.IsDir() {
		return nil, fmt.Errorf("%s is not a directory", confdir)
	}
	certs, err := ioutil.ReadFile(filepath.Join(confdir, "ee.cert"))
	if err != nil {
		return nil, err
	}

	enks, err := ioutil.ReadFile(filepath.Join(confdir, "ee.key"))
	if err != nil {
		return nil, err
	}
	blk, _ := pem.Decode(enks)
	blk_cert, _ := pem.Decode(certs)
	key_der, err := x509.DecryptPEMBlock(blk, passwd)
	var ck CertKey
	switch blk.Type {
	case "RSA PRIVATE KEY":
		ee := new(RSACertKey)
		eekey, err := x509.ParsePKCS1PrivateKey(key_der)
		if err != nil {
			return nil, err
		}
		ee.key = eekey
		cert, err := x509.ParseCertificates(blk_cert.Bytes)
		if err != nil {
			return nil, err
		}
		ee.cert = *cert[0]
		ck = ee
		return &ck, nil
	case "EC PRIVATE KEY":
		ecck := new(ECCertKey)
		priv_key, err := x509.ParseECPrivateKey(key_der)
		if err != nil {
			return nil, err
		}
		ecck.key = priv_key
		cert, err := x509.ParseCertificates(blk_cert.Bytes)
		if err != nil {
			return nil, err
		}
		ecck.cert = *cert[0]
		ck = ecck
		return &ck, nil
	}
	return nil, errors.New("Unussal error, you shouldn't see this")
}

func LoadManpassCA(uname string, passwd []byte) (*CertKey, error) {
	//generate a new EE cert/key with specified uname's CA
	//return encrypted CA cert, EE cert and encrypted EE key in a string
	confdir := common.GetConfDir(uname)
	fi, err := os.Stat(confdir)
	if err != nil {
		return nil, err
	}
	if !fi.IsDir() {
		return nil, fmt.Errorf("%s is not a directory", confdir)
	}
	encacs, err := ioutil.ReadFile(filepath.Join(confdir, "ca.cert"))
	if err != nil {
		return nil, err
	}
	cacert_der, err := passcrypto.DecryptMeBase32(string(encacs), passwd)
	if err != nil {
		return nil, err
	}
	encaks, err := ioutil.ReadFile(filepath.Join(confdir, "ca.key"))
	if err != nil {
		return nil, err
	}
	blk, _ := pem.Decode(encaks)
	cakey_der, err := x509.DecryptPEMBlock(blk, passwd)
	var ck CertKey
	switch blk.Type {
	case "RSA PRIVATE KEY":
		ca := new(RSACertKey)
		cakey, err := x509.ParsePKCS1PrivateKey(cakey_der)
		if err != nil {
			return nil, err
		}
		ca.key = cakey
		cert, err := x509.ParseCertificates(cacert_der)
		if err != nil {
			return nil, err
		}
		ca.cert = *cert[0]
		ck = ca
		return &ck, nil
	case "EC PRIVATE KEY":
		ecck := new(ECCertKey)
		priv_key, err := x509.ParseECPrivateKey(cakey_der)
		if err != nil {
			return nil, err
		}
		ecck.key = priv_key
		cert, err := x509.ParseCertificates(cacert_der)
		if err != nil {
			return nil, err
		}
		ecck.cert = *cert[0]
		ck = ecck
		return &ck, nil
	}
	return nil, errors.New("Unussal error, you shouldn't see this")
}

func DumpCAEE(ca CertKey, ee CertKey, passwd []byte) (string, error) {
	rs := ""
	cas, err := ca.EncBaseCert(passwd)
	if err != nil {
		return "", err
	}
	rs += string(cas)
	rs += "\n______\n"
	rs += string(ee.PemCert())
	rs += "\n______\n"
	enckeys, err := ee.EncPemKey(passwd)
	if err != nil {
		return "", err
	}
	rs += string(enckeys)
	return rs, nil
}
