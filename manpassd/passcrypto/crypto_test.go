// crypto_test
package passcrypto

import (
	"testing"
)

func TestCrypto(t *testing.T) {
	cases := []struct {
		clear string
		pass  string
	}{
		{"Hello, world", "Zid21ukj&^31A(#"},
		{"Hello, 世界", "dis@34SAX,87d"},
		{"", "dki#@21@99,kkK"},
	}
	for _, c := range cases {
		cstr, err := EncryptMe([]byte(c.clear), &c.pass)
		if err != nil {
			t.Errorf("Error encrypting for %q", c.clear)
		}
		t.Log(cstr)
		ostr, err := DecryptMe(cstr, &c.pass)
		if err != nil {
			t.Errorf("Error decrypting for %q", c.clear)
		}
		if string(ostr) != c.clear {
			t.Errorf("decrypted text is not equal to clear text: %q", c.clear)
		}
		b32, err := EncryptMeBase32([]byte(c.clear), &c.pass)
		if err != nil {
			t.Errorf("error encrypting B32 for %q ", c.clear)
		}
		t.Log(b32)
	}
}
