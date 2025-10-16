package sctp_go

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"strings"
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

// Endianness returns the byte order (big-endian or little-endian) of the system.
func Endianness() binary.ByteOrder {
	i := uint16(1)
	if *(*byte)(unsafe.Pointer(&i)) == 0 {
		return binary.BigEndian
	}
	return binary.LittleEndian
}

// DetectAddrFamily detects the address family (AF_INET or AF_INET6) from the network string.
func DetectAddrFamily(network string) int {
	if strings.HasSuffix(network, "4") {
		return syscall.AF_INET
	}
	return syscall.AF_INET6
}

// Pack serializes the given value into a byte slice using the system's endianness.
func Pack(v interface{}) []byte {
	var buf bytes.Buffer
	if err := binary.Write(&buf, endian, v); err != nil {
		return nil
	}
	return buf.Bytes()
}

// SCTPSocket creates a new SCTP socket with the specified address family and type.
func SCTPSocket(family, flag int) (int, error) {
	if family == syscall.AF_INET {
		return syscall.Socket(syscall.AF_INET, flag, syscall.IPPROTO_SCTP)
	}
	return syscall.Socket(syscall.AF_INET6, flag, syscall.IPPROTO_SCTP)
}

// SCTPBind binds the SCTP socket to the specified address with the given flags.
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

	buffer := MakeSockaddr(addr)
	if len(buffer) == 0 {
		return syscall.EINVAL
	}
	_, _, errno := syscall.Syscall6(
		syscall.SYS_SETSOCKOPT,
		uintptr(sock),
		SOL_SCTP,
		option,
		uintptr(unsafe.Pointer(&buffer[0])),
		uintptr(len(buffer)),
		0,
	)
	if errno != 0 {
		return errno
	}
	return nil
}

// SCTPConnect connects the SCTP socket to the specified address and returns the association ID.
func SCTPConnect(sock int, addr *SCTPAddr) (int, error) {
	buffer := MakeSockaddr(addr)
	if len(buffer) == 0 {
		return 0, syscall.EINVAL
	}
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
	if errno == 0 || errno == syscall.EINPROGRESS {
		return int(addrs.AssocId), nil
	}
	if errno != syscall.ENOPROTOOPT {
		return 0, errno
	}
	// Fallback to CONNECTX
	assoc, _, errno := syscall.Syscall6(
		syscall.SYS_SETSOCKOPT,
		uintptr(sock),
		syscall.IPPROTO_SCTP,
		SCTP_SOCKOPT_CONNECTX,
		uintptr(unsafe.Pointer(&buffer[0])),
		uintptr(len(buffer)),
		0,
	)
	if errno != 0 {
		return 0, errno
	}
	return int(assoc), nil
}

// SCTPPeelOffFlag peels off an association from the SCTP socket with the specified flags and returns a new SCTPConn.
func SCTPPeelOffFlag(sock, assoc, flags int) (*SCTPConn, error) {
	if flags == 0 {
		return SCTPPeelOff(sock, assoc)
	}
	var params = SCTPPeelOffFlagsArg{
		Arg: SCTPPeelOffArg{
			AssocId: int32(assoc),
			Sd:      0,
		},
		Flags: uint32(flags),
	}
	length := unsafe.Sizeof(params)
	_, _, errno := syscall.Syscall6(
		syscall.SYS_GETSOCKOPT,
		uintptr(sock),
		syscall.IPPROTO_SCTP,
		SCTP_SOCKOPT_PEELOFF_FLAGS,
		uintptr(unsafe.Pointer(&params)),
		uintptr(unsafe.Pointer(&length)),
		0,
	)
	if errno != 0 {
		return nil, errno
	}
	return NewSCTPConn(int(params.Arg.Sd)), nil
}

// SCTPPeelOff peels off an association from the SCTP socket and returns a new SCTPConn.
func SCTPPeelOff(sock, assoc int) (*SCTPConn, error) {
	var params = SCTPPeelOffArg{
		AssocId: int32(assoc),
		Sd:      0,
	}
	length := unsafe.Sizeof(params)
	_, _, errno := syscall.Syscall6(
		syscall.SYS_GETSOCKOPT,
		uintptr(sock),
		syscall.IPPROTO_SCTP,
		SCTP_SOCKOPT_PEELOFF,
		uintptr(unsafe.Pointer(&params)),
		uintptr(unsafe.Pointer(&length)),
		0,
	)
	if errno != 0 {
		return nil, errno
	}
	return NewSCTPConn(int(params.Sd)), nil
}

// AddrFamily returns the address family (AF_INET or AF_INET6) based on the network string.
func AddrFamily(network string) int {
	if strings.HasSuffix(network, "4") {
		return syscall.AF_INET
	}
	return syscall.AF_INET6
}

// Clone performs a deep copy of the 'from' interface to the 'to' interface using gob encoding.
func Clone(from, to interface{}) error {
	var (
		buffer  = new(bytes.Buffer)
		encoder = gob.NewEncoder(buffer)
		decoder = gob.NewDecoder(buffer)
	)
	if err := encoder.Encode(from); err != nil {
		return err
	}
	return decoder.Decode(to)
}

// ParseSndRcvInfo parses the socket control message data and populates the SCTPSndRcvInfo struct.
func ParseSndRcvInfo(info *SCTPSndRcvInfo, data []byte) {
	if info == nil || len(data) == 0 {
		return
	}
	messages, err := syscall.ParseSocketControlMessage(data)
	if err != nil {
		return
	}
	for _, message := range messages {
		if message.Header.Level == IPPROTO_SCTP && message.Header.Type == SCTP_SNDRCV {
			temp := (*SCTPSndRcvInfo)(unsafe.Pointer(&message.Data[0]))
			*info = *temp
			break
		}
	}
}

// ParseDataIOEvent parses the notification data into a SCTPNotificationHeader.
func ParseDataIOEvent(data []byte) (Notification, error) {
	if len(data) < SCTPNotificationHeaderSize {
		return nil, fmt.Errorf("invalid data len, too small")
	}
	temp := (*SCTPNotificationHeader)(unsafe.Pointer(&data[0]))
	return &SCTPNotificationHeader{
		Type:   temp.Type,
		Flags:  temp.Flags,
		Length: temp.Length,
	}, nil
}

// ParseAssocChangeEvent parses the notification data into a SCTPAssocChange.
func ParseAssocChangeEvent(data []byte) (Notification, error) {
	if len(data) < int(unsafe.Sizeof(SCTPAssocChange{})) {
		return nil, fmt.Errorf("invalid data len, too small")
	}
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

// ParsePeerAddrChangeEvent parses the notification data into a SCTPPAddrChange.
func ParsePeerAddrChangeEvent(data []byte) (Notification, error) {
	if len(data) < int(unsafe.Sizeof(SCTPPAddrChange{})) {
		return nil, fmt.Errorf("invalid data len, too small")
	}
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

// ParseSendFailedEvent parses the notification data into a SCTPSendFailed.
func ParseSendFailedEvent(data []byte) (Notification, error) {
	if len(data) < int(unsafe.Sizeof(SCTPSendFailed{})) {
		return nil, fmt.Errorf("invalid data len, too small")
	}
	temp := (*SCTPSendFailed)(unsafe.Pointer(&data[0]))
	return &SCTPSendFailed{
		Type:    temp.Type,
		Flags:   temp.Flags,
		Length:  temp.Length,
		Error:   temp.Error,
		Info:    temp.Info,
		AssocId: temp.AssocId,
	}, nil
}

// ParseRemoteErrorEvent parses the notification data into a SCTPRemoteError.
func ParseRemoteErrorEvent(data []byte) (Notification, error) {
	if len(data) < int(unsafe.Sizeof(SCTPRemoteError{})) {
		return nil, fmt.Errorf("invalid data len, too small")
	}
	temp := (*SCTPRemoteError)(unsafe.Pointer(&data[0]))
	return &SCTPRemoteError{
		Type:    temp.Type,
		Flags:   temp.Flags,
		Length:  temp.Length,
		Error:   temp.Error,
		AssocId: temp.AssocId,
	}, nil
}

// ParseShutdownEvent parses the notification data into a SCTPShutdownEvent.
func ParseShutdownEvent(data []byte) (Notification, error) {
	if len(data) < int(unsafe.Sizeof(SCTPShutdownEvent{})) {
		return nil, fmt.Errorf("invalid data len, too small")
	}
	temp := (*SCTPShutdownEvent)(unsafe.Pointer(&data[0]))
	return &SCTPShutdownEvent{
		Type:    temp.Type,
		Flags:   temp.Flags,
		Length:  temp.Length,
		AssocId: temp.AssocId,
	}, nil
}

// ParsePartialDeliveryEvent parses the notification data into a SCTPPDApiEvent.
func ParsePartialDeliveryEvent(data []byte) (Notification, error) {
	if len(data) < int(unsafe.Sizeof(SCTPPDApiEvent{})) {
		return nil, fmt.Errorf("invalid data len, too small")
	}
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

// ParseAdaptationIndicationEvent parses the notification data into a SCTPAdaptationEvent.
func ParseAdaptationIndicationEvent(data []byte) (Notification, error) {
	if len(data) < int(unsafe.Sizeof(SCTPAdaptationEvent{})) {
		return nil, fmt.Errorf("invalid data len, too small")
	}
	temp := (*SCTPAdaptationEvent)(unsafe.Pointer(&data[0]))
	return &SCTPAdaptationEvent{
		Type:          temp.Type,
		Flags:         temp.Flags,
		Length:        temp.Length,
		AdaptationInd: temp.AdaptationInd,
		AssocId:       temp.AssocId,
	}, nil
}

// ParseAuthenticationEvent parses the notification data into a SCTPAuthKeyEvent.
func ParseAuthenticationEvent(data []byte) (Notification, error) {
	if len(data) < int(unsafe.Sizeof(SCTPAuthKeyEvent{})) {
		return nil, fmt.Errorf("invalid data len, too small")
	}
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

// ParseSenderDryEvent parses the notification data into a SCTPSenderDryEvent.
func ParseSenderDryEvent(data []byte) (Notification, error) {
	if len(data) < int(unsafe.Sizeof(SCTPSenderDryEvent{})) {
		return nil, fmt.Errorf("invalid data len, too small")
	}
	temp := (*SCTPSenderDryEvent)(unsafe.Pointer(&data[0]))
	return &SCTPSenderDryEvent{
		Type:    temp.Type,
		Flags:   temp.Flags,
		Length:  temp.Length,
		AssocId: temp.AssocId,
	}, nil
}

// ParseStreamResetEvent parses the notification data into a SCTPStreamResetEvent.
func ParseStreamResetEvent(data []byte) (Notification, error) {
	if len(data) < int(unsafe.Sizeof(SCTPStreamResetEvent{})) {
		return nil, fmt.Errorf("invalid data len, too small")
	}
	temp := (*SCTPStreamResetEvent)(unsafe.Pointer(&data[0]))
	return &SCTPStreamResetEvent{
		Type:    temp.Type,
		Flags:   temp.Flags,
		Length:  temp.Length,
		AssocId: temp.AssocId,
	}, nil
}

// ParseAssocResetEvent parses the notification data into a SCTPAssocResetEvent.
func ParseAssocResetEvent(data []byte) (Notification, error) {
	if len(data) < int(unsafe.Sizeof(SCTPAssocResetEvent{})) {
		return nil, fmt.Errorf("invalid data len, too small")
	}
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

// ParseStreamChangeEvent parses the notification data into a SCTPStreamChangeEvent.
func ParseStreamChangeEvent(data []byte) (Notification, error) {
	if len(data) < int(unsafe.Sizeof(SCTPStreamChangeEvent{})) {
		return nil, fmt.Errorf("invalid data len, too small")
	}
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

// NotificationName returns the string name of the given SCTP notification type.
func NotificationName(notification uint16) string {
	names := map[uint16]string{
		SCTP_DATA_IO_EVENT:          "SCTP_DATA_IO_EVENT",
		SCTP_ASSOC_CHANGE:           "SCTP_ASSOC_CHANGE",
		SCTP_PEER_ADDR_CHANGE:       "SCTP_PEER_ADDR_CHANGE",
		SCTP_SEND_FAILED:            "SCTP_SEND_FAILED",
		SCTP_REMOTE_ERROR:           "SCTP_REMOTE_ERROR",
		SCTP_SHUTDOWN_EVENT:         "SCTP_SHUTDOWN_EVENT",
		SCTP_PARTIAL_DELIVERY_EVENT: "SCTP_PARTIAL_DELIVERY_EVENT",
		SCTP_ADAPTATION_INDICATION:  "SCTP_ADAPTATION_INDICATION",
		SCTP_AUTHENTICATION_EVENT:   "SCTP_AUTHENTICATION_EVENT",
		SCTP_SENDER_DRY_EVENT:       "SCTP_SENDER_DRY_EVENT",
		SCTP_STREAM_RESET_EVENT:     "SCTP_STREAM_RESET_EVENT",
		SCTP_ASSOC_RESET_EVENT:      "SCTP_ASSOC_RESET_EVENT",
		SCTP_STREAM_CHANGE_EVENT:    "SCTP_STREAM_CHANGE_EVENT",
	}
	return names[notification]
}

// ParseNotification parses the SCTP notification data based on its type and returns the appropriate Notification.
func ParseNotification(data []byte) (Notification, error) {
	if len(data) < SCTPNotificationHeaderSize {
		return nil, fmt.Errorf("invalid data len, too small")
	}
	if len(data) > SCTPNotificationSize {
		return nil, fmt.Errorf("invalid data len, too large")
	}
	temp := (*SCTPNotificationHeader)(unsafe.Pointer(&data[0]))
	parsers := map[uint16]func([]byte) (Notification, error){
		SCTP_DATA_IO_EVENT:          ParseDataIOEvent,
		SCTP_ASSOC_CHANGE:           ParseAssocChangeEvent,
		SCTP_PEER_ADDR_CHANGE:       ParsePeerAddrChangeEvent,
		SCTP_SEND_FAILED:            ParseSendFailedEvent,
		SCTP_REMOTE_ERROR:           ParseRemoteErrorEvent,
		SCTP_SHUTDOWN_EVENT:         ParseShutdownEvent,
		SCTP_PARTIAL_DELIVERY_EVENT: ParsePartialDeliveryEvent,
		SCTP_ADAPTATION_INDICATION:  ParseAdaptationIndicationEvent,
		SCTP_AUTHENTICATION_EVENT:   ParseAuthenticationEvent,
		SCTP_SENDER_DRY_EVENT:       ParseSenderDryEvent,
		SCTP_STREAM_RESET_EVENT:     ParseStreamResetEvent,
		SCTP_ASSOC_RESET_EVENT:      ParseAssocResetEvent,
		SCTP_STREAM_CHANGE_EVENT:    ParseStreamChangeEvent,
	}
	if parser, ok := parsers[temp.Type]; ok {
		return parser(data)
	}
	return nil, fmt.Errorf("invalid notification type")
}

// SCTPSendMsg sends a message with optional control data over the SCTP socket.
func SCTPSendMsg(sock int, buffer, control []byte, flags int) (int, error) {
	var (
		msg syscall.Msghdr
		iov syscall.Iovec
	)
	msg.Name = nil
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
	if errno != 0 {
		return 0, errno
	}
	if len(control) > 0 && len(buffer) == 0 {
		return 0, nil
	}
	return int(length), nil
}
