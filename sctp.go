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

func ParseDataIOEvent(data []byte) (Notification, error) {
	temp := (*SCTPNotificationHeader)(unsafe.Pointer(&data[0]))
	return &SCTPNotificationHeader{
		Type:   temp.Type,
		Flags:  temp.Flags,
		Length: temp.Length,
	}, nil
}

func ParseAssocChangeEvent(data []byte) (Notification, error) {
	temp := (*SCTPAssocChange)(unsafe.Pointer(&data[0]))
	return &SCTPAssocChange{
		Type:            temp.Type,
		Flags:           temp.Flags,
		Length:          temp.Length,
		State:           temp.State,
		Error:           temp.Error,
		OutboundStreams: temp.OutboundStreams,
		InboundStreams:  temp.InboundStreams,
		AssocId:         temp.AssocId,
	}, nil
}

func ParsePeerAddrChangeEvent(data []byte) (Notification, error) {
	var (
		temp = (*SCTPPAddrChange)(unsafe.Pointer(&data[0]))
		addr [128]byte
	)
	copy(addr[:], temp.Addr[:])
	return &SCTPPAddrChange{
		Type:    temp.Type,
		Flags:   temp.Flags,
		Length:  temp.Length,
		Addr:    addr,
		State:   temp.State,
		Error:   temp.Error,
		AssocId: temp.AssocId,
	}, nil
}

func ParseSendFailedEvent(data []byte) (Notification, error) {
	temp := (*SCTPSendFailed)(unsafe.Pointer(&data[0]))
	return &SCTPSendFailed{
		Type:   temp.Type,
		Flags:  temp.Flags,
		Length: temp.Length,
		Error:  temp.Error,
		Info: SCTPSndRcvInfo{
			Stream:     temp.Info.Stream,
			Ssn:        temp.Info.Ssn,
			Flags:      temp.Info.Flags,
			Ppid:       temp.Info.Ppid,
			Context:    temp.Info.Context,
			TimeToLive: temp.Info.TimeToLive,
			Tsn:        temp.Info.Tsn,
			CumTsn:     temp.Info.CumTsn,
			AssocId:    temp.Info.AssocId,
		},
		AssocId: temp.AssocId,
	}, nil
}

func ParseRemoteErrorEvent(data []byte) (Notification, error) {
	temp := (*SCTPRemoteError)(unsafe.Pointer(&data[0]))
	return &SCTPRemoteError{
		Type:    temp.Type,
		Flags:   temp.Flags,
		Length:  temp.Length,
		Error:   temp.Error,
		AssocId: temp.AssocId,
	}, nil
}

func ParseShutdownEvent(data []byte) (Notification, error) {
	temp := (*SCTPShutdownEvent)(unsafe.Pointer(&data[0]))
	return &SCTPShutdownEvent{
		Type:    temp.Type,
		Flags:   temp.Flags,
		Length:  temp.Length,
		AssocId: temp.AssocId,
	}, nil
}

func ParsePartialDeliveryEvent(data []byte) (Notification, error) {
	temp := (*SCTPPDApiEvent)(unsafe.Pointer(&data[0]))
	return &SCTPPDApiEvent{
		Type:       temp.Type,
		Flags:      temp.Flags,
		Length:     temp.Length,
		Indication: temp.Indication,
		AssocId:    temp.AssocId,
		Stream:     temp.Stream,
		Sequence:   temp.Sequence,
	}, nil
}

func ParseAdaptationIndicationEvent(data []byte) (Notification, error) {
	temp := (*SCTPAdaptationEvent)(unsafe.Pointer(&data[0]))
	return &SCTPAdaptationEvent{
		Type:          temp.Type,
		Flags:         temp.Flags,
		Length:        temp.Length,
		AdaptationInd: temp.AdaptationInd,
		AssocId:       temp.AssocId,
	}, nil
}

func ParseAuthenticationEvent(data []byte) (Notification, error) {
	temp := (*SCTPAuthKeyEvent)(unsafe.Pointer(&data[0]))
	return &SCTPAuthKeyEvent{
		Type:         temp.Type,
		Flags:        temp.Flags,
		Length:       temp.Length,
		KeyNumber:    temp.KeyNumber,
		AltKeyNumber: temp.AltKeyNumber,
		Indication:   temp.Indication,
		AssocId:      temp.AssocId,
	}, nil
}

func ParseSenderDryEvent(data []byte) (Notification, error) {
	temp := (*SCTPSenderDryEvent)(unsafe.Pointer(&data[0]))
	return &SCTPSenderDryEvent{
		Type:    temp.Type,
		Flags:   temp.Flags,
		Length:  temp.Length,
		AssocId: temp.AssocId,
	}, nil
}

func ParseStreamResetEvent(data []byte) (Notification, error) {
	temp := (*SCTPStreamResetEvent)(unsafe.Pointer(&data[0]))
	return &SCTPStreamResetEvent{
		Type:    temp.Type,
		Flags:   temp.Flags,
		Length:  temp.Length,
		AssocId: temp.AssocId,
	}, nil
}

func ParseAssocResetEvent(data []byte) (Notification, error) {
	temp := (*SCTPAssocResetEvent)(unsafe.Pointer(&data[0]))
	return &SCTPAssocResetEvent{
		Type:      temp.Type,
		Flags:     temp.Flags,
		Length:    temp.Length,
		AssocId:   temp.AssocId,
		LocalTsn:  temp.LocalTsn,
		RemoteTsn: temp.RemoteTsn,
	}, nil
}

func ParseStreamChangeEvent(data []byte) (Notification, error) {
	temp := (*SCTPStreamChangeEvent)(unsafe.Pointer(&data[0]))
	return &SCTPStreamChangeEvent{
		Type:       temp.Type,
		Flags:      temp.Flags,
		Length:     temp.Length,
		AssocId:    temp.AssocId,
		InStreams:  temp.InStreams,
		OutStreams: temp.OutStreams,
	}, nil
}

func NotificationName(notification uint16) string {
	switch notification {
	case SCTP_DATA_IO_EVENT:
		return "SCTP_DATA_IO_EVENT"
	case SCTP_ASSOC_CHANGE:
		return "SCTP_ASSOC_CHANGE"
	case SCTP_PEER_ADDR_CHANGE:
		return "SCTP_PEER_ADDR_CHANGE"
	case SCTP_SEND_FAILED:
		return "SCTP_SEND_FAILED"
	case SCTP_REMOTE_ERROR:
		return "SCTP_REMOTE_ERROR"
	case SCTP_SHUTDOWN_EVENT:
		return "SCTP_SHUTDOWN_EVENT"
	case SCTP_PARTIAL_DELIVERY_EVENT:
		return "SCTP_PARTIAL_DELIVERY_EVENT"
	case SCTP_ADAPTATION_INDICATION:
		return "SCTP_ADAPTATION_INDICATION"
	case SCTP_AUTHENTICATION_EVENT:
		return "SCTP_AUTHENTICATION_EVENT"
	case SCTP_SENDER_DRY_EVENT:
		return "SCTP_SENDER_DRY_EVENT"
	case SCTP_STREAM_RESET_EVENT:
		return "SCTP_STREAM_RESET_EVENT"
	case SCTP_ASSOC_RESET_EVENT:
		return "SCTP_ASSOC_RESET_EVENT"
	case SCTP_STREAM_CHANGE_EVENT:
		return "SCTP_STREAM_CHANGE_EVENT"
	}
	return ""
}

func ParseNotification(data []byte) (Notification, error) {
	if len(data) < SCTPNotificationHeaderSize {
		return nil, fmt.Errorf("invalid data len, too small")
	}
	if len(data) > SCTPNotificationSize {
		return nil, fmt.Errorf("invalid data len, too large")
	}
	temp := (*SCTPNotificationHeader)(unsafe.Pointer(&data[0]))
	switch temp.Type {
	case SCTP_DATA_IO_EVENT:
		return ParseDataIOEvent(data)
	case SCTP_ASSOC_CHANGE:
		return ParseAssocChangeEvent(data)
	case SCTP_PEER_ADDR_CHANGE:
		return ParsePeerAddrChangeEvent(data)
	case SCTP_SEND_FAILED:
		return ParseSendFailedEvent(data)
	case SCTP_REMOTE_ERROR:
		return ParseRemoteErrorEvent(data)
	case SCTP_SHUTDOWN_EVENT:
		return ParseShutdownEvent(data)
	case SCTP_PARTIAL_DELIVERY_EVENT:
		return ParsePartialDeliveryEvent(data)
	case SCTP_ADAPTATION_INDICATION:
		return ParseAdaptationIndicationEvent(data)
	case SCTP_AUTHENTICATION_EVENT:
		return ParseAuthenticationEvent(data)
	case SCTP_SENDER_DRY_EVENT:
		return ParseSenderDryEvent(data)
	case SCTP_STREAM_RESET_EVENT:
		return ParseStreamResetEvent(data)
	case SCTP_ASSOC_RESET_EVENT:
		return ParseAssocResetEvent(data)
	case SCTP_STREAM_CHANGE_EVENT:
		return ParseStreamChangeEvent(data)
	}
	return nil, fmt.Errorf("invalid notification type")
}

func SCTPSendMsg(sock int, buffer, control []byte, flags int) (int, error) {
	var (
		ptr unsafe.Pointer
		msg syscall.Msghdr
		iov syscall.Iovec
	)
	msg.Name = (*byte)(ptr)
	msg.Namelen = uint32(0)
	if len(buffer) > 0 {
		iov.Base = &buffer[0]
		iov.SetLen(len(buffer))
	}
	if len(control) > 0 {
		msg.Control = &control[0]
		msg.SetControllen(len(control))
	}
	msg.Iov = &iov
	msg.Iovlen = 1
	length, _, errno := syscall.Syscall(
		syscall.SYS_SENDMSG,
		uintptr(sock),
		uintptr(unsafe.Pointer(&msg)),
		uintptr(flags),
	)
	if 0 != errno {
		return 0, errno
	}
	if len(control) > 0 && len(buffer) == 0 {
		return 0, nil
	}
	return int(length), nil
}
