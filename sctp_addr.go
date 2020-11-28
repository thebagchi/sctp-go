package sctp_go

import (
	"bytes"
	"net"
	"strconv"
)

type SCTPAddr struct {
	addresses []net.IP
	port      int
}

func (addr *SCTPAddr) String() string {
	var b bytes.Buffer
	for n, address := range addr.addresses {
		if n > 0 {
			b.WriteRune('/')
		}
		if address.To4() != nil {
			b.WriteString(address.String())
		} else if address.To16() != nil {
			b.WriteRune('[')
			b.WriteString(address.String())
			b.WriteRune(']')
		}
	}
	b.WriteRune(':')
	b.WriteString(strconv.Itoa(addr.port))
	return b.String()
}

func (addr *SCTPAddr) Network() string {
	return "sctp"
}

func (addr *SCTPAddr) Addr() net.Addr {
	if addr == nil {
		return nil
	}
	return addr
}