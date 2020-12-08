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
type SCTPSndInfo struct {
	Sid     uint16
	Flags   uint16
	Ppid    uint32
	Context uint32
	AssocId int32
}
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
type SCTPNxtInfo struct {
	Sid     uint16
	Flags   uint16
	Ppid    uint32
	Length  uint32
	AssocId int32
}
type SCTPPrInfo struct {
	Policy uint16
	Value  uint32
}
type SCTPAuthInfo struct {
	KeyNumber uint16
}
type SCTPCmsgData [32]byte
type SCTPGetAddrsOld struct {
	AssocId int32
	Num     int32
	Addrs   uintptr
}
type SCTPGetAddrs struct {
	AssocId int32
	Num     uint32
	// Member Addr is removed as it variable sized array.
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
type SCTPPeelOffArg struct {
	AssocId int32
	Sd      int32
}
type SCTPPeelOffFlagsArg struct {
	Arg   SCTPPeelOffArg
	Flags uint32
}
type Notification interface {
	GetType() uint16
	GetFlags() uint16
	GetLength() uint32
}
type SCTPNotificationHeader struct {
	Type   uint16
	Flags  uint16
	Length uint32
}

func (n *SCTPNotificationHeader) GetType() uint16 {
	return n.Type
}
func (n *SCTPNotificationHeader) GetFlags() uint16 {
	return n.Flags
}
func (n *SCTPNotificationHeader) GetLength() uint32 {
	return n.Length
}

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

func (n *SCTPAssocChange) GetType() uint16 {
	return n.Type
}
func (n *SCTPAssocChange) GetFlags() uint16 {
	return n.Flags
}
func (n *SCTPAssocChange) GetLength() uint32 {
	return n.Length
}

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

func (n *SCTPPAddrChange) GetType() uint16 {
	return n.Type
}
func (n *SCTPPAddrChange) GetFlags() uint16 {
	return n.Flags
}
func (n *SCTPPAddrChange) GetLength() uint32 {
	return n.Length
}
func (n *SCTPPAddrChange) GetAddr() *SCTPAddr {
	return nil
}

type SCTPRemoteError struct {
	Type    uint16
	Flags   uint16
	Length  uint32
	Error   uint16
	AssocId int32
}

func (n *SCTPRemoteError) GetType() uint16 {
	return n.Type
}
func (n *SCTPRemoteError) GetFlags() uint16 {
	return n.Flags
}
func (n *SCTPRemoteError) GetLength() uint32 {
	return n.Length
}

type SCTPSendFailed struct {
	Type    uint16
	Flags   uint16
	Length  uint32
	Error   uint32
	Info    SCTPSndRcvInfo
	AssocId int32
}

func (n *SCTPSendFailed) GetType() uint16 {
	return n.Type
}
func (n *SCTPSendFailed) GetFlags() uint16 {
	return n.Flags
}
func (n *SCTPSendFailed) GetLength() uint32 {
	return n.Length
}

type SCTPShutdownEvent struct {
	Type    uint16
	Flags   uint16
	Length  uint32
	AssocId int32
}

func (n *SCTPShutdownEvent) GetType() uint16 {
	return n.Type
}
func (n *SCTPShutdownEvent) GetFlags() uint16 {
	return n.Flags
}
func (n *SCTPShutdownEvent) GetLength() uint32 {
	return n.Length
}

type SCTPAdaptationEvent struct {
	Type          uint16
	Flags         uint16
	Length        uint32
	AdaptationInd uint32
	AssocId       int32
}

func (n *SCTPAdaptationEvent) GetType() uint16 {
	return n.Type
}
func (n *SCTPAdaptationEvent) GetFlags() uint16 {
	return n.Flags
}
func (n *SCTPAdaptationEvent) GetLength() uint32 {
	return n.Length
}

type SCTPPDApiEvent struct {
	Type       uint16
	Flags      uint16
	Length     uint32
	Indication uint32
	AssocId    int32
	Stream     uint32
	Sequence   uint32
}

func (n *SCTPPDApiEvent) GetType() uint16 {
	return n.Type
}
func (n *SCTPPDApiEvent) GetFlags() uint16 {
	return n.Flags
}
func (n *SCTPPDApiEvent) GetLength() uint32 {
	return n.Length
}

type SCTPAuthKeyEvent struct {
	Type         uint16
	Flags        uint16
	Length       uint32
	KeyNumber    uint16
	AltKeyNumber uint16
	Indication   uint32
	AssocId      int32
}

func (n *SCTPAuthKeyEvent) GetType() uint16 {
	return n.Type
}
func (n *SCTPAuthKeyEvent) GetFlags() uint16 {
	return n.Flags
}
func (n *SCTPAuthKeyEvent) GetLength() uint32 {
	return n.Length
}

type SCTPSenderDryEvent struct {
	Type    uint16
	Flags   uint16
	Length  uint32
	AssocId int32
}

func (n *SCTPSenderDryEvent) GetType() uint16 {
	return n.Type
}
func (n *SCTPSenderDryEvent) GetFlags() uint16 {
	return n.Flags
}
func (n *SCTPSenderDryEvent) GetLength() uint32 {
	return n.Length
}

type SCTPStreamResetEvent struct {
	Type    uint16
	Flags   uint16
	Length  uint32
	AssocId int32
}

func (n *SCTPStreamResetEvent) GetType() uint16 {
	return n.Type
}
func (n *SCTPStreamResetEvent) GetFlags() uint16 {
	return n.Flags
}
func (n *SCTPStreamResetEvent) GetLength() uint32 {
	return n.Length
}

type SCTPAssocResetEvent struct {
	Type      uint16
	Flags     uint16
	Length    uint32
	AssocId   int32
	LocalTsn  uint32
	RemoteTsn uint32
}

func (n *SCTPAssocResetEvent) GetType() uint16 {
	return n.Type
}
func (n *SCTPAssocResetEvent) GetFlags() uint16 {
	return n.Flags
}
func (n *SCTPAssocResetEvent) GetLength() uint32 {
	return n.Length
}

type SCTPStreamChangeEvent struct {
	Type       uint16
	Flags      uint16
	Length     uint32
	AssocId    int32
	InStreams  uint16
	OutStreams uint16
}

func (n *SCTPStreamChangeEvent) GetType() uint16 {
	return n.Type
}
func (n *SCTPStreamChangeEvent) GetFlags() uint16 {
	return n.Flags
}
func (n *SCTPStreamChangeEvent) GetLength() uint32 {
	return n.Length
}
