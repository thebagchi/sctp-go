package sctp_go

import (
	"net"
	"syscall"
	"unsafe"
)

type SCTPListener struct {
	sock int
}

func (listener *SCTPListener) Addr() net.Addr {
	return nil
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

func ListenSCTP(network string, local *SCTPAddr, init *SCTPInitMsg) (*SCTPListener, error) {
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
		sock, err = SCTPSocket(family, syscall.SOCK_STREAM)
		if nil != err {
			break
		}
		err = syscall.SetsockoptInt(sock, syscall.IPPROTO_IPV6, syscall.IPV6_V6ONLY, 0)
		err = syscall.SetsockoptInt(sock, syscall.SOL_SOCKET, syscall.SO_BROADCAST, 1)
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
