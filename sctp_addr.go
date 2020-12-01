package sctp_go

import (
	"bytes"
	"net"
	"strconv"
	"strings"
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
			continue
		}
		if ip6 := address.To16(); ip6 != nil {
			sa := syscall.RawSockaddrInet6{
				Family: syscall.AF_INET6,
				Port:   htons(uint16(addr.port)),
			}
			copy(sa.Addr[:], ip6)
			buffer = append(buffer, Pack(sa)...)
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

func FromSCTPGetAddrs(addr *SCTPGetAddrs) *SCTPAddr {
	if nil != addr {
		address := &SCTPAddr{
			addresses: make([]net.IP, addr.Num),
		}
		ptr := unsafe.Pointer(addr.Addr)
		for i := uint32(0); i < addr.Num;  i++ {
			addr := (*SockAddr)(unsafe.Pointer(addr))
			size := uintptr(0)
			switch addr.Family {
			case syscall.AF_INET:
				addr := (*SockAddrIn)(unsafe.Pointer(addr))
				address.port = int(ntohs(addr.Port))
				copy(address.addresses[i], addr.Addr.Addr[:])
				size = SockAddrInSize
			case syscall.AF_INET6:
				addr := (*SockAddrIn6)(unsafe.Pointer(addr))
				address.port = int(ntohs(addr.Port))
				copy(address.addresses[i], addr.Addr.Addr[:])
				size = SockAddrIn6Size
			default:
				return nil
			}
			ptr = unsafe.Pointer(uintptr(ptr) + size * uintptr(i))
		}
	}
	return nil
}

func MakeSCTPAddr(network, addr string) (*SCTPAddr, error) {
	switch network {
	case "", "sctp":
		network = "sctp"
	case "sctp4":
		network = "sctp"
	case "sctp6":
		network = "sctp"
	default:
		return nil, net.UnknownNetworkError(network)
	}

	if strings.LastIndex(addr, ":") < 0 {
		return nil, &net.AddrError{
			Err:  "missing port in address",
			Addr: addr,
		}
	}

	var (
		index           = strings.LastIndex(addr, ":")
		addrs           = strings.Split(addr[:index], "/")
		port            = 0
		err       error = nil
		addresses       = make([]net.IP, 0)
	)

	if index < 0 || index == len(addr) {
		return nil, &net.AddrError{
			Err:  "missing port in address",
			Addr: addr,
		}
	} else if port, err = net.LookupPort(network, addr[index+1:]); err != nil {
		return nil, &net.AddrError{
			Err:  "missing port in address: " + err.Error(),
			Addr: addr,
		}
	}

	for _, addr := range addrs {
		if len(addr) == 0 {
			if network == "sctp" {
				addresses = append(addresses, net.IPv4zero)
				continue
			}
			if network == "sctp6" {
				addresses = append(addresses, net.IPv6zero)
			}
		} else {
			address := net.ParseIP(addr)
			if network == "sctp6" && address.To16() != nil {
				addresses = append(addresses, address)
				continue
			}
			if network == "sctp" && address.To4() != nil {
				addresses = append(addresses, address)
			}
		}
	}

	if len(addresses) == 0 {
		return nil, net.InvalidAddrError(addr)
	}
	address := &SCTPAddr{
		addresses: addresses,
		port:      port,
	}
	return address, nil
}