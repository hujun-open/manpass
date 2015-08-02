/*
This package provide crypto services, most high level entry functions are:
- EncryptMe: use scrypt to generate encrypt key for each encryption, and append salt to the result cipher; input clear text is a slice of byte
- DecryptMe
- EncryptWithoutSalt: doesn't use scrypt, expect user to re-use the scrypt generated key;much faster but less secure
- DecryptWithoutSalt

The underlying crypto is nacl/secretbox
*/
package passcrypto

import (
	"crypto/rand"
	"encoding/base32"
	"errors"
	"golang.org/x/crypto/nacl/secretbox"
	"golang.org/x/crypto/scrypt"
	"io"
)

var (
	ErrEncrypt = errors.New("secret: encryption failed")
	ErrDecrypt = errors.New("secret: decryption failed")
)

const (
	KeySize   = 32
	NonceSize = 24
	SaltSize  = 16
	Scrypt_N  = 1024 //1048576 is for file encryption, 16384 is for interactive logins
	Scrypt_r  = 1
	Scrypt_p  = 1
)

func ZeroMe(myr []byte) {
	//zero a byte slice
	for k, _ := range myr {
		myr[k] = 0
	}
	return
}

func GenerateNonce() (*[NonceSize]byte, error) {
	nonce := new([NonceSize]byte)
	_, err := io.ReadFull(rand.Reader, nonce[:])
	if err != nil {
		return nil, err
	}
	return nonce, nil
}

func GenerateSalt() (*[SaltSize]byte, error) {
	salt := new([SaltSize]byte)
	_, err := io.ReadFull(rand.Reader, salt[:])
	if err != nil {
		return nil, err
	}
	return salt, nil
}

func GenerateEncKey(passwd []byte, salt *[SaltSize]byte) (*[KeySize]byte, error) {
	var rkey = new([KeySize]byte)
	skey, err := scrypt.Key(passwd, salt[:], Scrypt_N, Scrypt_r, Scrypt_p, KeySize)
	if err != nil {

		return nil, err
	}
	copy(rkey[:], skey)
	ZeroMe(skey)
	return rkey, nil
}

func EncryptWithoutSalt(msg []byte, skey *[KeySize]byte) ([]byte, error) {
	nonce, err := GenerateNonce()
	if err != nil {
		return nil, ErrEncrypt
	}
	out := make([]byte, len(msg)+NonceSize)
	out = secretbox.Seal(nonce[:], []byte(msg), nonce, skey)
	return out, nil

}

func EncryptMe(msg []byte, inpasswd []byte) ([]byte, error) {
	salt, err := GenerateSalt()
	if err != nil {
		return nil, ErrEncrypt
	}
	skey, err := GenerateEncKey(inpasswd, salt)
	if err != nil {
		return nil, ErrEncrypt
	}
	rstr, err := EncryptWithoutSalt(msg, skey)
	if err != nil {
		return nil, ErrEncrypt
	}
	return append(salt[:], rstr...), nil
}

func EncryptMeBase32(msg []byte, inpasswd []byte) (string, error) {
	cipher, err := EncryptMe(msg, inpasswd)
	if err != nil {
		return "", err
	}
	return base32.StdEncoding.EncodeToString(cipher), nil
}

func DecryptWithoutSalt(cipher []byte, skey *[KeySize]byte) ([]byte, error) {
	var nonce = [NonceSize]byte{}
	copy(nonce[:], cipher[:NonceSize])
	out, ok := secretbox.Open(nil, cipher[NonceSize:], &nonce, skey)
	if !ok {
		return nil, ErrDecrypt
	}
	return out, nil
}

func DecryptMe(cipher []byte, inpasswd []byte) ([]byte, error) {
	var salt = [SaltSize]byte{}
	copy(salt[:], cipher[:SaltSize])
	skey, err := GenerateEncKey(inpasswd, &salt)
	if err != nil {
		return nil, ErrDecrypt
	}
	ostr, err := DecryptWithoutSalt(cipher[SaltSize:], skey)
	if err != nil {
		return nil, ErrDecrypt
	}
	return ostr, nil
}

func DecryptMeBase32(cipher_b32 string, inpasswd []byte) ([]byte, error) {
	cipher, err := base32.StdEncoding.DecodeString(cipher_b32)
	if err != nil {
		return nil, err
	}
	return DecryptMe(cipher, inpasswd)

}
