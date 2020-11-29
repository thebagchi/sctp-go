package sctp_go

import (
	"bytes"
	"encoding/binary"
	"syscall"
	"unsafe"
)

var (
	endian binary.ByteOrder = binary.LittleEndian
)

func init() {
	endian = Endianness()
}

func htons(port uint16) uint16 {
	if endian == binary.LittleEndian {
		return (port << 8 & 0xff00) | (port >> 8 & 0xff)
	}
	return port
}

func ntohs(port uint16) uint16 {
	if endian == binary.LittleEndian {
		return (port << 8 & 0xff00) | (port >> 8 & 0xff)
	}
	return port
}


func Endianness() binary.ByteOrder {
	i := uint16(1)
	if *(*byte)(unsafe.Pointer(&i)) == 0 {
		return binary.BigEndian
	}
	return binary.LittleEndian
}

func DetectAddrFamily(network string) int {
	family := syscall.AF_INET6
	switch network[len(network)-1] {
	case '4':
		family = syscall.AF_INET
	case '6':
		family = syscall.AF_INET6
	}
	return family
}

func Pack(v interface{}) []byte {
	var buf bytes.Buffer
	_ = binary.Write(&buf, endian, v)
	return buf.Bytes()
}

func SCTPSocket(family int) (int, error) {
	switch family {
	case syscall.AF_INET6:
		return syscall.Socket(syscall.AF_INET6, syscall.SOCK_SEQPACKET, syscall.IPPROTO_SCTP)
	case syscall.AF_INET:
		return syscall.Socket(syscall.AF_INET, syscall.SOCK_SEQPACKET, syscall.IPPROTO_SCTP)
	default:
		return syscall.Socket(syscall.AF_INET6, syscall.SOCK_SEQPACKET, syscall.IPPROTO_SCTP)
	}
}

func SCTPBind(sock int, addr *SCTPAddr, flags int) error {
	var option uintptr
	switch flags {
	case SCTP_BINDX_ADD_ADDR:
		option = SCTP_SOCKOPT_BINDX_ADD
	case SCTP_BINDX_REM_ADDR:
		option = SCTP_SOCKOPT_BINDX_REM
	default:
		return syscall.EINVAL
	}

	var (
		buffer       = MakeSockaddr(addr)
		err    error = nil
	)
	if len(buffer) > 0 {
		_, _, err = syscall.Syscall6(
			syscall.SYS_SETSOCKOPT,
			uintptr(sock),
			SOL_SCTP,
			option,
			uintptr(unsafe.Pointer(&buffer[0])),
			unsafe.Sizeof(len(buffer)),
			0,
		)
	} else {
		err = syscall.EINVAL
	}
	return err
}

func SCTPConnect(sock int, addr *SCTPAddr) (int, error) {
	var (
		buffer         = MakeSockaddr(addr)
		err    error   = nil
		assoc  uintptr = 0
	)
	if len(buffer) > 0 {
		addrs := &SCTPGetAddrsOld{
			AssocId: 0,
			Num:     int32(len(buffer)),
			Addrs:   (*SockAddr)(unsafe.Pointer(&buffer[0])),
		}
		_, _, err = syscall.Syscall6(
			syscall.SYS_GETSOCKOPT,
			uintptr(sock),
			syscall.IPPROTO_SCTP,
			SCTP_SOCKOPT_CONNECTX3,
			uintptr(unsafe.Pointer(addrs)),
			unsafe.Sizeof(*addrs),
			0,
		)
		if nil == err {
			return int(addrs.AssocId), nil
		} else {
			if err == syscall.EINPROGRESS {
				return int(addrs.AssocId), nil
			}
			if err != syscall.ENOPROTOOPT {
				return 0, err
			}
		}
		assoc, _, err = syscall.Syscall6(
			syscall.SYS_SETSOCKOPT,
			uintptr(sock),
			syscall.IPPROTO_SCTP,
			SCTP_SOCKOPT_CONNECTX,
			uintptr(unsafe.Pointer(&buffer[0])),
			uintptr(len(buffer)),
			0,
		)
	} else {
		return 0, syscall.EINVAL
	}
	return int(assoc), err
}

func AddrFamily(network string) int {
	family := syscall.AF_INET6
	switch network[len(network)-1] {
	case '4':
		family = syscall.AF_INET
	case '6':
		family = syscall.AF_INET6
	}
	return family
}
