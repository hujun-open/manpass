// tls-svr
package tlsvr

import (
	"crypto/tls"
	//"fmt"
	"manpassd/transport"
	"net"
	"net/http"
	"strconv"
	"time"
)

func GetServer(ipaddr string, port int, uname string, upass []byte) (svr http.Server, err error) {
	ca_f, ee_f, key_f, err := transport.LoadCertsKeys(uname, upass)
	if err != nil {
		return svr, err
	}
	config, err := transport.GetConfig(ca_f, ee_f, key_f)

	if err != nil {
		return svr, err
	}
	config.ServerName = "agent-1.hj.com"
	svr.Addr = ipaddr + ":" + strconv.Itoa(port)
	svr.TLSConfig = config

	return
}

//following is copied from net/http/server.go
type tcpKeepAliveListener struct {
	*net.TCPListener
}

func (ln tcpKeepAliveListener) Accept() (c net.Conn, err error) {
	tc, err := ln.AcceptTCP()
	if err != nil {
		return
	}
	tc.SetKeepAlive(true)
	tc.SetKeepAlivePeriod(3 * time.Minute)
	return tc, nil
}

//end of copy

func ListenWithConfig(srv http.Server) error {
	//this is a modified version of http.Server.ListenAndServeTLS
	//using http.Server.TLSConfig without specify cert/key as arguments
	addr := srv.Addr
	if addr == "" {
		addr = ":https"
	}
	config := srv.TLSConfig
	var err error
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	tlsListener := tls.NewListener(tcpKeepAliveListener{ln.(*net.TCPListener)}, config)
	go srv.Serve(tlsListener)
	return nil
}
