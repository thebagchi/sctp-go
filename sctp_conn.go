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

func (conn *SCTPConn) GetPrimaryPeerAddr() (*SCTPAddr, error) {
	param := &SCTPPrimaryAddr{
		AssocId: int32(0),
	}
	length := unsafe.Sizeof(*param)
	_, _, err := syscall.Syscall6(
		syscall.SYS_GETSOCKOPT,
		uintptr(conn.sock),
		syscall.IPPROTO_SCTP,
		SCTP_PRIMARY_ADDR,
		uintptr(unsafe.Pointer(param)),
		uintptr(unsafe.Pointer(&length)),
		0,
	)
	if 0 != err {
		return nil, err
	}
	addr := FromSockAddrStorage((*SockAddrStorage)(unsafe.Pointer(&param.Addr)))
	return addr, nil
}

func (conn *SCTPConn) Read(b []byte) (n int, err error) {
	var (
		flags = 0
		info  = &SCTPSndRcvInfo{}
	)
	for {
		n, err = conn.RecvMsg(b, info, &flags)
		if flags&SCTP_MSG_NOTIFICATION <= 0 {
			return n, err
		}
	}
}

func (conn *SCTPConn) RecvMsg(b []byte, info *SCTPSndRcvInfo, flags *int) (n int, err error) {
	oob := make([]byte, syscall.CmsgSpace(SCTPSndRcvInfoSize))
	n, noob, flag, _, err := syscall.Recvmsg(int(conn.sock), b, oob, 0)
	if nil != err {
		return n, err
	}
	*flags = flag
	if noob > 0 {
		ParseSndRcvInfo(info, oob[:noob])
	}
	return n, nil
}

func (conn *SCTPConn) Write(b []byte) (n int, err error) {
	return conn.SendMsg(b, nil)
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
	var (
		data   = make([]byte, 4096)
		addrs  = (*SCTPGetAddrs)(unsafe.Pointer(&data))
		length = len(data)
	)
	addrs.AssocId = 0
	_, _, err := syscall.Syscall6(
		syscall.SYS_GETSOCKOPT,
		uintptr(conn.sock),
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

func (conn *SCTPConn) RemoteAddr() net.Addr {
	var (
		data   = make([]byte, 4096)
		addrs  = (*SCTPGetAddrs)(unsafe.Pointer(&data))
		length = len(data)
	)
	addrs.AssocId = 0

	_, _, err := syscall.Syscall6(
		syscall.SYS_GETSOCKOPT,
		uintptr(conn.sock),
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

func (conn *SCTPConn) SetDeadline(t time.Time) error {
	return syscall.ENOPROTOOPT
}

func (conn *SCTPConn) SetReadDeadline(t time.Time) error {
	return syscall.ENOPROTOOPT
}

func (conn *SCTPConn) SetWriteDeadline(t time.Time) error {
	return syscall.ENOPROTOOPT
}

func (conn *SCTPConn) SetWriteBufferSize(bytes int) error {
	return syscall.SetsockoptInt(int(conn.sock), syscall.SOL_SOCKET, syscall.SO_SNDBUF, bytes)
}

func (conn *SCTPConn) GetWriteBufferSize() (int, error) {
	return syscall.GetsockoptInt(int(conn.sock), syscall.SOL_SOCKET, syscall.SO_SNDBUF)
}

func (conn *SCTPConn) SetReadBufferSize(bytes int) error {
	return syscall.SetsockoptInt(int(conn.sock), syscall.SOL_SOCKET, syscall.SO_RCVBUF, bytes)
}

func (conn *SCTPConn) GetReadBufferSize() (int, error) {
	return syscall.GetsockoptInt(int(conn.sock), syscall.SOL_SOCKET, syscall.SO_RCVBUF)
}

func (conn *SCTPConn) SetEventSubscribe(events *SCTPEventSubscribe) error {
	_, _, err := syscall.Syscall6(
		syscall.SYS_SETSOCKOPT,
		uintptr(conn.sock),
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

func (conn *SCTPConn) GetEventSubscribe() (*SCTPEventSubscribe, error) {
	var (
		events = &SCTPEventSubscribe{}
		length = unsafe.Sizeof(*events)
	)
	_, _, err := syscall.Syscall6(
		syscall.SYS_GETSOCKOPT,
		uintptr(conn.sock),
		SOL_SCTP,
		SCTP_EVENTS,
		uintptr(unsafe.Pointer(events)),
		uintptr(unsafe.Pointer(&length)),
		0,
	)
	if 0 != err {
		return nil, err
	}
	return events, err
}

func (conn *SCTPConn) SetInitMsg(init *SCTPInitMsg) error {
	_, _, err := syscall.Syscall6(
		syscall.SYS_SETSOCKOPT,
		uintptr(conn.sock),
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

func (conn *SCTPConn) GetInitMsg() (*SCTPInitMsg, error) {
	var (
		init   = &SCTPInitMsg{}
		length = unsafe.Sizeof(*init)
	)
	_, _, err := syscall.Syscall6(
		syscall.SYS_GETSOCKOPT,
		uintptr(conn.sock),
		SOL_SCTP,
		SCTP_INITMSG,
		uintptr(unsafe.Pointer(init)),
		uintptr(unsafe.Pointer(&length)),
		0,
	)
	if 0 != err {
		return nil, err
	}
	return init, nil
}

func (conn *SCTPConn) SetDefaultSendParam(param *SCTPSndRcvInfo) error {
	_, _, err := syscall.Syscall6(
		syscall.SYS_SETSOCKOPT,
		uintptr(conn.sock),
		SOL_SCTP,
		SCTP_DEFAULT_SEND_PARAM,
		uintptr(unsafe.Pointer(param)),
		unsafe.Sizeof(*param),
		0,
	)
	if 0 != err {
		return err
	}
	return nil
}

func (conn *SCTPConn) GetDefaultSendParam() (*SCTPSndRcvInfo, error) {
	var (
		param  = &SCTPSndRcvInfo{}
		length = unsafe.Sizeof(*param)
	)
	_, _, err := syscall.Syscall6(
		syscall.SYS_SETSOCKOPT,
		uintptr(conn.sock),
		SOL_SCTP,
		SCTP_DEFAULT_SEND_PARAM,
		uintptr(unsafe.Pointer(param)),
		uintptr(unsafe.Pointer(&length)),
		0,
	)
	if 0 != err {
		return nil, err
	}
	return param, nil
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
	//syscall.SOCK_SEQPACKET vs syscall.SOCK_STREAM
	sock, err := SCTPSocket(AddrFamily(network), syscall.SOCK_STREAM)
	if nil != err {
		return nil, err
	}
	conn := &SCTPConn{
		sock: int64(sock),
	}
	for {
		err = syscall.SetsockoptInt(sock, syscall.IPPROTO_IPV6, syscall.IPV6_V6ONLY, 0)
		err = syscall.SetsockoptInt(sock, syscall.SOL_SOCKET, syscall.SO_BROADCAST, 1)
		if nil != err {
			break
		}
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
