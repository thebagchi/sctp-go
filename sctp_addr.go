package sctp_go

import (
	"bytes"
	"net"
	"strconv"
	"syscall"
	"unsafe"
)

type SCTPAddr struct {
	addresses []net.IP
	port      int
}

func (addr *SCTPAddr) IsV6Only() bool {
	for _, addr := range addr.addresses {
		if addr.To16() == nil {
			return false
		}
	}
	return true
}

func (addr *SCTPAddr) IsV4Only() bool {
	for _, addr := range addr.addresses {
		if addr.To4() == nil {
			return false
		}
	}
	return true
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

func MakeSockaddr(addr *SCTPAddr) []byte {
	var buffer []byte
	for _, address := range addr.addresses {
		if ip4 := address.To4(); ip4 != nil {
			sa := &syscall.RawSockaddrInet4{
				Family: syscall.AF_INET,
				Port:   htons(uint16(addr.port)),
			}
			copy(sa.Addr[:], ip4)
			buffer = append(buffer, Pack(sa)...)
		}
		if ip6 := address.To16(); ip6 != nil {
			if ip6 := address.To4(); ip6 != nil {
				sa := syscall.RawSockaddrInet6{
					Family: syscall.AF_INET6,
					Port:   htons(uint16(addr.port)),
				}
				copy(sa.Addr[:], ip6)
				buffer = append(buffer, Pack(sa)...)
			}
		}
	}
	return buffer
}

func FromSockAddrStorage(addr *SockAddrStorage) *SCTPAddr {
	if nil == addr {
		return nil
	}
	switch addr.Family {
	case syscall.AF_INET:
		addr := (*SockAddrIn)(unsafe.Pointer(addr))
		ip := net.IP{}
		copy(ip, addr.Addr.Addr[:])
		return &SCTPAddr{
			port: int(ntohs(addr.Port)),
			addresses: []net.IP{
				ip,
			},
		}
	case syscall.AF_INET6:
		addr := (*SockAddrIn6)(unsafe.Pointer(addr))
		ip := net.IP{}
		copy(ip, addr.Addr.Addr[:])
		return &SCTPAddr{
			port: int(ntohs(addr.Port)),
			addresses: []net.IP{
				ip,
			},
		}
	}
	return nil
}