// common
package common

import (
	"bytes"
	"fmt"
	"golang.org/x/crypto/ssh/terminal"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
)

func GetConfDir(uname string) string {
	var defDir string
	switch runtime.GOOS {
	case "windows":
		defDir = filepath.Join(os.Getenv("APPDATA"), "manpass", uname)
	case "linux", "darwin":
		defDir = filepath.Join(os.Getenv("HOME"), ".manpass", uname)
	}
	redirectfilename := filepath.Join(defDir, "redirection.conf")
	redir, err := ioutil.ReadFile(redirectfilename)
	if err != nil {
		return defDir
	} else {
		return string(redir)
	}
	return ""
}

func GoodUname(uname string) bool {
	if len(uname) < 3 {
		return false
	} else {
		return true
	}
}

func GoodPass(pass []byte) bool {
	if len(pass) < 6 {
		return false
	} else {
		return true
	}
}

func InputNewUser() (string, []byte, error) {
	var uname string
	for true {
		fmt.Print("Input your username:")
		_, err := fmt.Scan(&uname)
		if err != nil {
			fmt.Println("\nInvalid input!")
			continue
		}
		if GoodUname(uname) {
			break
		} else {
			fmt.Println("Invalid username!")
		}
	}
	os.Stdin.Read(make([]byte, 1024))
	for true {
		fmt.Print("\nInput your password:")
		upass, err := terminal.ReadPassword(int(os.Stdin.Fd()))
		if err != nil {
			fmt.Println("\nInvalid input!")
			continue
		}
		fmt.Print("\nType the password again:")
		pass2, err := terminal.ReadPassword(int(os.Stdin.Fd()))
		if err != nil {
			fmt.Println("\nInvalid input!")
			continue
		}
		if !bytes.Equal(upass, pass2) {
			fmt.Println("\nTwo typing are not equal!")
			continue
		}
		if GoodPass(upass) {
			fmt.Println("\n")
			return uname, upass, nil
		} else {
			fmt.Println("\nNot a good password!")
		}
	}

	return "", nil, nil
}

func InputNewPassword(uname string) []byte {
	fmt.Printf("creating password for user %s\n", uname)
	for true {
		fmt.Print("\nInput your password:")
		upass, err := terminal.ReadPassword(int(os.Stdin.Fd()))
		if err != nil {
			fmt.Println("\nInvalid input!")
			continue
		}
		fmt.Print("\nType the password again:")
		pass2, err := terminal.ReadPassword(int(os.Stdin.Fd()))
		if err != nil {
			fmt.Println("\nInvalid input!")
			continue
		}
		if !bytes.Equal(upass, pass2) {
			fmt.Println("\nTwo typing are not equal!")
			continue
		}
		if GoodPass(upass) {
			fmt.Println("\n")
			return upass
		} else {
			fmt.Println("\nNot a good password!")
		}
	}
	fmt.Println("\n")
	return nil
}
