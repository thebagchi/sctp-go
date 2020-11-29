package sctp_go

import (
	"net"
	"sync/atomic"
	"syscall"
	"time"
	"unsafe"
)

type SCTPConn struct {
	sock  int64
	assoc int
}

func NewSCTPConn(sock int) *SCTPConn {
	return &SCTPConn{
		sock: int64(sock),
	}
}

func (conn *SCTPConn) Read(b []byte) (n int, err error) {
	return 0, nil
}

func (conn *SCTPConn) Write(b []byte) (n int, err error) {
	return 0, nil
}

func (conn *SCTPConn) SendMsg(b []byte, info *SCTPSndRcvInfo) (int, error) {
	var buffer []byte
	if nil != info {
		hdr := &syscall.Cmsghdr{
			Level: syscall.IPPROTO_SCTP,
			Type:  SCTP_SNDRCV,
			Len:   uint64(syscall.CmsgSpace(SCTPSndRcvInfoSize)),
		}
		buffer = append(buffer, Pack(hdr)...)
		buffer = append(buffer, Pack(info)...)
	}
	return syscall.SendmsgN(int(conn.sock), b, buffer, nil, 0)
}

func (conn *SCTPConn) Close() error {
	if !conn.ok() {
		return syscall.EINVAL
	}
	sock := atomic.SwapInt64(&conn.sock, -1)
	if sock > 0 {
		msg := &SCTPSndRcvInfo{
			Flags: SCTP_EOF,
		}
		_, _ = conn.SendMsg(nil, msg)
		_ = syscall.Shutdown(int(sock), syscall.SHUT_RDWR)
		return syscall.Close(int(sock))
	}
	return syscall.EBADFD
}

func (conn *SCTPConn) LocalAddr() net.Addr {
	return nil
}

func (conn *SCTPConn) RemoteAddr() net.Addr {
	return nil
}

func (conn *SCTPConn) SetDeadline(t time.Time) error {
	return nil
}

func (conn *SCTPConn) SetReadDeadline(t time.Time) error {
	return nil
}

func (conn *SCTPConn) SetWriteDeadline(t time.Time) error {
	return nil
}

func (conn *SCTPConn) SetInitMsg(init *SCTPInitMsg) error {
	_, _, err := syscall.Syscall6(
		syscall.SYS_SETSOCKOPT,
		uintptr(conn.sock),
		SOL_SCTP,
		SCTP_INITMSG,
		uintptr(unsafe.Pointer(&init)),
		unsafe.Sizeof(init),
		0,
	)
	return err
}

func (conn *SCTPConn) ok() bool {
	if nil != conn && conn.sock > 0 {
		return true
	}
	return false
}

func DialSCTP(network string, local, remote *SCTPAddr, init *SCTPInitMsg) (*SCTPConn, error) {
	switch network {
	case "sctp", "sctp4", "sctp6":
	default:
		return nil, &net.OpError{
			Op:     "dial",
			Net:    network,
			Source: local.Addr(),
			Addr:   remote.Addr(),
			Err:    net.UnknownNetworkError(network),
		}
	}
	if nil == remote {
		return nil, &net.OpError{
			Op:     "dial",
			Net:    network,
			Source: local.Addr(),
			Addr:   remote.Addr(),
			Err:    net.InvalidAddrError("invalid remote addr"),
		}
	}
	sock, err := SCTPSocket(AddrFamily(network))
	if nil != err {
		return nil, err
	}
	conn := &SCTPConn{
		sock: int64(sock),
	}
	for {
		err = conn.SetInitMsg(init)
		if nil != err {
			break
		}
		if nil != local {
			err = SCTPBind(sock, local, SCTP_BINDX_ADD_ADDR)
			if nil != err {
				break
			}
		}
		conn.assoc, err = SCTPConnect(sock, remote)
        if nil != err {
        	break
		}
		break
	}
	if nil != err {
		_ = conn.Close()
		return nil, err
	}
	return conn, nil
}
