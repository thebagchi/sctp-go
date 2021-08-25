package sctp_go

import (
	"net"
	"syscall"
	"unsafe"
)

type SCTPListener struct {
	sock int
}

func (listener *SCTPListener) FD() int {
	return listener.sock
}

func (listener *SCTPListener) Addr() net.Addr {
	var (
		data   = make([]byte, 4096)
		addrs  = (*SCTPGetAddrs)(unsafe.Pointer(&data))
		length = len(data)
	)
	addrs.AssocId = 0
	_, _, err := syscall.Syscall6(
		syscall.SYS_GETSOCKOPT,
		uintptr(listener.sock),
		syscall.IPPROTO_SCTP,
		SCTP_GET_LOCAL_ADDRS,
		uintptr(unsafe.Pointer(addrs)),
		uintptr(unsafe.Pointer(&length)),
		0,
	)
	if 0 == err {
		return FromSCTPGetAddrs(addrs)
	}
	return nil
}

func (listener *SCTPListener) RemoteAddr(assoc int) net.Addr {
	var (
		data   = make([]byte, 4096)
		addrs  = (*SCTPGetAddrs)(unsafe.Pointer(&data))
		length = len(data)
	)
	addrs.AssocId = int32(assoc)
	_, _, err := syscall.Syscall6(
		syscall.SYS_GETSOCKOPT,
		uintptr(listener.sock),
		syscall.IPPROTO_SCTP,
		SCTP_GET_PEER_ADDRS,
		uintptr(unsafe.Pointer(addrs)),
		uintptr(unsafe.Pointer(&length)),
		0,
	)
	if 0 == err {
		return FromSCTPGetAddrs(addrs)
	}
	return nil
}

func (listener *SCTPListener) Connect(remote *SCTPAddr) (int, error) {
	return SCTPConnect(listener.sock, remote)
}

func (listener *SCTPListener) Abort(assoc int) error {
	msg := &SCTPSndRcvInfo{
		Flags:   SCTP_ABORT,
		AssocId: int32(assoc),
	}
	_, err := listener.SendMsg(nil, msg)
	return err
}

func (listener *SCTPListener) Disconnect(assoc int) error {
	msg := &SCTPSndRcvInfo{
		Stream:  0,
		Ppid:    0,
		Flags:   SCTP_EOF,
		AssocId: int32(assoc),
	}
	_, err := listener.SendMsg(nil, msg)
	return err
}

func (listener *SCTPListener) PeelOff(assoc int) (*SCTPConn, error) {
	return SCTPPeelOff(listener.sock, assoc)
}

func (listener *SCTPListener) PeelOffFlags(assoc, flag int) (*SCTPConn, error) {
	return SCTPPeelOffFlag(listener.sock, assoc, flag)
}

func (listener *SCTPListener) AcceptSCTP() (*SCTPConn, error) {
	sock, _, err := syscall.Accept4(listener.sock, 0)
	if nil != err {
		return nil, err
	}
	return NewSCTPConn(sock), err
}

func (listener *SCTPListener) Accept() (net.Conn, error) {
	return listener.AcceptSCTP()
}

func (listener *SCTPListener) Close() error {
	_ = syscall.Shutdown(listener.sock, syscall.SHUT_RDWR)
	return syscall.Close(listener.sock)
}

func (listener *SCTPListener) SetEventSubscribe(events *SCTPEventSubscribe) error {
	_, _, err := syscall.Syscall6(
		syscall.SYS_SETSOCKOPT,
		uintptr(listener.sock),
		SOL_SCTP,
		SCTP_EVENTS,
		uintptr(unsafe.Pointer(&events)),
		unsafe.Sizeof(*events),
		0,
	)
	if 0 != err {
		return err
	}
	return nil
}

func (listener *SCTPListener) GetEventSubscribe() (*SCTPEventSubscribe, error) {
	var (
		events = &SCTPEventSubscribe{}
		length = unsafe.Sizeof(*events)
	)
	_, _, err := syscall.Syscall6(
		syscall.SYS_GETSOCKOPT,
		uintptr(listener.sock),
		SOL_SCTP,
		SCTP_EVENTS,
		uintptr(unsafe.Pointer(events)),
		uintptr(unsafe.Pointer(&length)),
		0,
	)
	if 0 != err {
		return nil, err
	}
	return events, nil
}

func (listener *SCTPListener) RecvMsg(b []byte, info *SCTPSndRcvInfo, flags *int) (n int, err error) {
	oob := make([]byte, syscall.CmsgSpace(SCTPSndRcvInfoSize))
	n, noob, flag, _, err := syscall.Recvmsg(listener.sock, b, oob, 0)
	if nil != err {
		return n, err
	}
	*flags = flag
	if noob > 0 {
		ParseSndRcvInfo(info, oob[:noob])
	}
	return n, nil
}

func (listener *SCTPListener) SendMsg(b []byte, info *SCTPSndRcvInfo) (int, error) {
	var buffer []byte
	if nil != info {
		hdr := &syscall.Cmsghdr{
			Level: syscall.IPPROTO_SCTP,
			Type:  SCTP_SNDRCV,
			Len:   uint64(syscall.CmsgSpace(SCTPSndRcvInfoSize)),
		}
		buffer = append(Pack(hdr), Pack(info)...)
	}
	// return syscall.SendmsgN(listener.sock, b, buffer, nil, 0)
	return SCTPSendMsg(listener.sock, b, buffer, 0)
}

func (listener *SCTPListener) SetInitMsg(init *SCTPInitMsg) error {
	_, _, err := syscall.Syscall6(
		syscall.SYS_SETSOCKOPT,
		uintptr(listener.sock),
		SOL_SCTP,
		SCTP_INITMSG,
		uintptr(unsafe.Pointer(init)),
		unsafe.Sizeof(*init),
		0,
	)
	if 0 != err {
		return err
	}
	return nil
}

func ListenSCTP(network string, sockettype int, local *SCTPAddr, init *SCTPInitMsg) (*SCTPListener, error) {
	switch network {
	case "sctp", "sctp4", "sctp6":
	default:
		return nil, &net.OpError{
			Op:     "dial",
			Net:    network,
			Source: local.Addr(),
			Err:    net.UnknownNetworkError(network),
		}
	}
	var (
		sock     int           = 0
		err      error         = nil
		listener *SCTPListener = nil
	)
	for {
		family := DetectAddrFamily(network)
		//syscall.SOCK_SEQPACKET vs syscall.SOCK_STREAM
		sock, err = SCTPSocket(family, sockettype)
		if nil != err {
			break
		}
		err = syscall.SetsockoptInt(sock, syscall.IPPROTO_IPV6, syscall.IPV6_V6ONLY, 0)
		err = syscall.SetsockoptInt(sock, syscall.SOL_SOCKET, syscall.SO_BROADCAST, 1)
		if nil != err {
			break
		}
		err = syscall.SetsockoptInt(sock, syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1)
		if nil != err {
			break
		}
		_, _, errno := syscall.Syscall6(
			syscall.SYS_SETSOCKOPT,
			uintptr(sock),
			SOL_SCTP,
			SCTP_INITMSG,
			uintptr(unsafe.Pointer(init)),
			uintptr(unsafe.Sizeof(*init)),
			0,
		)
		if 0 != errno {
			err = errno
			break
		}
		err = SCTPBind(sock, local, SCTP_BINDX_ADD_ADDR)
		if nil != err {
			break
		}
		err = syscall.Listen(sock, syscall.SOMAXCONN)
		if err != nil {
			break
		}
		listener = &SCTPListener{
			sock: sock,
		}
		break
	}
	if nil != err {
		if sock > 0 {
			_ = syscall.Close(sock)
		}
		listener = nil
	}
	return listener, err
}
