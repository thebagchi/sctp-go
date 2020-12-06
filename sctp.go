package sctp_go

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"fmt"
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

func SCTPSocket(family, flag int) (int, error) {
	switch family {
	case syscall.AF_INET:
		return syscall.Socket(syscall.AF_INET, flag, syscall.IPPROTO_SCTP)
	case syscall.AF_INET6:
		fallthrough
	default:
		return syscall.Socket(syscall.AF_INET6, flag, syscall.IPPROTO_SCTP)
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
		_, _, errno := syscall.Syscall6(
			syscall.SYS_SETSOCKOPT,
			uintptr(sock),
			SOL_SCTP,
			option,
			uintptr(unsafe.Pointer(&buffer[0])),
			uintptr(len(buffer)),
			0,
		)
		if 0 != errno {
			err = errno
		}
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
			Addrs:   uintptr(unsafe.Pointer(&buffer[0])),
		}
		length := unsafe.Sizeof(*addr)
		_, _, errno := syscall.Syscall6(
			syscall.SYS_GETSOCKOPT,
			uintptr(sock),
			syscall.IPPROTO_SCTP,
			SCTP_SOCKOPT_CONNECTX3,
			uintptr(unsafe.Pointer(addrs)),
			uintptr(unsafe.Pointer(&length)),
			0,
		)
		if 0 == errno {
			return int(addrs.AssocId), nil
		} else {
			if errno == syscall.EINPROGRESS {
				return int(addrs.AssocId), nil
			}
			if errno != syscall.ENOPROTOOPT {
				return 0, errno
			}
		}
		assoc, _, errno = syscall.Syscall6(
			syscall.SYS_SETSOCKOPT,
			uintptr(sock),
			syscall.IPPROTO_SCTP,
			SCTP_SOCKOPT_CONNECTX,
			uintptr(unsafe.Pointer(&buffer[0])),
			uintptr(len(buffer)),
			0,
		)
		if 0 != errno {
			err = errno
		}
	} else {
		return 0, syscall.EINVAL
	}
	return int(assoc), err
}

func SCTPPeelOffFlag(sock, assoc, flags int) (*SCTPConn, error) {
	if flags == 0 {
		return SCTPPeelOff(sock, assoc)
	}
	params := &SCTPPeelOffFlagsArg{
		Arg: SCTPPeelOffArg{
			AssocId: int32(assoc),
			Sd:      0,
		},
		Flags: uint32(flags),
	}
	length := unsafe.Sizeof(*params)
	_, _, errno := syscall.Syscall6(
		syscall.SYS_GETSOCKOPT,
		uintptr(sock),
		syscall.IPPROTO_SCTP,
		SCTP_SOCKOPT_PEELOFF_FLAGS,
		uintptr(unsafe.Pointer(params)),
		uintptr(unsafe.Pointer(&length)),
		0,
	)
	if 0 != errno {
		return nil, errno
	}
	return NewSCTPConn(int(params.Arg.Sd)), nil
}

func SCTPPeelOff(sock, assoc int) (*SCTPConn, error) {
	params := &SCTPPeelOffArg{
		AssocId: int32(assoc),
		Sd:      0,
	}
	length := unsafe.Sizeof(*params)
	_, _, errno := syscall.Syscall6(
		syscall.SYS_GETSOCKOPT,
		uintptr(sock),
		syscall.IPPROTO_SCTP,
		SCTP_SOCKOPT_PEELOFF,
		uintptr(unsafe.Pointer(params)),
		uintptr(unsafe.Pointer(&length)),
		0,
	)
	if 0 != errno {
		return nil, errno
	}
	return NewSCTPConn(int(params.Sd)), nil
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

func Clone(from, to interface{}) {
	buffer := new(bytes.Buffer)
	encoder := gob.NewEncoder(buffer)
	decoder := gob.NewDecoder(buffer)
	_ = encoder.Encode(from)
	_ = decoder.Decode(to)
}

func ParseSndRcvInfo(info *SCTPSndRcvInfo, data []byte) {
	if nil == info || len(data) == 0 {
		return
	}
    messages, err := syscall.ParseSocketControlMessage(data)
    if nil != err {
		return
	}
	for _, message := range messages {
		if message.Header.Level == IPPROTO_SCTP && message.Header.Type == SCTP_SNDRCV {
			temp := (*SCTPSndRcvInfo)(unsafe.Pointer(&message.Data[0]))
			Clone(temp, info)
			break
		}
	}
}

func ParseNotification(data []byte) (*Notification, error){
	if len(data) < SCTPNotificationHeaderSize {
		return nil, fmt.Errorf("data too small")
	}
	temp := (*SCTPNotificationHeader)(unsafe.Pointer(&data[0]))
	switch temp.Type() {
	
	}
	return nil, fmt.Errorf("invalid notification type")
}