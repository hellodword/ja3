// meek-client is the client transport plugin for the meek pluggable transport.
//
// Sample usage in torrc:
// 	Bridge meek 0.0.2.0:1 url=https://forbidden.example/ front=allowed.example
// 	ClientTransportPlugin meek exec ./meek-client
// The transport ignores the bridge address 0.0.2.0:1 and instead connects to
// the URL given by url=. When front= is given, the domain in the URL is
// replaced by the front domain for the purpose of the DNS lookup, TCP
// connection, and TLS SNI, but the HTTP Host header in the request will be the
// one in url=.
//
// Most user configuration can happen either through SOCKS args (i.e., args on a
// Bridge line) or through command line options. SOCKS args take precedence
// per-connection over command line options. For example, this configuration
// using SOCKS args:
// 	Bridge meek 0.0.2.0:1 url=https://forbidden.example/ front=allowed.example
// 	ClientTransportPlugin meek exec ./meek-client
// is the same as this one using command line options:
// 	Bridge meek 0.0.2.0:1
// 	ClientTransportPlugin meek exec ./meek-client --url=https://forbidden.example/ --front=allowed.example
// The command-line configuration interface is for compatibility with tor 0.2.4
// and older, which doesn't support parameters on Bridge lines.
//
// The --helper option prevents this program from doing any network operations
// itself. Rather, it will send all requests through a browser extension that
// makes HTTP requests.
package ja3

import (
	"net/http"
	"time"
)

const (
	ptMethodName = "meek"
	// A session ID is a randomly generated string that identifies a
	// long-lived session. We split a TCP stream across multiple HTTP
	// requests, and those with the same session ID belong to the same
	// stream.
	sessionIDLength = 8
	// The size of the largest chunk of data we will read from the SOCKS
	// port before forwarding it in a request, and the maximum size of a
	// body we are willing to handle in a reply.
	maxPayloadLength = 0x10000
	// We must poll the server to see if it has anything to send; there is
	// no way for the server to push data back to us until we send an HTTP
	// request. When a timer expires, we send a request even if it has an
	// empty body. The interval starts at this value and then grows.
	initPollInterval = 100 * time.Millisecond
	// Maximum polling interval.
	maxPollInterval = 5 * time.Second
	// Geometric increase in the polling interval each time we fail to read
	// data.
	pollIntervalMultiplier = 1.5
	// Try an HTTP roundtrip at most this many times.
	maxTries = 10
	// Wait this long between retries.
	retryDelay = 30 * time.Second
	// Safety limits on interaction with the HTTP helper.
	maxHelperResponseLength = 10000000
	helperReadTimeout       = 60 * time.Second
	helperWriteTimeout      = 2 * time.Second
)

// We use this RoundTripper to make all our requests when neither --helper nor
// utls is in effect. We use the defaults, except we take control of the Proxy
// setting (notably, disabling the default ProxyFromEnvironment).
var httpRoundTripper *http.Transport = http.DefaultTransport.(*http.Transport).Clone()
