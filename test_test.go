package sctp_go

import (
	"encoding/hex"
	"fmt"
	"syscall"
	"testing"
	"unsafe"
)

func TestSizes(t *testing.T) {
	{
		temp := &IoVector{}
		if unsafe.Sizeof(*temp) != IoVectorSize {
			fmt.Println(unsafe.Sizeof(*temp))
			fmt.Println(IoVectorSize)
			t.Error("IoVector sizes don't match")
		} else {
			fmt.Println("IoVectorSize: ", IoVectorSize)
		}
	}
	{
		temp := &MsgHeader{}
		if unsafe.Sizeof(*temp) != MsgHeaderSize {
			fmt.Println(unsafe.Sizeof(*temp))
			fmt.Println(MsgHeaderSize)
			t.Error("MsgHeader sizes don't match")
		} else {
			fmt.Println("MsgHeaderSize: ", MsgHeaderSize)
		}
	}
	{
		temp := &CMsgHeader{}
		if unsafe.Sizeof(*temp) != CMsgHeaderSize {
			fmt.Println(unsafe.Sizeof(*temp))
			fmt.Println(CMsgHeaderSize)
			t.Error("CMsgHeader sizes don't match")
		} else {
			fmt.Println("CMsgHeaderSize: ", CMsgHeaderSize)
		}
	}
	{
		temp := &SCTPSndRcvInfo{}
		if unsafe.Sizeof(*temp) != SCTPSndRcvInfoSize {
			fmt.Println(unsafe.Sizeof(*temp))
			fmt.Println(SCTPSndRcvInfoSize)
			t.Error("SCTPSndRcvInfo sizes don't match")
		} else {
			fmt.Println("SCTPSndRcvInfoSize: ", SCTPSndRcvInfoSize)
		}
	}
	{
		temp := &SCTPInitMsg{}
		if unsafe.Sizeof(*temp) != SCTPInitMsgSize {
			fmt.Println(unsafe.Sizeof(*temp))
			fmt.Println(SCTPInitMsgSize)
			t.Error("SCTPInitMsg sizes don't match")
		} else {
			fmt.Println("SCTPInitMsgSize: ", SCTPInitMsgSize)
		}
	}
	{
		temp := &SCTPGetAddrsOld{}
		if unsafe.Sizeof(*temp) != SCTPGetAddrsOldSize {
			fmt.Println(unsafe.Sizeof(*temp))
			fmt.Println(SCTPGetAddrsOldSize)
			t.Error("SCTPGetAddrsOld sizes don't match")
		} else {
			fmt.Println("SCTPGetAddrsOldSize: ", SCTPGetAddrsOldSize)
		}
	}
	{
		temp := &SCTPEventSubscribe{}
		if unsafe.Sizeof(*temp) != SCTPEventSubscribeSize {
			fmt.Println(unsafe.Sizeof(*temp))
			fmt.Println(SCTPEventSubscribeSize)
			t.Error("SCTPEventSubscribe sizes don't match")
		} else {
			fmt.Println("SCTPEventSubscribeSize: ", SCTPEventSubscribeSize)
		}
	}
	{
		temp := &SockAddrStorage{}
		if unsafe.Sizeof(*temp) != SockAddrStorageSize {
			fmt.Println(unsafe.Sizeof(*temp))
			fmt.Println(SockAddrStorageSize)
			t.Error("SockAddrStorage sizes don't match")
		} else {
			fmt.Println("SockAddrStorageSize: ", SockAddrStorageSize)
		}
	}
	{
		temp := &SockAddrIn{}
		if unsafe.Sizeof(*temp) != SockAddrInSize {
			fmt.Println(unsafe.Sizeof(*temp))
			fmt.Println(SockAddrInSize)
			t.Error("SockAddrIn sizes don't match")
		} else {
			fmt.Println("SockAddrInSize: ", SockAddrInSize)
		}
	}
	{
		temp := &SockAddrIn6{}
		if unsafe.Sizeof(*temp) != SockAddrIn6Size {
			fmt.Println(unsafe.Sizeof(*temp))
			fmt.Println(SockAddrIn6Size)
			t.Error("SockAddrIn6 sizes don't match")
		} else {
			fmt.Println("SockAddrIn6Size: ", SockAddrIn6Size)
		}
	}
	{
		temp := &InAddr{}
		if unsafe.Sizeof(*temp) != InAddrSize {
			fmt.Println(unsafe.Sizeof(*temp))
			fmt.Println(InAddrSize)
			t.Error("InAddr sizes don't match")
		} else {
			fmt.Println("InAddrSize: ", InAddrSize)
		}
	}
	{
		temp := &In6Addr{}
		if unsafe.Sizeof(*temp) != In6AddrSize {
			fmt.Println(unsafe.Sizeof(*temp))
			fmt.Println(In6AddrSize)
			t.Error("In6AddrSize sizes don't match")
		} else {
			fmt.Println("In6AddrSize: ", In6AddrSize)
		}
	}
	{
		temp := &SCTPSetPeerPrimary{}
		if unsafe.Sizeof(*temp) != SCTPSetPeerPrimarySize {
			fmt.Println(unsafe.Sizeof(*temp))
			fmt.Println(SCTPSetPeerPrimarySize)
			t.Error("SCTPSetPeerPrimary sizes don't match")
		} else {
			fmt.Println("SCTPSetPeerPrimarySize: ", SCTPSetPeerPrimarySize)
		}
	}
	{
		temp := &SCTPPrimaryAddr{}
		if unsafe.Sizeof(*temp) != SCTPPrimaryAddrSize {
			fmt.Println(unsafe.Sizeof(*temp))
			fmt.Println(SCTPPrimaryAddrSize)
			t.Error("SCTPPrimaryAddr sizes don't match")
		} else {
			fmt.Println("SCTPPrimaryAddrSize: ", SCTPPrimaryAddrSize)
		}
	}
}

func TestPacking(t *testing.T) {
	var buffer []byte
	{
		hdr := &syscall.Cmsghdr{
			Level: syscall.IPPROTO_SCTP,
			Type:  SCTP_SNDRCV,
			Len:   uint64(syscall.CmsgSpace(SCTPSndRcvInfoSize)),
		}
		buffer = append(buffer, Pack(unsafe.Pointer(hdr))...)
	}
	fmt.Printf("%s", hex.Dump(buffer))
	fmt.Println(len(buffer))
	fmt.Println(CMsgHeaderSize)
}
