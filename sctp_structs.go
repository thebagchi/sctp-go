package sctp_go

type IoVector struct {
	Base *byte
	Len  uint64
}
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
type CMsgHeader struct {
	Len   uint64
	Level int32
	Type  int32
}
type InAddr struct {
	Addr [4]byte
}
type In6Addr struct {
	Addr [16]byte
}
type SockAddrIn6 struct {
	Family   uint16
	Port     uint16
	FlowInfo uint32
	Addr     In6Addr
	ScopeId  uint32
}
type SockAddrIn struct {
	Family uint16
	Port   uint16
	Addr   InAddr
	Zero   [8]uint8
}
type SockAddr struct {
	Family uint16
	Data   [14]int8
}
type SockAddrStorage struct {
	Family  uint16
	Padding [118]int8
	Align   uint64
}
type SCTPAssocId int32
type SCTPInitMsg struct {
	NumOutStreams  uint16
	MaxInStreams   uint16
	MaxAttempts    uint16
	MaxInitTimeout uint16
}
type SCTPSndRcvInfo struct {
	Stream     uint16
	Ssn        uint16
	Flags      uint16
	Ppid       uint32
	Context    uint32
	TimeToLive uint32
	Tsn        uint32
	CumTsn     uint32
	AssocId    int32
}
type SCTPGetAddrsOld struct {
	AssocId int32
	Num     int32
	Addrs   *SockAddr
}
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
type SCTPSetPeerPrimary struct {
	AssocId int32
	Addr    [128]byte
	// Cannot use SockAddrStorage as type for Addr.
	// Inner structures have alignment requirement of 8 bytes.
	// This structure has alignment requirement of 4 bytes.
}
type SCTPPrimaryAddr SCTPSetPeerPrimary
