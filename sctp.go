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