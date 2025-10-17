package sctp_go

import (
	"net"
	"sync/atomic"
	"syscall"
	"time"
	"unsafe"
)

// SCTPConn represents an SCTP connection.
type SCTPConn struct {
	sock  int64
	assoc int
}

// NewSCTPConn creates a new SCTPConn from a socket file descriptor.
func NewSCTPConn(sock int) *SCTPConn {
	return &SCTPConn{
		sock: int64(sock),
	}
}

// FD returns the socket file descriptor.
func (conn *SCTPConn) FD() int64 {
	return conn.sock
}

// AssocId returns the association ID.
func (conn *SCTPConn) AssocId() int {
	return conn.assoc
}

// GetPrimaryPeerAddr returns the primary peer address for the association.
func (conn *SCTPConn) GetPrimaryPeerAddr() (*SCTPAddr, error) {
	param := SCTPPrimaryAddr{
		AssocId: int32(conn.assoc),
	}
	length := unsafe.Sizeof(param)
	_, _, err := syscall.Syscall6(
		syscall.SYS_GETSOCKOPT,
		uintptr(conn.sock),
		syscall.IPPROTO_SCTP,
		SCTP_PRIMARY_ADDR,
		uintptr(unsafe.Pointer(&param)),
		uintptr(unsafe.Pointer(&length)),
		0,
	)
	if err != 0 {
		return nil, err
	}
	addr := FromSockAddrStorage((*SockAddrStorage)(unsafe.Pointer(&param.Addr)))
	return addr, nil
}

// Read reads data from the connection, skipping notifications.
func (conn *SCTPConn) Read(b []byte) (n int, err error) {
	if !conn.ok() {
		return 0, syscall.EINVAL
	}
	var (
		flags = 0
		info  = &SCTPSndRcvInfo{}
	)
	for {
		n, err = conn.RecvMsg(b, info, &flags)
		if flags&SCTP_MSG_NOTIFICATION == 0 {
			return n, err
		}
	}
}

// RecvMsg receives a message from the connection.
func (conn *SCTPConn) RecvMsg(b []byte, info *SCTPSndRcvInfo, flags *int) (n int, err error) {
	if !conn.ok() {
		return 0, syscall.EINVAL
	}
	oob := make([]byte, syscall.CmsgSpace(SCTPSndRcvInfoSize))
	n, noob, flag, _, err := syscall.Recvmsg(int(conn.sock), b, oob, 0)
	if err != nil {
		return n, err
	}
	*flags = flag
	if noob > 0 {
		ParseSndRcvInfo(info, oob[:noob])
	}
	return n, nil
}

// Write writes data to the connection.
func (conn *SCTPConn) Write(b []byte) (n int, err error) {
	return conn.SendMsg(b, nil)
}

// SendMsg sends a message on the connection.
func (conn *SCTPConn) SendMsg(b []byte, info *SCTPSndRcvInfo) (int, error) {
	if !conn.ok() {
		return 0, syscall.EINVAL
	}
	var buffer []byte
	if info != nil {
		hdr := &syscall.Cmsghdr{
			Level: syscall.IPPROTO_SCTP,
			Type:  SCTP_SNDRCV,
			Len:   uint64(syscall.CmsgSpace(SCTPSndRcvInfoSize)),
		}
		buffer = append(Pack(hdr), Pack(info)...)
	}
	// return syscall.SendmsgN(int(conn.sock), b, buffer, nil, 0)
	return SCTPSendMsg(int(conn.sock), b, buffer, 0)
}

// Abort aborts the SCTP association.
func (conn *SCTPConn) Abort() error {
	sock := atomic.SwapInt64(&conn.sock, -1)
	if sock > 0 {
		linger := syscall.Linger{
			Onoff:  1,
			Linger: 0,
		}
		_, _, _ = syscall.Syscall6(
			syscall.SYS_SETSOCKOPT,
			uintptr(sock),
			syscall.SOL_SOCKET,
			syscall.SO_LINGER,
			uintptr(unsafe.Pointer(&linger)),
			unsafe.Sizeof(linger),
			0,
		)
		return syscall.Close(int(sock))
	}
	return syscall.EBADFD
}

// Close closes the SCTP connection.
func (conn *SCTPConn) Close() error {
	if !conn.ok() {
		return syscall.EINVAL
	}
	msg := &SCTPSndRcvInfo{
		Flags: SCTP_EOF,
	}
	_, _ = conn.SendMsg(nil, msg)
	sock := atomic.SwapInt64(&conn.sock, -1)
	if sock > 0 {
		_ = syscall.Shutdown(int(sock), syscall.SHUT_RDWR)
		return syscall.Close(int(sock))
	}
	return syscall.EBADFD
}

// LocalAddr returns the local network address.
func (conn *SCTPConn) LocalAddr() net.Addr {
	if !conn.ok() {
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
		uintptr(conn.sock),
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

// RemoteAddr returns the remote network address.
func (conn *SCTPConn) RemoteAddr() net.Addr {
	if !conn.ok() {
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
		uintptr(conn.sock),
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

// SetDeadline sets the read and write deadlines. Not supported for SCTP.
func (conn *SCTPConn) SetDeadline(_ time.Time) error {
	return syscall.ENOPROTOOPT
}

// SetReadDeadline sets the read deadline. Not supported for SCTP.
func (conn *SCTPConn) SetReadDeadline(_ time.Time) error {
	return syscall.ENOPROTOOPT
}

// SetWriteDeadline sets the write deadline. Not supported for SCTP.
func (conn *SCTPConn) SetWriteDeadline(_ time.Time) error {
	return syscall.ENOPROTOOPT
}

// SetWriteBufferSize sets the size of the send buffer.
func (conn *SCTPConn) SetWriteBufferSize(bytes int) error {
	return syscall.SetsockoptInt(int(conn.sock), syscall.SOL_SOCKET, syscall.SO_SNDBUF, bytes)
}

// GetWriteBufferSize gets the size of the send buffer.
func (conn *SCTPConn) GetWriteBufferSize() (int, error) {
	return syscall.GetsockoptInt(int(conn.sock), syscall.SOL_SOCKET, syscall.SO_SNDBUF)
}

// SetReadBufferSize sets the size of the receive buffer.
func (conn *SCTPConn) SetReadBufferSize(bytes int) error {
	return syscall.SetsockoptInt(int(conn.sock), syscall.SOL_SOCKET, syscall.SO_RCVBUF, bytes)
}

// GetReadBufferSize gets the size of the receive buffer.
func (conn *SCTPConn) GetReadBufferSize() (int, error) {
	return syscall.GetsockoptInt(int(conn.sock), syscall.SOL_SOCKET, syscall.SO_RCVBUF)
}

// SetEventSubscribe sets the SCTP event subscriptions.
func (conn *SCTPConn) SetEventSubscribe(events *SCTPEventSubscribe) error {
	_, _, err := syscall.Syscall6(
		syscall.SYS_SETSOCKOPT,
		uintptr(conn.sock),
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

// GetEventSubscribe gets the current SCTP event subscriptions.
func (conn *SCTPConn) GetEventSubscribe() (*SCTPEventSubscribe, error) {
	if !conn.ok() {
		return nil, syscall.EINVAL
	}
	var (
		events SCTPEventSubscribe
		length = unsafe.Sizeof(events)
	)
	_, _, err := syscall.Syscall6(
		syscall.SYS_GETSOCKOPT,
		uintptr(conn.sock),
		SOL_SCTP,
		SCTP_EVENTS,
		uintptr(unsafe.Pointer(&events)),
		uintptr(unsafe.Pointer(&length)),
		0,
	)
	if err != 0 {
		return nil, err
	}
	return &events, nil
}

// SetInitMsg sets the SCTP initialization message parameters.
func (conn *SCTPConn) SetInitMsg(init *SCTPInitMsg) error {
	if !conn.ok() {
		return syscall.EINVAL
	}
	_, _, err := syscall.Syscall6(
		syscall.SYS_SETSOCKOPT,
		uintptr(conn.sock),
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

// GetInitMsg gets the current SCTP initialization message parameters.
func (conn *SCTPConn) GetInitMsg() (*SCTPInitMsg, error) {
	if !conn.ok() {
		return nil, syscall.EINVAL
	}
	var (
		init   SCTPInitMsg
		length = unsafe.Sizeof(init)
	)
	_, _, err := syscall.Syscall6(
		syscall.SYS_GETSOCKOPT,
		uintptr(conn.sock),
		SOL_SCTP,
		SCTP_INITMSG,
		uintptr(unsafe.Pointer(&init)),
		uintptr(unsafe.Pointer(&length)),
		0,
	)
	if err != 0 {
		return nil, err
	}
	return &init, nil
}

// SetDefaultSendParam sets the default send parameters for the association.
func (conn *SCTPConn) SetDefaultSendParam(param *SCTPSndRcvInfo) error {
	if !conn.ok() {
		return syscall.EINVAL
	}
	_, _, err := syscall.Syscall6(
		syscall.SYS_SETSOCKOPT,
		uintptr(conn.sock),
		SOL_SCTP,
		SCTP_DEFAULT_SEND_PARAM,
		uintptr(unsafe.Pointer(param)),
		unsafe.Sizeof(*param),
		0,
	)
	if err != 0 {
		return err
	}
	return nil
}

// GetDefaultSendParam gets the default send parameters for the association.
func (conn *SCTPConn) GetDefaultSendParam() (*SCTPSndRcvInfo, error) {
	if !conn.ok() {
		return nil, syscall.EINVAL
	}
	var (
		param  SCTPSndRcvInfo
		length = unsafe.Sizeof(param)
	)
	_, _, err := syscall.Syscall6(
		syscall.SYS_GETSOCKOPT,
		uintptr(conn.sock),
		SOL_SCTP,
		SCTP_DEFAULT_SEND_PARAM,
		uintptr(unsafe.Pointer(&param)),
		uintptr(unsafe.Pointer(&length)),
		0,
	)
	if err != 0 {
		return nil, err
	}
	return &param, nil
}

func (conn *SCTPConn) ok() bool {
	if nil != conn && conn.sock > 0 {
		return true
	}
	return false
}

// DialSCTP dials an SCTP connection to the remote address.
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
	if remote == nil {
		return nil, &net.OpError{
			Op:     "dial",
			Net:    network,
			Source: local.Addr(),
			Addr:   remote.Addr(),
			Err:    net.InvalidAddrError("invalid remote addr"),
		}
	}
	// syscall.SOCK_SEQPACKET vs syscall.SOCK_STREAM
	sock, err := SCTPSocket(AddrFamily(network), syscall.SOCK_STREAM)
	if err != nil {
		return nil, err
	}
	conn := &SCTPConn{
		sock: int64(sock),
	}
	if err = syscall.SetsockoptInt(sock, syscall.SOL_SOCKET, syscall.SO_BROADCAST, 1); err != nil {
		_ = conn.Close()
		return nil, err
	}
	if err = conn.SetInitMsg(init); err != nil {
		_ = conn.Close()
		return nil, err
	}
	if local != nil {
		if err = SCTPBind(sock, local, SCTP_BINDX_ADD_ADDR); err != nil {
			_ = conn.Close()
			return nil, err
		}
	}
	conn.assoc, err = SCTPConnect(sock, remote)
	if err != nil {
		_ = conn.Close()
		return nil, err
	}
	return conn, nil
}
