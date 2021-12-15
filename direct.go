// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ja3

import (
	"context"
	"net"
)

type ProxyDirect struct {
	Dialer net.Dialer
}

// Dial directly invokes net.Dial with the supplied parameters.
func (d ProxyDirect) Dial(network, addr string) (net.Conn, error) {
	return d.Dialer.Dial(network, addr)
}

// DialContext instantiates a net.Dialer and invokes its DialContext receiver with the supplied parameters.
func (d ProxyDirect) DialContext(ctx context.Context, network, addr string) (net.Conn, error) {
	return d.Dialer.DialContext(ctx, network, addr)
}
