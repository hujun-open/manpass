// main
package main

import (
	"fmt"
	"log"
	//"manpass/transport/tlsvr"
	"flag"
	"golang.org/x/crypto/ssh/terminal"
	"manpassd/api"
	"manpassd/common"
	"manpassd/passsql"
	"manpassd/pki"
	"manpassd/transport/tlsvr"
	"os"
	"path/filepath"
)

func main() {
	version_str := "Manpass ver1.0 - server daemon"
	fmt.Println(version_str)
	var uname = flag.String("username", "", "specify the username")
	var svr_port = flag.Int("svrport", 9000, "specify the server listening port")
	var svr_ip = flag.String("svrip", "127.0.0.1", "specify the server listening IP address")
	var pipe_pass = flag.Bool("pipepass", false, "enable to use pipe to pass password")
	var createuser = flag.Bool("create", false, "create a new user")
	flag.Parse()
	if *uname == "" || !common.GoodUname(*uname) {
		fmt.Println("Error: Missing username or invalid username")
		flag.PrintDefaults()
		return
	}

	confDir := common.GetConfDir(*uname)
	fi, err := os.Stat(confDir)
	var upass []byte
	if (err != nil || !fi.IsDir()) && *createuser {
		//if user directory doesn't exist
		if *pipe_pass == false {
			fmt.Printf("user %s does not exisit! Creating a new user:%s\n", *uname, *uname)
			upass = common.InputNewPassword(*uname)
		} else {
			var pass_str string
			_, err := fmt.Scan(&pass_str)
			if err != nil {
				log.Println(err)
				log.Fatal("invalid input")
			}
			upass = []byte(pass_str)
		}
		os.Remove(confDir)
		fmt.Println("creating needed files...\n")
		err = os.MkdirAll(confDir, 700)
		if err != nil {
			log.Fatalf("Failed to create directory:%s\n", confDir)
		}
		err = pki.GenerateCAandEEFiles(*uname, upass)
		if err != nil {
			log.Fatalf("Failed to generate CA and EE cert/keys:%s\n", err)
		}
		err = passsql.InitPassDB(*uname)
		if err != nil {
			log.Fatalf("Failed to init PassDB:%s\n", err)
		}
		fmt.Println("new user created.\n")
	} else {
		//load exisiting user directory
		if *pipe_pass == false {
			fmt.Printf("Input password for user %s:", *uname)
			upass, err = terminal.ReadPassword(int(os.Stdin.Fd()))
			if err != nil {
				log.Println(err)
				log.Fatal("invalid input")

			}
		} else {
			var pass_str string
			_, err := fmt.Scan(&pass_str)
			if err != nil {
				log.Println(err)
				log.Fatal("invalid input")
			}
			upass = []byte(pass_str)

		}

	}
	log.Println("Starting server...")
	dbfile := filepath.Join(confDir, *uname+".db")
	passdb, err := passsql.LoadDB(dbfile)
	if err != nil {
		log.Fatalf("Failed to load db %s, %s", dbfile, err)
	}
	svr, err := api.NewClientAPISVR(*svr_ip, *svr_port, *uname, upass, *passdb, *uname)
	if err != nil {
		log.Fatal(err)
	}
	err = tlsvr.ListenWithConfig(svr.HttpSvr)
	if err != nil {
		log.Fatal(err)
	}
	//	httpsvr, err := api.NewProvisionSVR(*svr_ip, *svr_port+1, *uname, upass)
	//	httpsvr.Serve()
	var c chan int
	<-c
}
