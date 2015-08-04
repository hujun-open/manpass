// api
package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"manpassd/passsql"
	"manpassd/pki"
	"manpassd/transport/tlsvr"
	"net/http"
	"strconv"
	"strings"
)

type APIError struct {
	req http.Request
	msg string
}

func (ae APIError) Error() string {
	s := fmt.Sprintf("API Transcation Error with client %[1]s", ae.req.RemoteAddr)
	return s + ae.msg
}

type ClientAPISVR struct {
	HttpSvr   http.Server
	PDB       passsql.PassDB
	Tablename string
}

func (csvr ClientAPISVR) routeClient(resp http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "POST":
		csvr.addRecord(resp, req)
	case "GET":
		if req.URL.Path == "/client/meta-id" {
			csvr.getAllMetaId(resp, req)

		} else {
			csvr.getRecord(resp, req)
		}

	case "DELETE":
		csvr.delRecord(resp, req)
	case "PUT":
		csvr.replaceAll(resp, req)
	}

}

func (csvr ClientAPISVR) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	remote_addr := strings.Split(req.RemoteAddr, ":")[0]
	if remote_addr != "127.0.0.1" {
		log.Printf("Client %s is not from local host", remote_addr)
	}
	if strings.HasPrefix(req.URL.Path, "/client") {
		csvr.routeClient(resp, req)
	} else {
		resp.WriteHeader(http.StatusNotFound)
		return
	}

}

func (csvr ClientAPISVR) parseReq(req *http.Request) (map[string]interface{}, error) {
	//parse the json in the HTTP request body, return a map
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	if string(body) == "" {
		return nil, nil
	}
	x := make(map[string]interface{})
	err = json.Unmarshal(body, &x)
	if err != nil {
		return nil, err
	}
	return x, nil

}

func (csvr ClientAPISVR) parseReqList(req *http.Request) ([]map[string]interface{}, error) {
	//parse the json in the HTTP request body, return a slice of map
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}
	x := make([]map[string]interface{}, 4096)
	err = json.Unmarshal(body, &x)
	if err != nil {
		return nil, err
	}
	return x, nil

}

func (csvr ClientAPISVR) replaceAll(resp http.ResponseWriter, req *http.Request) {
	reqlist, err := csvr.parseReqList(req)
	if err != nil {
		log.Println(err)
		resp.WriteHeader(http.StatusBadRequest)
		return
	}
	var rlist []passsql.PassRecord
	var r passsql.PassRecord
	for _, x := range reqlist {
		r.Meta = x["meta"].(string)
		r.Meta_id = x["meta_id"].(string)
		r.Uname = x["uname"].(string)
		r.Pass = x["pass"].(string)
		r.Pass_rev = int(x["pass_rev"].(float64))
		rlist = append(rlist, r)
	}
	err = csvr.PDB.ReplaceAll(csvr.Tablename, rlist)
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}
	resp.WriteHeader(http.StatusOK)
	return
}

func (csvr ClientAPISVR) addRecord(resp http.ResponseWriter, req *http.Request) {
	x, err := csvr.parseReq(req)
	if err != nil {
		resp.WriteHeader(http.StatusBadRequest)
		return
	}
	var r passsql.PassRecord
	r.Meta = x["meta"].(string)
	r.Meta_id = x["meta_id"].(string)
	r.Uname = x["uname"].(string)
	r.Pass = x["pass"].(string)
	err = csvr.PDB.Insert(csvr.Tablename, r)
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}
	resp.WriteHeader(http.StatusCreated)
	return

}

func (csvr ClientAPISVR) delRecord(resp http.ResponseWriter, req *http.Request) {
	x, err := csvr.parseReq(req)
	if err != nil {
		resp.WriteHeader(http.StatusBadRequest)
		return
	}
	if _, ok := x["pass_rev"]; ok {
		err := csvr.PDB.RemovePassForRev(csvr.Tablename, x["meta_id"].(string), int(x["pass_rev"].(float64)))
		if err != nil {
			resp.WriteHeader(http.StatusInternalServerError)
			return
		}
		resp.WriteHeader(http.StatusOK)

	} else {
		err := csvr.PDB.RemovePass(csvr.Tablename, x["meta_id"].(string))
		if err != nil {
			resp.WriteHeader(http.StatusInternalServerError)
			return
		}
		resp.WriteHeader(http.StatusOK)
	}

	return

}

func (csvr ClientAPISVR) getAllMetaId(resp http.ResponseWriter, req *http.Request) {
	rlist, err := csvr.PDB.GetAllMetaId(csvr.Tablename)
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}
	if len(rlist) == 0 {
		resp.WriteHeader(http.StatusNoContent)
		return

	}
	js, err := json.Marshal(rlist)
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}
	resp.WriteHeader(http.StatusOK)
	fmt.Fprintf(resp, string(js))
	return

}

func (csvr ClientAPISVR) getRecord(resp http.ResponseWriter, req *http.Request) {
	x, err := csvr.parseReq(req)
	if err != nil {
		resp.WriteHeader(http.StatusBadRequest)
		return
	}
	if x["meta_id"] == "__ALLLATESTPASS__" {
		r, err := csvr.PDB.GetAllLatest(csvr.Tablename)
		if err != nil {
			resp.WriteHeader(http.StatusInternalServerError)
			return
		}
		if len(r) == 0 {
			resp.WriteHeader(http.StatusNoContent)
			return
		}
		js, err := json.Marshal(r)
		if err != nil {
			resp.WriteHeader(http.StatusInternalServerError)
			return
		}
		resp.WriteHeader(http.StatusOK)
		fmt.Fprintf(resp, string(js))
		return
	}

	if _, ok := x["pass_rev"]; ok {
		r, err := csvr.PDB.GetRecord(csvr.Tablename, x["meta_id"].(string), int(x["pass_rev"].(float64)))
		if err != nil {
			resp.WriteHeader(http.StatusInternalServerError)
			return
		}
		if r == nil {
			resp.WriteHeader(http.StatusNoContent)
			return
		}
		js, err := json.Marshal(*r)
		if err != nil {
			resp.WriteHeader(http.StatusInternalServerError)
			return
		}
		resp.WriteHeader(http.StatusOK)
		fmt.Fprintf(resp, string(js))

	} else if _, ok = x["meta_id"]; ok {
		r, err := csvr.PDB.GetAllRevForMetaId(csvr.Tablename, x["meta_id"].(string))
		if err != nil {
			resp.WriteHeader(http.StatusInternalServerError)
			return
		}
		if len(r) == 0 {
			resp.WriteHeader(http.StatusNoContent)
			return
		}
		js, err := json.Marshal(r)
		if err != nil {
			resp.WriteHeader(http.StatusInternalServerError)
			return
		}
		resp.WriteHeader(http.StatusOK)
		fmt.Fprintf(resp, string(js))
	} else {
		r, err := csvr.PDB.GetAll(csvr.Tablename)
		if err != nil {
			log.Println(err)
			resp.WriteHeader(http.StatusInternalServerError)
			return
		}
		if len(r) == 0 {
			resp.WriteHeader(http.StatusNoContent)
			return
		}
		js, err := json.Marshal(r)
		if err != nil {
			log.Println(err)
			resp.WriteHeader(http.StatusInternalServerError)
			return
		}
		resp.WriteHeader(http.StatusOK)
		fmt.Fprintf(resp, string(js))
	}

	return
}

func (csvr ClientAPISVR) Serve() {
	go csvr.HttpSvr.ListenAndServe()
}

func NewClientAPISVR(ipaddr string, port int, uname string, upass []byte, pdb passsql.PassDB, tablename string) (*ClientAPISVR, error) {
	csvr := new(ClientAPISVR)
	tls_svr, err := tlsvr.GetServer(ipaddr, port, uname, upass)
	if err != nil {
		return nil, err
	}
	csvr.HttpSvr = tls_svr
	csvr.HttpSvr.Handler = csvr
	csvr.PDB = pdb
	csvr.Tablename = tablename
	return csvr, nil

}

type ProvisionSVR struct {
	Uname  string
	passwd []byte
	Svr    http.Server
	ca     *pki.CertKey
}

func NewProvisionSVR(ipaddr string, port int, uname string, passwd []byte) (*ProvisionSVR, error) {
	psvr := new(ProvisionSVR)
	psvr.Svr.Addr = ipaddr + ":" + strconv.Itoa(port)
	psvr.Svr.Handler = psvr
	psvr.Uname = uname
	psvr.passwd = passwd
	ca, err := pki.LoadManpassCA(uname, passwd)
	if err != nil {
		return nil, err
	}
	psvr.ca = ca
	return psvr, nil
}
func (psvr ProvisionSVR) Serve() {
	go psvr.Svr.ListenAndServe()
}
func (psvr ProvisionSVR) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	remote_addr := strings.Split(req.RemoteAddr, ":")[0]
	if remote_addr != "127.0.0.1" {
		log.Printf("Client %s is not from local host", remote_addr)
	}
	if req.Method == "GET" {
		ee := new(pki.RSACertKey)
		err := pki.CreateCertWithCA(*(psvr.ca), ee)
		if err != nil {
			resp.WriteHeader(http.StatusInternalServerError)
			return
		}
		rs, err := pki.DumpCAEE(*(psvr.ca), ee, psvr.passwd)
		if err != nil {
			resp.WriteHeader(http.StatusInternalServerError)
			return
		}
		resp.WriteHeader(http.StatusOK)
		fmt.Fprintf(resp, rs)
	}

}
