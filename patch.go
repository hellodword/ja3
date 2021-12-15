package ja3

import (
	"bufio"
	"fmt"
	tls "github.com/refraction-networking/utls"
	"golang.org/x/net/http2"
	"golang.org/x/net/proxy"
	"net"
	"net/http"
)

func UTLSConnViaProxy(hostname string, addr string, dialer proxy.Dialer, id tls.ClientHelloID) (*tls.UConn, error) {
	tcpConn, err := dialer.Dial("tcp", addr)
	if err != nil {
		err = fmt.Errorf("utls dial error: %s", err.Error())
		return nil, err
	}

	config := tls.Config{ServerName: hostname}
	uTlsConn := tls.UClient(tcpConn, &config, id)
	err = uTlsConn.Handshake()
	if err != nil {
		err = fmt.Errorf("utls handshake error: %s", err.Error())
		return nil, err
	} else {
		return uTlsConn, nil
	}
}

func HTTPOverConn(conn net.Conn, alpn string, req *http.Request) (*http.Response, error) {

	switch alpn {
	case "h2":
		req.Proto = "HTTP/2.0"
		req.ProtoMajor = 2
		req.ProtoMinor = 0

		tr := http2.Transport{}
		cConn, err := tr.NewClientConn(conn)
		if err != nil {
			return nil, err
		}
		return cConn.RoundTrip(req)
	case "http/1.1", "":
		req.Proto = "HTTP/1.1"
		req.ProtoMajor = 1
		req.ProtoMinor = 1

		err := req.Write(conn)
		if err != nil {
			return nil, err
		}
		return http.ReadResponse(bufio.NewReader(conn), req)
	default:
		return nil, fmt.Errorf("unsupported ALPN: %v", alpn)
	}
}
