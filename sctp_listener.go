package sctp_go

import (
	"bytes"
	"encoding/binary"
	"errors"
	"net"
	"syscall"
	"unsafe"
)

type SCTPListener struct {
	sock int
}

// FD returns the file descriptor of the listener socket.
func (listener *SCTPListener) FD() int {
	return listener.sock
}

// Addr returns the local network address.
func (listener *SCTPListener) Addr() net.Addr {
	if listener.sock <= 0 {
		return nil
	}
	var (
		data   [4096]byte
		addrs  = (*SCTPGetAddrs)(unsafe.Pointer(&data[0]))
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
	if err == 0 {
		return FromSCTPGetAddrs(addrs)
	}
	return nil
}

// RemoteAddr returns the remote network address for the specified association.
func (listener *SCTPListener) RemoteAddr(assoc int) net.Addr {
	if listener.sock <= 0 {
		return nil
	}
	var (
		data   [4096]byte
		addrs  = (*SCTPGetAddrs)(unsafe.Pointer(&data[0]))
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
	if err == 0 {
		return FromSCTPGetAddrs(addrs)
	}
	return nil
}

// Connect connects to the remote SCTP address.
func (listener *SCTPListener) Connect(remote *SCTPAddr) (int, error) {
	if listener.sock <= 0 {
		return 0, errors.New("invalid listener")
	}
	return SCTPConnect(listener.sock, remote)
}

// Abort aborts the SCTP association specified by assoc.
func (listener *SCTPListener) Abort(assoc int) error {
	if listener.sock <= 0 {
		return errors.New("invalid listener")
	}
	var msg = SCTPSndRcvInfo{
		Flags:   SCTP_ABORT,
		AssocId: int32(assoc),
	}
	_, err := listener.SendMsg(nil, &msg)
	return err
}

// Disconnect disconnects the SCTP association specified by assoc.
func (listener *SCTPListener) Disconnect(assoc int) error {
	if listener.sock <= 0 {
		return errors.New("invalid listener")
	}
	var msg = SCTPSndRcvInfo{
		Stream:  0,
		Ppid:    0,
		Flags:   SCTP_EOF,
		AssocId: int32(assoc),
	}
	_, err := listener.SendMsg(nil, &msg)
	return err
}

// PeelOff peels off the SCTP association specified by assoc.
func (listener *SCTPListener) PeelOff(assoc int) (*SCTPConn, error) {
	if listener.sock <= 0 {
		return nil, errors.New("invalid listener")
	}
	return SCTPPeelOff(listener.sock, assoc)
}

// PeelOffFlags peels off the SCTP association specified by assoc with flags.
func (listener *SCTPListener) PeelOffFlags(assoc, flag int) (*SCTPConn, error) {
	if listener.sock <= 0 {
		return nil, errors.New("invalid listener")
	}
	return SCTPPeelOffFlag(listener.sock, assoc, flag)
}

// AcceptSCTP accepts an incoming SCTP connection.
func (listener *SCTPListener) AcceptSCTP() (*SCTPConn, error) {
	if listener.sock <= 0 {
		return nil, errors.New("invalid listener")
	}
	sock, _, err := syscall.Accept4(listener.sock, 0)
	if err != nil {
		return nil, err
	}
	return NewSCTPConn(sock), nil
}

// Accept accepts a generic network connection.
func (listener *SCTPListener) Accept() (net.Conn, error) {
	return listener.AcceptSCTP()
}

// Close closes the listener.
func (listener *SCTPListener) Close() error {
	if listener.sock <= 0 {
		return errors.New("invalid listener")
	}
	// Shutdown doesn't work on RAW Sockets, SCTP Socket is essentially a RAW Socket.
	_ = syscall.Shutdown(listener.sock, syscall.SHUT_RDWR)
	return syscall.Close(listener.sock)
}

// SetEventSubscribe sets the SCTP event subscription.
func (listener *SCTPListener) SetEventSubscribe(events *SCTPEventSubscribe) error {
	if listener.sock <= 0 {
		return errors.New("invalid listener")
	}
	if events == nil {
		return errors.New("events cannot be nil")
	}
	_, _, err := syscall.Syscall6(
		syscall.SYS_SETSOCKOPT,
		uintptr(listener.sock),
		SOL_SCTP,
		SCTP_EVENTS,
		uintptr(unsafe.Pointer(events)),
		unsafe.Sizeof(*events),
		0,
	)
	if err != 0 {
		return err
	}
	return nil
}

// GetEventSubscribe gets the SCTP event subscription.
func (listener *SCTPListener) GetEventSubscribe() (*SCTPEventSubscribe, error) {
	if listener.sock <= 0 {
		return nil, errors.New("invalid listener")
	}
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
	if err != 0 {
		return nil, err
	}
	return events, nil
}

// RecvMsg receives a message from the SCTP socket.
func (listener *SCTPListener) RecvMsg(b []byte, info *SCTPSndRcvInfo, flags *int) (n int, err error) {
	if listener.sock <= 0 {
		return 0, errors.New("invalid listener")
	}
	var (
		oob  = make([]byte, syscall.CmsgSpace(SCTPSndRcvInfoSize))
		flag = 0
	)
	if flags != nil {
		flag = *flags
	}
	n, noob, flag, _, err := syscall.Recvmsg(listener.sock, b, oob[:], flag)
	if err != nil {
		return n, err
	}
	*flags = flag
	if noob > 0 {
		ParseSndRcvInfo(info, oob[:noob])
	}
	return n, nil
}

// SendMsg sends a message on the SCTP socket.
func (listener *SCTPListener) SendMsg(b []byte, info *SCTPSndRcvInfo) (int, error) {
	if listener.sock <= 0 {
		return 0, errors.New("invalid listener")
	}
	var buffer bytes.Buffer
	if info != nil {
		size := int(unsafe.Sizeof(syscall.Cmsghdr{})) + int(unsafe.Sizeof(*info))
		buffer.Grow(size)
		hdr := syscall.Cmsghdr{
			Level: syscall.IPPROTO_SCTP,
			Type:  SCTP_SNDRCV,
			Len:   uint64(syscall.CmsgSpace(SCTPSndRcvInfoSize)),
		}
		binary.Write(&buffer, endian, hdr)
		binary.Write(&buffer, endian, *info)
	}
	return SCTPSendMsg(listener.sock, b, buffer.Bytes(), 0)
}

// SetInitMsg sets the SCTP initialization message.
func (listener *SCTPListener) SetInitMsg(init *SCTPInitMsg) error {
	if listener.sock <= 0 {
		return errors.New("invalid listener")
	}
	if init == nil {
		return errors.New("init cannot be nil")
	}
	_, _, err := syscall.Syscall6(
		syscall.SYS_SETSOCKOPT,
		uintptr(listener.sock),
		SOL_SCTP,
		SCTP_INITMSG,
		uintptr(unsafe.Pointer(init)),
		unsafe.Sizeof(*init),
		0,
	)
	if err != 0 {
		return err
	}
	return nil
}

// SetNonblock sets the listener socket to non-blocking mode.
func (listener *SCTPListener) SetNonblock() error {
	if listener.sock <= 0 {
		return errors.New("invalid listener")
	}
	return syscall.SetNonblock(listener.sock, true)
}

// ListenSCTP creates an SCTP listener on the specified network and address.
func ListenSCTP(network string, sockettype int, local *SCTPAddr, init *SCTPInitMsg) (*SCTPListener, error) {
	if local == nil {
		return nil, errors.New("local address cannot be nil")
	}
	if init == nil {
		return nil, errors.New("init message cannot be nil")
	}
	switch network {
	case "sctp", "sctp4", "sctp6":
	default:
		return nil, &net.OpError{
			Op:     "listen",
			Net:    network,
			Source: local.Addr(),
			Err:    net.UnknownNetworkError(network),
		}
	}
	var (
		sock     int
		err      error
		listener *SCTPListener
	)
	for {
		family := DetectAddrFamily(network)
		// syscall.SOCK_SEQPACKET vs syscall.SOCK_STREAM
		sock, err = SCTPSocket(family, sockettype)
		if err != nil {
			break
		}
		err = syscall.SetsockoptInt(sock, syscall.IPPROTO_IPV6, syscall.IPV6_V6ONLY, 0)
		if err != nil {
			break
		}
		err = syscall.SetsockoptInt(sock, syscall.SOL_SOCKET, syscall.SO_BROADCAST, 1)
		if err != nil {
			break
		}
		err = syscall.SetsockoptInt(sock, syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1)
		if err != nil {
			break
		}
		_, _, errno := syscall.Syscall6(
			syscall.SYS_SETSOCKOPT,
			uintptr(sock),
			SOL_SCTP,
			SCTP_INITMSG,
			uintptr(unsafe.Pointer(init)),
			unsafe.Sizeof(*init),
			0,
		)
		if errno != 0 {
			err = errno
			break
		}
		err = SCTPBind(sock, local, SCTP_BINDX_ADD_ADDR)
		if err != nil {
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
	if err != nil {
		if sock > 0 {
			_ = syscall.Close(sock)
		}
		return nil, err
	}
	return listener, nil
}
