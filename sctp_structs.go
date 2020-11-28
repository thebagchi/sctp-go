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
