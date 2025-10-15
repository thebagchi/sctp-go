package sctp_go

import (
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

// Address returns the address part of the SCTPAddr without the port.
// Multiple addresses are separated by '/'.
func (addr *SCTPAddr) Address() string {
	if len(addr.addresses) == 0 {
		return ""
	}
	var b strings.Builder
	b.Grow(64)
	for n, address := range addr.addresses {
		if n > 0 {
			b.WriteByte('/')
		}
		if ip4 := address.To4(); ip4 != nil {
			b.WriteString(ip4.String())
			continue
		}
		if ip6 := address.To16(); ip6 != nil {
			b.WriteByte('[')
			b.WriteString(address.String())
			b.WriteByte(']')
		}
	}
	return b.String()
}

// Port returns the port number of the SCTPAddr.
func (addr *SCTPAddr) Port() int {
	return addr.port
}

// IsV6Only reports whether the SCTPAddr contains only IPv6 addresses.
func (addr *SCTPAddr) IsV6Only() bool {
	for _, addr := range addr.addresses {
		if addr.To16() == nil {
			return false
		}
	}
	return true
}

// IsV4Only reports whether the SCTPAddr contains only IPv4 addresses.
func (addr *SCTPAddr) IsV4Only() bool {
	for _, addr := range addr.addresses {
		if addr.To4() == nil {
			return false
		}
	}
	return true
}

// String returns the string representation of the SCTPAddr in the format "address:port".
// Multiple addresses are separated by '/', and IPv6 addresses are enclosed in brackets.
func (addr *SCTPAddr) String() string {
	var b strings.Builder
	b.Grow(128)
	for n, address := range addr.addresses {
		if n > 0 {
			b.WriteByte('/')
		}
		if ip4 := address.To4(); ip4 != nil {
			b.WriteString(ip4.String())
		}
		if ip6 := address.To16(); ip6 != nil {
			b.WriteByte('[')
			b.WriteString(address.String())
			b.WriteByte(']')
		}
	}
	b.WriteByte(':')
	b.WriteString(strconv.Itoa(addr.port))
	return b.String()
}

// Network returns the network type, which is always "sctp".
func (addr *SCTPAddr) Network() string {
	return "sctp"
}

// Addr returns the SCTPAddr itself, implementing the net.Addr interface.
func (addr *SCTPAddr) Addr() net.Addr {
	if addr == nil {
		return nil
	}
	return addr
}

// MakeSockaddr converts an SCTPAddr to a byte slice containing socket address structures
// suitable for use with SCTP system calls. It handles both IPv4 and IPv6 addresses.
func MakeSockaddr(addr *SCTPAddr) []byte {
	if len(addr.addresses) == 0 {
		return nil
	}
	// Pre-calculate buffer capacity and cache port conversion
	port := htons(uint16(addr.port))
	capacity := 0
	for _, address := range addr.addresses {
		if address.To4() != nil {
			capacity += SockAddrInSize
		} else {
			capacity += SockAddrIn6Size
		}
	}
	buffer := make([]byte, 0, capacity)
	for _, address := range addr.addresses {
		if ip4 := address.To4(); ip4 != nil {
			/*
				sa := &syscall.RawSockaddrInet4{
					Family: syscall.AF_INET,
					Port: port,
				}
				copy(sa.Addr[:], ip4)
				buffer = append(buffer, Pack(sa)...)
			*/
			// IPv4: Family(2) + Port(2) + Addr(4) + Zero(8) = 16 bytes
			var sockaddr [16]byte
			endian.PutUint16(sockaddr[0:2], syscall.AF_INET)
			endian.PutUint16(sockaddr[2:4], port)
			copy(sockaddr[4:8], ip4)
			// sockaddr[8:16] remains zero
			buffer = append(buffer, sockaddr[:]...)
			continue
		}
		if ip6 := address.To16(); ip6 != nil {
			/*
				sa := syscall.RawSockaddrInet6{
					Family: syscall.AF_INET6,
					Port:   port,
				}
				copy(sa.Addr[:], address.To16())
				buffer = append(buffer, Pack(sa)...)
			*/
			// IPv6: Family(2) + Port(2) + FlowInfo(4) + Addr(16) + ScopeId(4) = 28 bytes
			var sockaddr [28]byte
			endian.PutUint16(sockaddr[0:2], syscall.AF_INET6)
			endian.PutUint16(sockaddr[2:4], port)
			// sockaddr[4:8] remains zero (FlowInfo)
			copy(sockaddr[8:24], address.To16())
			// sockaddr[24:28] remains zero (ScopeId)
			buffer = append(buffer, sockaddr[:]...)
		}
	}
	return buffer
}

// FromSockAddrStorage converts a socket address storage structure to an SCTPAddr.
// It supports both IPv4 (AF_INET) and IPv6 (AF_INET6) address families.
func FromSockAddrStorage(addr *SockAddrStorage) *SCTPAddr {
	if addr != nil {
		switch addr.Family {
		case syscall.AF_INET:
			sa := (*SockAddrIn)(unsafe.Pointer(addr))
			return &SCTPAddr{
				port: int(ntohs(sa.Port)),
				addresses: []net.IP{
					sa.Addr.Addr[:],
				},
			}
		case syscall.AF_INET6:
			sa := (*SockAddrIn6)(unsafe.Pointer(addr))
			return &SCTPAddr{
				port: int(ntohs(sa.Port)),
				addresses: []net.IP{
					sa.Addr.Addr[:],
				},
			}
		}
	}
	return nil
}

// FromSCTPGetAddrs converts a SCTPGetAddrs structure to an SCTPAddr.
// It handles both IPv4 and IPv6 addresses, extracting the port from the first address.
func FromSCTPGetAddrs(addr *SCTPGetAddrs) *SCTPAddr {
	if addr != nil && addr.Num > 0 {
		address := &SCTPAddr{
			addresses: make([]net.IP, addr.Num),
		}
		ptr := unsafe.Add(unsafe.Pointer(addr), SCTPGetAddrsSize)
		for i := uint32(0); i < addr.Num; i++ {
			sockaddr := (*SockAddr)(unsafe.Pointer(ptr))
			switch sockaddr.Family {
			case syscall.AF_INET:
				sa := (*SockAddrIn)(unsafe.Pointer(ptr))
				if i == 0 {
					address.port = int(ntohs(sa.Port))
				}
				address.addresses[i] = sa.Addr.Addr[:]
				ptr = unsafe.Add(ptr, SockAddrInSize)
			case syscall.AF_INET6:
				sa := (*SockAddrIn6)(unsafe.Pointer(ptr))
				if i == 0 {
					address.port = int(ntohs(sa.Port))
				}
				address.addresses[i] = sa.Addr.Addr[:]
				ptr = unsafe.Add(ptr, SockAddrIn6Size)
			default:
				return nil
			}
		}
		return address
	}
	return nil
}

// MakeSCTPAddr parses a network address string and creates an SCTPAddr.
// It supports "sctp", "sctp4", and "sctp6" networks, parsing comma-separated IP addresses.
func MakeSCTPAddr(network, addr string) (*SCTPAddr, error) {
	// Normalize network
	switch network {
	case "", "sctp":
		network = "sctp"
	case "sctp4":
		network = "sctp4"
	case "sctp6":
		network = "sctp6"
	default:
		return nil, net.UnknownNetworkError(network)
	}

	// Find port separator
	index := strings.LastIndex(addr, ":")
	if index <= 0 || index == len(addr)-1 {
		return nil, &net.AddrError{
			Err:  "missing port in address",
			Addr: addr,
		}
	}

	// Parse port
	port, err := net.LookupPort(network, addr[index+1:])
	if err != nil {
		return nil, &net.AddrError{
			Err:  "missing port in address: " + err.Error(),
			Addr: addr,
		}
	}

	// Split addresses
	addrs := strings.Split(addr[:index], "/")
	addresses := make([]net.IP, 0, len(addrs))

	for _, addrPart := range addrs {
		if len(addrPart) == 0 {
			// Empty address part
			switch network {
			case "sctp4":
				addresses = append(addresses, net.IPv4zero)
			case "sctp", "sctp6":
				addresses = append(addresses, net.IPv6zero)
			}
		} else {
			// Parse IP address
			ip := net.ParseIP(addrPart)
			if ip == nil {
				continue // Skip invalid addresses
			}

			// Validate based on network type
			switch network {
			case "sctp4":
				if ip.To4() != nil {
					addresses = append(addresses, ip)
				}
			case "sctp6":
				if ip.To16() != nil {
					addresses = append(addresses, ip)
				}
			case "sctp":
				// Accept both IPv4 and IPv6
				addresses = append(addresses, ip)
			}
		}
	}

	if len(addresses) == 0 {
		return nil, net.InvalidAddrError(addr)
	}

	return &SCTPAddr{
		addresses: addresses,
		port:      port,
	}, nil
}
