package sctp_go

import (
	"bytes"
	"encoding/binary"
	"unsafe"
)

// IoVector represents the C struct iovec for scatter/gather I/O operations.
type IoVector struct {
	Base *byte
	Len  uint64
}

// MsgHeader represents the C struct msghdr for socket message headers.
type MsgHeader struct {
	Name       *byte
	NameLen    uint32
	Iov        *IoVector
	IovLen     uint64
	Control    *byte
	ControlLen uint64
	Flags      int32
	Padding    [4]byte
}

// CMsgHeader represents the C struct cmsghdr for control message headers.
type CMsgHeader struct {
	Len   uint64
	Level int32
	Type  int32
}

// InAddr represents the C struct in_addr for IPv4 addresses.
type InAddr struct {
	Addr [4]byte
}

// In6Addr represents the C struct in6_addr for IPv6 addresses.
type In6Addr struct {
	Addr [16]byte
}

// SockAddrIn6 represents the C struct sockaddr_in6 for IPv6 socket addresses.
type SockAddrIn6 struct {
	Family   uint16
	Port     uint16
	FlowInfo uint32
	Addr     In6Addr
	ScopeId  uint32
}

// SockAddrIn represents the C struct sockaddr_in for IPv4 socket addresses.
type SockAddrIn struct {
	Family uint16
	Port   uint16
	Addr   InAddr
	Zero   [8]uint8
}

// SockAddr represents the C struct sockaddr for generic socket addresses.
type SockAddr struct {
	Family uint16
	Data   [14]int8
}

// SockAddrStorage represents the C struct sockaddr_storage for storing socket addresses.
type SockAddrStorage struct {
	Family  uint16
	Padding [118]int8
	Align   uint64
}

// SCTPAssocId represents the C type sctp_assoc_t for SCTP association IDs.
type SCTPAssocId int32

// SCTPInitMsg represents the C struct sctp_initmsg for SCTP initialization parameters.
type SCTPInitMsg struct {
	NumOutStreams  uint16
	MaxInStreams   uint16
	MaxAttempts    uint16
	MaxInitTimeout uint16
}

// SCTPSndRcvInfo represents the C struct sctp_sndrcvinfo for SCTP send/receive information.
type SCTPSndRcvInfo struct {
	Stream     uint16
	Ssn        uint16
	Flags      uint16
	_          uint16
	Ppid       uint32
	Context    uint32
	TimeToLive uint32
	Tsn        uint32
	CumTsn     uint32
	AssocId    int32
}

// SCTPSndInfo represents the C struct sctp_sndinfo for SCTP send information.
type SCTPSndInfo struct {
	Sid     uint16
	Flags   uint16
	Ppid    uint32
	Context uint32
	AssocId int32
}

// SCTPRcvInfo represents the C struct sctp_rcvinfo for SCTP receive information.
type SCTPRcvInfo struct {
	Sid     uint16
	Ssn     uint16
	Flags   uint16
	Ppid    uint32
	Tsn     uint32
	CumTsn  uint32
	Context uint32
	AssocId int32
}

// SCTPNxtInfo represents the C struct sctp_nxtinfo for SCTP next information.
type SCTPNxtInfo struct {
	Sid     uint16
	Flags   uint16
	Ppid    uint32
	Length  uint32
	AssocId int32
}

// SCTPPrInfo represents the C struct sctp_prinfo for SCTP partial reliability information.
type SCTPPrInfo struct {
	Policy uint16
	Value  uint32
}

// SCTPAuthInfo represents the C struct sctp_authinfo for SCTP authentication information.
type SCTPAuthInfo struct {
	KeyNumber uint16
}

// SCTPCmsgData represents the C type sctp_cmsg_data_t for SCTP control message data.
type SCTPCmsgData [32]byte

// SCTPGetAddrsOld represents the C struct sctp_getaddrs_old for getting SCTP addresses (old version).
type SCTPGetAddrsOld struct {
	AssocId int32
	Num     int32
	Addrs   uintptr
}

// SCTPGetAddrs represents the C struct sctp_getaddrs for getting SCTP addresses.
type SCTPGetAddrs struct {
	AssocId int32
	Num     uint32
	// Member Addr is removed as it is variable sized array.
}

// SCTPEventSubscribe represents the C struct sctp_event_subscribe for SCTP event subscription.
type SCTPEventSubscribe struct {
	DataIoEvent          uint8
	AssociationEvent     uint8
	AddressEvent         uint8
	SendFailureEvent     uint8
	PeerErrorEvent       uint8
	ShutdownEvent        uint8
	PartialDeliveryEvent uint8
	AdaptationLayerEvent uint8
	AuthenticationEvent  uint8
	SenderDryEvent       uint8
	StreamResetEvent     uint8
	AssocResetEvent      uint8
	StreamChangeEvent    uint8
}

// SCTPSetPeerPrimary represents the C struct sctp_setpeerprim for setting SCTP peer primary address.
type SCTPSetPeerPrimary struct {
	AssocId int32
	Addr    [128]byte
	// Cannot use SockAddrStorage as type for Addr.
	// Inner structures have alignment requirement of 8 bytes.
	// This structure has alignment requirement of 4 bytes.
}

// SCTPPrimaryAddr represents the C struct sctp_prim for SCTP primary address.
type SCTPPrimaryAddr SCTPSetPeerPrimary

// SCTPPeelOffArg represents the C type sctp_peeloff_arg_t for SCTP peel-off arguments.
type SCTPPeelOffArg struct {
	AssocId int32
	Sd      int32
}

// SCTPPeelOffFlagsArg represents the C type sctp_peeloff_flags_arg_t for SCTP peel-off arguments with flags.
type SCTPPeelOffFlagsArg struct {
	Arg   SCTPPeelOffArg
	Flags uint32
}

// SCTPNotification represents the C union sctp_notification for SCTP notifications.
type SCTPNotification [148]byte

// Notification represents the interface for SCTP notification types.
type Notification interface {
	GetType() uint16
	GetFlags() uint16
	GetLength() uint32
}

// SCTPNotificationHeader represents the C struct sn_header for SCTP notification headers.
type SCTPNotificationHeader struct {
	Type   uint16
	Flags  uint16
	Length uint32
}

// GetType returns the type field of the SCTP notification header.
func (n *SCTPNotificationHeader) GetType() uint16 {
	return n.Type
}

// GetFlags returns the flags field of the SCTP notification header.
func (n *SCTPNotificationHeader) GetFlags() uint16 {
	return n.Flags
}

// GetLength returns the length field of the SCTP notification header.
func (n *SCTPNotificationHeader) GetLength() uint32 {
	return n.Length
}

// SCTPAssocChange represents the C struct sctp_assoc_change for SCTP association change notifications.
type SCTPAssocChange struct {
	Type            uint16
	Flags           uint16
	Length          uint32
	State           uint16
	Error           uint16
	OutboundStreams uint16
	InboundStreams  uint16
	AssocId         int32
}

// GetType returns the type field of the SCTP association change notification.
func (n *SCTPAssocChange) GetType() uint16 {
	return n.Type
}

// GetFlags returns the flags field of the SCTP association change notification.
func (n *SCTPAssocChange) GetFlags() uint16 {
	return n.Flags
}

// GetLength returns the length field of the SCTP association change notification.
func (n *SCTPAssocChange) GetLength() uint32 {
	return n.Length
}

// SCTPPAddrChange represents the C struct sctp_paddr_change for SCTP peer address change notifications.
type SCTPPAddrChange struct {
	Type    uint16
	Flags   uint16
	Length  uint32
	Addr    [128]byte
	State   int32
	Error   int32
	AssocId int32
	// Cannot use SockAddrStorage as type for Addr.
	// Inner structures have alignment requirement of 8 bytes.
	// This structure has alignment requirement of 4 bytes.
}

// GetType returns the type field of the SCTP peer address change notification.
func (n *SCTPPAddrChange) GetType() uint16 {
	return n.Type
}

// GetFlags returns the flags field of the SCTP peer address change notification.
func (n *SCTPPAddrChange) GetFlags() uint16 {
	return n.Flags
}

// GetLength returns the length field of the SCTP peer address change notification.
func (n *SCTPPAddrChange) GetLength() uint32 {
	return n.Length
}

// GetAddr returns the peer address from the SCTP peer address change notification.
func (n *SCTPPAddrChange) GetAddr() *SCTPAddr {
	return FromSockAddrStorage((*SockAddrStorage)(unsafe.Pointer(&n.Addr)))
}

// SCTPRemoteError represents the C struct sctp_remote_error for SCTP remote error notifications.
type SCTPRemoteError struct {
	Type    uint16
	Flags   uint16
	Length  uint32
	Error   uint16
	AssocId int32
}

// GetType returns the type field of the SCTP remote error notification.
func (n *SCTPRemoteError) GetType() uint16 {
	return n.Type
}

// GetFlags returns the flags field of the SCTP remote error notification.
func (n *SCTPRemoteError) GetFlags() uint16 {
	return n.Flags
}

// GetLength returns the length field of the SCTP remote error notification.
func (n *SCTPRemoteError) GetLength() uint32 {
	return n.Length
}

// SCTPSendFailed represents the C struct sctp_send_failed for SCTP send failure notifications.
type SCTPSendFailed struct {
	Type    uint16
	Flags   uint16
	Length  uint32
	Error   uint32
	Info    SCTPSndRcvInfo
	AssocId int32
}

// GetType returns the type field of the SCTP send failed notification.
func (n *SCTPSendFailed) GetType() uint16 {
	return n.Type
}

// GetFlags returns the flags field of the SCTP send failed notification.
func (n *SCTPSendFailed) GetFlags() uint16 {
	return n.Flags
}

// GetLength returns the length field of the SCTP send failed notification.
func (n *SCTPSendFailed) GetLength() uint32 {
	return n.Length
}

// SCTPShutdownEvent represents the C struct sctp_shutdown_event for SCTP shutdown event notifications.
type SCTPShutdownEvent struct {
	Type    uint16
	Flags   uint16
	Length  uint32
	AssocId int32
}

// GetType returns the type field of the SCTP shutdown event notification.
func (n *SCTPShutdownEvent) GetType() uint16 {
	return n.Type
}

// GetFlags returns the flags field of the SCTP shutdown event notification.
func (n *SCTPShutdownEvent) GetFlags() uint16 {
	return n.Flags
}

// GetLength returns the length field of the SCTP shutdown event notification.
func (n *SCTPShutdownEvent) GetLength() uint32 {
	return n.Length
}

// SCTPAdaptationEvent represents the C struct sctp_adaptation_event for SCTP adaptation layer event notifications.
type SCTPAdaptationEvent struct {
	Type          uint16
	Flags         uint16
	Length        uint32
	AdaptationInd uint32
	AssocId       int32
}

// GetType returns the type field of the SCTP adaptation event notification.
func (n *SCTPAdaptationEvent) GetType() uint16 {
	return n.Type
}

// GetFlags returns the flags field of the SCTP adaptation event notification.
func (n *SCTPAdaptationEvent) GetFlags() uint16 {
	return n.Flags
}

// GetLength returns the length field of the SCTP adaptation event notification.
func (n *SCTPAdaptationEvent) GetLength() uint32 {
	return n.Length
}

// SCTPPDApiEvent represents the C struct sctp_pdapi_event for SCTP partial delivery API event notifications.
type SCTPPDApiEvent struct {
	Type       uint16
	Flags      uint16
	Length     uint32
	Indication uint32
	AssocId    int32
	Stream     uint32
	Sequence   uint32
}

// GetType returns the type field of the SCTP partial delivery API event notification.
func (n *SCTPPDApiEvent) GetType() uint16 {
	return n.Type
}

// GetFlags returns the flags field of the SCTP partial delivery API event notification.
func (n *SCTPPDApiEvent) GetFlags() uint16 {
	return n.Flags
}

// GetLength returns the length field of the SCTP partial delivery API event notification.
func (n *SCTPPDApiEvent) GetLength() uint32 {
	return n.Length
}

// SCTPAuthKeyEvent represents the C struct sctp_authkey_event for SCTP authentication key event notifications.
type SCTPAuthKeyEvent struct {
	Type         uint16
	Flags        uint16
	Length       uint32
	KeyNumber    uint16
	AltKeyNumber uint16
	Indication   uint32
	AssocId      int32
}

// GetType returns the type field of the SCTP authentication key event notification.
func (n *SCTPAuthKeyEvent) GetType() uint16 {
	return n.Type
}

// GetFlags returns the flags field of the SCTP authentication key event notification.
func (n *SCTPAuthKeyEvent) GetFlags() uint16 {
	return n.Flags
}

// GetLength returns the length field of the SCTP authentication key event notification.
func (n *SCTPAuthKeyEvent) GetLength() uint32 {
	return n.Length
}

// SCTPSenderDryEvent represents the C struct sctp_sender_dry_event for SCTP sender dry event notifications.
type SCTPSenderDryEvent struct {
	Type    uint16
	Flags   uint16
	Length  uint32
	AssocId int32
}

// GetType returns the type field of the SCTP sender dry event notification.
func (n *SCTPSenderDryEvent) GetType() uint16 {
	return n.Type
}

// GetFlags returns the flags field of the SCTP sender dry event notification.
func (n *SCTPSenderDryEvent) GetFlags() uint16 {
	return n.Flags
}

// GetLength returns the length field of the SCTP sender dry event notification.
func (n *SCTPSenderDryEvent) GetLength() uint32 {
	return n.Length
}

// SCTPStreamResetEvent represents the C struct sctp_stream_reset_event for SCTP stream reset event notifications.
type SCTPStreamResetEvent struct {
	Type    uint16
	Flags   uint16
	Length  uint32
	AssocId int32
}

// GetType returns the type field of the SCTP stream reset event notification.
func (n *SCTPStreamResetEvent) GetType() uint16 {
	return n.Type
}

// GetFlags returns the flags field of the SCTP stream reset event notification.
func (n *SCTPStreamResetEvent) GetFlags() uint16 {
	return n.Flags
}

// GetLength returns the length field of the SCTP stream reset event notification.
func (n *SCTPStreamResetEvent) GetLength() uint32 {
	return n.Length
}

// SCTPAssocResetEvent represents the C struct sctp_assoc_reset_event for SCTP association reset event notifications.
type SCTPAssocResetEvent struct {
	Type      uint16
	Flags     uint16
	Length    uint32
	AssocId   int32
	LocalTsn  uint32
	RemoteTsn uint32
}

// GetType returns the type field of the SCTP association reset event notification.
func (n *SCTPAssocResetEvent) GetType() uint16 {
	return n.Type
}

// GetFlags returns the flags field of the SCTP association reset event notification.
func (n *SCTPAssocResetEvent) GetFlags() uint16 {
	return n.Flags
}

// GetLength returns the length field of the SCTP association reset event notification.
func (n *SCTPAssocResetEvent) GetLength() uint32 {
	return n.Length
}

// SCTPStreamChangeEvent represents the C struct sctp_stream_change_event for SCTP stream change event notifications.
type SCTPStreamChangeEvent struct {
	Type       uint16
	Flags      uint16
	Length     uint32
	AssocId    int32
	InStreams  uint16
	OutStreams uint16
}

// GetType returns the type field of the SCTP stream change event notification.
func (n *SCTPStreamChangeEvent) GetType() uint16 {
	return n.Type
}

// GetFlags returns the flags field of the SCTP stream change event notification.
func (n *SCTPStreamChangeEvent) GetFlags() uint16 {
	return n.Flags
}

// GetLength returns the length field of the SCTP stream change event notification.
func (n *SCTPStreamChangeEvent) GetLength() uint32 {
	return n.Length
}

// SCTPRTOInfo represents the C struct sctp_rtoinfo for SCTP retransmission timeout information.
type SCTPRTOInfo struct {
	AssocId int32
	Initial uint32
	Max     uint32
	Min     uint32
}

// SCTPResetStreams represents the C struct sctp_reset_streams for SCTP reset streams parameters.
type SCTPResetStreams struct {
	AssocId       int32
	Flags         uint16
	NumberStreams uint16
}

// SCTPAddStreams represents the C struct sctp_add_streams for SCTP add streams parameters.
type SCTPAddStreams struct {
	AssocId    int32
	InStreams  uint16
	OutStreams uint16
}

// SCTPAssocParams represents the C struct sctp_assocparams for SCTP association parameters.
type SCTPAssocParams struct {
	AssocId                int32
	AssocMaxrxt            uint16
	NumberPeerDestinations uint16
	PeerRwnd               uint32
	LocalRwnd              uint32
	CookieLife             uint32
}

// SCTPSetAdaptation represents the C struct sctp_setadaptation for SCTP adaptation layer parameters.
type SCTPSetAdaptation struct {
	AdaptationInd uint32
}

// SCTPPeerAddrParams represents the C struct sctp_paddrparams for SCTP peer address parameters.
type SCTPPeerAddrParams struct {
	AssocId       int32
	Addr          [128]byte
	HbInterval    uint32
	PathMaxRxt    uint16
	PathMtu       uint32
	SackDelay     uint32
	Flags         uint32
	Ipv6FlowLabel uint32
	Dscp          uint8
	_             uint8
}

// Pack serializes the SCTPPeerAddrParams struct into a byte slice.
func (s *SCTPPeerAddrParams) Pack() []byte {
	buffer := bytes.NewBuffer(make([]byte, 0, SCTPPeerAddrParamsSize))
	_ = binary.Write(buffer, endian, s.AssocId)
	_ = binary.Write(buffer, endian, s.Addr)
	_ = binary.Write(buffer, endian, s.HbInterval)
	_ = binary.Write(buffer, endian, s.PathMaxRxt)
	_ = binary.Write(buffer, endian, s.PathMtu)
	_ = binary.Write(buffer, endian, s.SackDelay)
	_ = binary.Write(buffer, endian, s.Flags)
	_ = binary.Write(buffer, endian, s.Ipv6FlowLabel)
	_ = binary.Write(buffer, endian, s.Dscp)
	_ = binary.Write(buffer, endian, uint8(0))
	return buffer.Bytes()
}

// Unpack deserializes a byte slice into the SCTPPeerAddrParams struct.
func (s *SCTPPeerAddrParams) Unpack(data []byte) {
	if len(data) == SCTPPeerAddrParamsSize {
		buffer := bytes.NewReader(data)
		_ = binary.Read(buffer, endian, s.AssocId)
		_ = binary.Read(buffer, endian, s.Addr)
		_ = binary.Read(buffer, endian, s.HbInterval)
		_ = binary.Read(buffer, endian, s.PathMaxRxt)
		_ = binary.Read(buffer, endian, s.PathMtu)
		_ = binary.Read(buffer, endian, s.SackDelay)
		_ = binary.Read(buffer, endian, s.Flags)
		_ = binary.Read(buffer, endian, s.Ipv6FlowLabel)
		_ = binary.Read(buffer, endian, s.Dscp)
	}
}

// SCTPPeerAddrInfo represents the C struct sctp_paddrinfo for SCTP peer address information.
type SCTPPeerAddrInfo struct {
	AssocId int32
	Addr    [128]byte
	State   int32
	Cwnd    uint32
	Srtt    uint32
	Rto     uint32
	Mtu     uint32
}

// SCTPAssocValue represents the C struct sctp_assoc_value for SCTP association value parameters.
type SCTPAssocValue struct {
	Id    int32
	Value uint32
}

// SCTPSackInfo represents the C struct sctp_sack_info for SCTP selective acknowledgment information.
type SCTPSackInfo struct {
	AssocId int32
	Delay   uint32
	Freq    uint32
}

// SCTPStreamValue represents the C struct sctp_stream_value for SCTP stream value parameters.
type SCTPStreamValue struct {
	AssocId     int32
	StreamId    uint16
	StreamValue uint16
}

// SCTPStatus represents the C struct sctp_status for SCTP association status.
type SCTPStatus struct {
	AssocId            int32
	State              int32
	Rwnd               uint32
	UnackedData        uint16
	PendingData        uint16
	InStreams          uint16
	OutStreams         uint16
	FragmentationPoint uint32
	Primary            SCTPPeerAddrInfo
}

// SCTPAuthKeyId represents the C struct sctp_authkeyid for SCTP authentication key ID.
type SCTPAuthKeyId struct {
	AssocId   int32
	KeyNumber uint16
	_         uint16
}

// SCTPAuthKey represents the C struct sctp_authkey for SCTP authentication key.
type SCTPAuthKey struct {
	AssocId   int32
	KeyNumber uint16
	KeyLength uint16
	// Member Key is removed as it is variable sized array.
}

// SCTPAuthChunk represents the C struct sctp_authchunk for SCTP authentication chunk.
type SCTPAuthChunk struct {
	Chunk uint8
}

// SCTPHmacAlgo represents the C struct sctp_hmacalgo for SCTP HMAC algorithm.
type SCTPHmacAlgo struct {
	NumIdents uint32
	// Member Idents is removed as it is variable sized array.
}

// SCTPAuthChunks represents the C struct sctp_authchunks for SCTP authentication chunks.
type SCTPAuthChunks struct {
	AssocId      int32
	NumberChunks uint32
	// Member Chunks is removed as it is variable sized array.
}

// SCTPAssocIds represents the C struct sctp_assoc_ids for SCTP association IDs.
type SCTPAssocIds struct {
	NumberIds uint32
	// Member Ids is removed as it is variable sized array.
}

// SCTPAssocStats represents the C struct sctp_assoc_stats for SCTP association statistics.
type SCTPAssocStats struct {
	AssocId      int32
	Addr         SockAddrStorage
	MaxRto       uint64
	ISacks       uint64
	OSacks       uint64
	OPackets     uint64
	IPackets     uint64
	RtxChunks    uint64
	OutOfSeqTsns uint64
	IDupChunks   uint64
	GapCnt       uint64
	OUodChunks   uint64
	IUodChunks   uint64
	OodChunks    uint64
	IodChunks    uint64
	OCtrlChunks  uint64
	ICtrlChunks  uint64
}

// SCTPPeerAddrThresholds represents the C struct sctp_paddrthlds for SCTP peer address thresholds.
type SCTPPeerAddrThresholds struct {
	AssocId         int32
	Address         SockAddrStorage
	PathMaxRxt      uint16
	PathPfThreshold uint16
}

// SCTPPRStatus represents the C struct sctp_prstatus for SCTP partial reliability status.
type SCTPPRStatus struct {
	AssocId         int32
	Sid             uint16
	Policy          uint16
	AbandonedUnsent uint64
	AbandonedSent   uint64
}

// SCTPDefaultPRInfo represents the C struct sctp_default_prinfo for SCTP default partial reliability information.
type SCTPDefaultPRInfo struct {
	AssocId int32
	Value   uint32
	Policy  uint16
}

// SCTPEvent represents the C struct sctp_event for SCTP event parameters.
type SCTPEvent struct {
	AssocId int32
	Type    uint16
	On      uint8
}

// SCTPInfo represents the C struct sctp_info for SCTP association information.
//
//lint:ignore U1000 "auto-generated"
type SCTPInfo struct {
	Tag                      uint32
	State                    uint32
	Rwnd                     uint32
	UnackedData              uint16
	PendingData              uint16
	InStreams                uint16
	OutStreams               uint16
	FragmentationPoint       uint32
	InQueue                  uint32
	OutQueue                 uint32
	OverallError             uint32
	MaxBurst                 uint32
	MaxSeg                   uint32
	PeerRwnd                 uint32
	PeerTag                  uint32
	PeerCapable              uint8
	PeerSack                 uint8
	_1                       uint16
	ISacks                   uint64
	OSacks                   uint64
	OPackets                 uint64
	IPackets                 uint64
	RtxChunks                uint64
	OutOfSeqTsns             uint64
	IDupChunks               uint64
	GapCount                 uint64
	OUodChunks               uint64
	IUodChunks               uint64
	OOdChunks                uint64
	IOdChunks                uint64
	OCtrlChunks              uint64
	ICtrlChunks              uint64
	PrimaryAddress           SockAddrStorage
	PrimaryState             int32
	PrimaryCwnd              uint32
	PrimarySrtt              uint32
	PrimaryRto               uint32
	PrimaryHbInterval        uint32
	PrimaryPathMaxRxt        uint32
	PrimarySackDelay         uint32
	PrimarySackFreq          uint32
	PrimarySsThreshold       uint32
	PrimaryPartialBytesAcked uint32
	PrimaryFlightSize        uint32
	PrimaryError             uint16
	_2                       uint16
	SockAutoClose            uint32
	SockAdaptationInd        uint32
	SockPdPoint              uint32
	SockNodelay              uint8
	SockDisableFragments     uint8
	SockV4Mapped             uint8
	SockFragInterleave       uint8
	SockType                 uint32
	_3                       uint32
}
