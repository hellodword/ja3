package ja3

import (
	"encoding/json"
	"fmt"
	tls "github.com/refraction-networking/utls"
	"net"
	"net/http"
	"time"
)

var (
	ProxyAddr = "127.0.0.1:1081"
)

func ExampleHTTPOverConn() {
	req, err := http.NewRequest(http.MethodGet, "https://ja3er.com/json", nil)
	if err != nil {
		panic(err)
	}

	if req.URL.Scheme != "https" {
		panic(req.URL.Scheme)
	}

	addr, err := net.LookupHost(req.URL.Hostname())
	if err != nil || len(addr) == 0 {
		panic(err)
	}

	dialer := ProxyDirect{
		Dialer: net.Dialer{
			Timeout: time.Second * 5,
		},
	}
	proxyDial, err := ProxyHTTP("tcp", ProxyAddr, nil, time.Second*5, dialer)
	if err != nil {
		panic(err)
	}

	conn, err := UTLSConnViaProxy(req.URL.Hostname(), addr[0]+":443", proxyDial, tls.HelloChrome_62)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	response, err := HTTPOverConn(conn, conn.HandshakeState.ServerHello.AlpnProtocol, req)
	if err != nil {
		panic(err)
	}

	defer response.Body.Close()

	var ja3 struct {
		Hash string `json:"ja3_hash"`
	}

	err = json.NewDecoder(response.Body).Decode(&ja3)
	if err != nil {
		panic(err)
	}

	fmt.Println(ja3.Hash)
	// Output: bc6c386f480ee97b9d9e52d472b772d8
}
