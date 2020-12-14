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
	{
		temp := &SCTPPeelOffArg{}
		if unsafe.Sizeof(*temp) != SCTPPeelOffArgSize {
			fmt.Println(unsafe.Sizeof(*temp))
			fmt.Println(SCTPPeelOffArgSize)
			t.Error("SCTPPeelOffArg sizes don't match")
		} else {
			fmt.Println("SCTPPeelOffArgSize: ", SCTPPeelOffArgSize)
		}
	}
	{
		temp := &SCTPPeelOffFlagsArg{}
		if unsafe.Sizeof(*temp) != SCTPPeelOffFlagsArgSize {
			fmt.Println(unsafe.Sizeof(*temp))
			fmt.Println(SCTPPeelOffFlagsArgSize)
			t.Error("SCTPPeelOffFlagsArg sizes don't match")
		} else {
			fmt.Println("SCTPPeelOffFlagsArgSize: ", SCTPPeelOffFlagsArgSize)
		}
	}
	{
		temp := &SCTPAssocChange{}
		if unsafe.Sizeof(*temp) != SCTPAssocChangeSize {
			fmt.Println(unsafe.Sizeof(*temp))
			fmt.Println(SCTPAssocChangeSize)
			t.Error("SCTPAssocChange sizes don't match")
		} else {
			fmt.Println("SCTPAssocChangeSize: ", SCTPAssocChangeSize)
		}
	}
	{
		temp := &SCTPPAddrChange{}
		if unsafe.Sizeof(*temp) != SCTPPAddrChangeSize {
			fmt.Println(unsafe.Sizeof(*temp))
			fmt.Println(SCTPPAddrChangeSize)
			t.Error("SCTPPAddrChange sizes don't match")
		} else {
			fmt.Println("SCTPPAddrChangeSize: ", SCTPPAddrChangeSize)
		}
	}
	{
		temp := &SCTPRemoteError{}
		if unsafe.Sizeof(*temp) != SCTPRemoteErrorSize {
			fmt.Println(unsafe.Sizeof(*temp))
			fmt.Println(SCTPRemoteErrorSize)
			t.Error("SCTPRemoteError sizes don't match")
		} else {
			fmt.Println("SCTPRemoteErrorSize: ", SCTPRemoteErrorSize)
		}
	}
	{
		temp := &SCTPSendFailed{}
		if unsafe.Sizeof(*temp) != SCTPSendFailedSize {
			fmt.Println(unsafe.Sizeof(*temp))
			fmt.Println(SCTPSendFailedSize)
			t.Error("SCTPSendFailed sizes don't match")
		} else {
			fmt.Println("SCTPSendFailedSize: ", SCTPSendFailedSize)
		}
	}
	{
		temp := &SCTPShutdownEvent{}
		if unsafe.Sizeof(*temp) != SCTPShutdownEventSize {
			fmt.Println(unsafe.Sizeof(*temp))
			fmt.Println(SCTPShutdownEventSize)
			t.Error("SCTPShutdownEvent sizes don't match")
		} else {
			fmt.Println("SCTPShutdownEventSize: ", SCTPShutdownEventSize)
		}
	}
	{
		temp := &SCTPAdaptationEvent{}
		if unsafe.Sizeof(*temp) != SCTPAdaptationEventSize {
			fmt.Println(unsafe.Sizeof(*temp))
			fmt.Println(SCTPAdaptationEventSize)
			t.Error("SCTPAdaptationEvent sizes don't match")
		} else {
			fmt.Println("SCTPAdaptationEventSize: ", SCTPAdaptationEventSize)
		}
	}
	{
		temp := &SCTPPDApiEvent{}
		if unsafe.Sizeof(*temp) != SCTPPDApiEventSize {
			fmt.Println(unsafe.Sizeof(*temp))
			fmt.Println(SCTPPDApiEventSize)
			t.Error("SCTPPDApiEvent sizes don't match")
		} else {
			fmt.Println("SCTPPDApiEventSize: ", SCTPPDApiEventSize)
		}
	}
	{
		temp := &SCTPAuthKeyEvent{}
		if unsafe.Sizeof(*temp) != SCTPAuthKeyEventSize {
			fmt.Println(unsafe.Sizeof(*temp))
			fmt.Println(SCTPAuthKeyEventSize)
			t.Error("SCTPAuthKeyEvent sizes don't match")
		} else {
			fmt.Println("SCTPAuthKeyEventSize: ", SCTPAuthKeyEventSize)
		}
	}
	{
		temp := &SCTPSenderDryEvent{}
		if unsafe.Sizeof(*temp) != SCTPSenderDryEventSize {
			fmt.Println(unsafe.Sizeof(*temp))
			fmt.Println(SCTPSenderDryEventSize)
			t.Error("SCTPSenderDryEvent sizes don't match")
		} else {
			fmt.Println("SCTPSenderDryEventSize: ", SCTPSenderDryEventSize)
		}
	}
	{
		temp := &SCTPStreamResetEvent{}
		if unsafe.Sizeof(*temp) != SCTPStreamResetEventSize {
			fmt.Println(unsafe.Sizeof(*temp))
			fmt.Println(SCTPStreamResetEventSize)
			t.Error("SCTPStreamResetEvent sizes don't match")
		} else {
			fmt.Println("SCTPStreamResetEventSize: ", SCTPStreamResetEventSize)
		}
	}
	{
		temp := &SCTPAssocResetEvent{}
		if unsafe.Sizeof(*temp) != SCTPAssocResetEventSize {
			fmt.Println(unsafe.Sizeof(*temp))
			fmt.Println(SCTPAssocResetEventSize)
			t.Error("SCTPAssocResetEvent sizes don't match")
		} else {
			fmt.Println("SCTPAssocResetEventSize: ", SCTPAssocResetEventSize)
		}
	}
	{
		temp := &SCTPAssocResetEvent{}
		if unsafe.Sizeof(*temp) != SCTPAssocResetEventSize {
			fmt.Println(unsafe.Sizeof(*temp))
			fmt.Println(SCTPAssocResetEventSize)
			t.Error("SCTPAssocResetEvent sizes don't match")
		} else {
			fmt.Println("SCTPAssocResetEventSize: ", SCTPAssocResetEventSize)
		}
	}
	{
		temp := &SCTPStreamChangeEvent{}
		if unsafe.Sizeof(*temp) != SCTPStreamChangeEventSize {
			fmt.Println(unsafe.Sizeof(*temp))
			fmt.Println(SCTPStreamChangeEventSize)
			t.Error("SCTPStreamChangeEvent sizes don't match")
		} else {
			fmt.Println("SCTPStreamChangeEventSize: ", SCTPStreamChangeEventSize)
		}
	}
	{
		temp := &SCTPRTOInfo{}
		if unsafe.Sizeof(*temp) != SCTPRTOInfoSize {
			fmt.Println(unsafe.Sizeof(*temp))
			fmt.Println(SCTPRTOInfoSize)
			t.Error("SCTPRTOInfo sizes don't match")
		} else {
			fmt.Println("SCTPRTOInfoSize: ", SCTPRTOInfoSize)
		}
	}
	{
		temp := &SCTPResetStreams{}
		if unsafe.Sizeof(*temp) != SCTPResetStreamsSize {
			fmt.Println(unsafe.Sizeof(*temp))
			fmt.Println(SCTPResetStreamsSize)
			t.Error("SCTPResetStreams sizes don't match")
		} else {
			fmt.Println("SCTPResetStreamsSize: ", SCTPResetStreamsSize)
		}
	}
	{
		temp := &SCTPAddStreams{}
		if unsafe.Sizeof(*temp) != SCTPAddStreamsSize {
			fmt.Println(unsafe.Sizeof(*temp))
			fmt.Println(SCTPAddStreamsSize)
			t.Error("SCTPAddStreams sizes don't match")
		} else {
			fmt.Println("SCTPAddStreamsSize: ", SCTPAddStreamsSize)
		}
	}
	{
		temp := &SCTPAssocParams{}
		if unsafe.Sizeof(*temp) != SCTPAssocParamsSize {
			fmt.Println(unsafe.Sizeof(*temp))
			fmt.Println(SCTPAssocParamsSize)
			t.Error("SCTPAssocParams sizes don't match")
		} else {
			fmt.Println("SCTPAssocParamsSize: ", SCTPAssocParamsSize)
		}
	}
	{
		temp := &SCTPSetAdaptation{}
		if unsafe.Sizeof(*temp) != SCTPSetAdaptationSize {
			fmt.Println(unsafe.Sizeof(*temp))
			fmt.Println(SCTPSetAdaptationSize)
			t.Error("SCTPSetAdaptation sizes don't match")
		} else {
			fmt.Println("SCTPSetAdaptationSize: ", SCTPSetAdaptationSize)
		}
	}
	{
		temp := &SCTPPeerAddrParams{}
		if unsafe.Sizeof(*temp) != SCTPPeerAddrParamsSize {
			fmt.Println(unsafe.Sizeof(*temp))
			fmt.Println(SCTPPeerAddrParamsSize)
			{
				if len(Pack(temp)) != SCTPPeerAddrParamsSize {
					t.Error("SCTPPeerAddrParams sizes don't match")
				} else {
					fmt.Println("SCTPPeerAddrParams sizes don't match")
				}
			}
		} else {
			fmt.Println("SCTPPeerAddrParamsSize: ", SCTPPeerAddrParamsSize)
		}
	}
	{
		temp := &SCTPPeerAddrInfo{}
		if unsafe.Sizeof(*temp) != SCTPPeerAddrInfoSize {
			fmt.Println(unsafe.Sizeof(*temp))
			fmt.Println(SCTPPeerAddrInfoSize)
			t.Error("SCTPPeerAddrInfo sizes don't match")
		} else {
			fmt.Println("SCTPPeerAddrInfoSize: ", SCTPPeerAddrInfoSize)
		}
	}
	{
		temp := &SCTPAssocValue{}
		if unsafe.Sizeof(*temp) != SCTPAssocValueSize {
			fmt.Println(unsafe.Sizeof(*temp))
			fmt.Println(SCTPAssocValueSize)
			t.Error("SCTPAssocValue sizes don't match")
		} else {
			fmt.Println("SCTPAssocValueSize: ", SCTPAssocValueSize)
		}
	}
	{
		temp := &SCTPSackInfo{}
		if unsafe.Sizeof(*temp) != SCTPSackInfoSize {
			fmt.Println(unsafe.Sizeof(*temp))
			fmt.Println(SCTPSackInfoSize)
			t.Error("SCTPSackInfo sizes don't match")
		} else {
			fmt.Println("SCTPSackInfoSize: ", SCTPSackInfoSize)
		}
	}
	{
		temp := &SCTPStreamValue{}
		if unsafe.Sizeof(*temp) != SCTPStreamValueSize {
			fmt.Println(unsafe.Sizeof(*temp))
			fmt.Println(SCTPStreamValueSize)
			t.Error("SCTPStreamValue sizes don't match")
		} else {
			fmt.Println("SCTPStreamValueSize: ", SCTPStreamValueSize)
		}
	}
	{
		temp := &SCTPStatus{}
		if unsafe.Sizeof(*temp) != SCTPStatusSize {
			fmt.Println(unsafe.Sizeof(*temp))
			fmt.Println(SCTPStatusSize)
			t.Error("SCTPStatus sizes don't match")
		} else {
			fmt.Println("SCTPStatusSize: ", SCTPStatusSize)
		}
	}
	{
		temp := &SCTPAuthKeyId{}
		if unsafe.Sizeof(*temp) != SCTPAuthKeyIdSize {
			fmt.Println(unsafe.Sizeof(*temp))
			fmt.Println(SCTPAuthKeyIdSize)
			t.Error("SCTPAuthKeyId sizes don't match")
		} else {
			fmt.Println("SCTPAuthKeyIdSize: ", SCTPAuthKeyIdSize)
		}
	}
	{
		temp := &SCTPAssocStats{}
		if unsafe.Sizeof(*temp) != SCTPAssocStatsSize {
			fmt.Println(unsafe.Sizeof(*temp))
			fmt.Println(SCTPAssocStatsSize)
			t.Error("SCTPAssocStats sizes don't match")
		} else {
			fmt.Println("SCTPAssocStatsSize: ", SCTPAssocStatsSize)
		}
	}
	{
		temp := &SCTPPeerAddrThresholds{}
		if unsafe.Sizeof(*temp) != SCTPPeerAddrThresholdsSize {
			fmt.Println(unsafe.Sizeof(*temp))
			fmt.Println(SCTPPeerAddrThresholdsSize)
			t.Error("SCTPPeerAddrThresholds sizes don't match")
		} else {
			fmt.Println("SCTPPeerAddrThresholdsSize: ", SCTPPeerAddrThresholdsSize)
		}
	}
	{
		temp := &SCTPPRStatus{}
		if unsafe.Sizeof(*temp) != SCTPPRStatusSize {
			fmt.Println(unsafe.Sizeof(*temp))
			fmt.Println(SCTPPRStatusSize)
			t.Error("SCTPPRStatus sizes don't match")
		} else {
			fmt.Println("SCTPPRStatusSize: ", SCTPPRStatusSize)
		}
	}
	{
		temp := &SCTPDefaultPRInfo{}
		if unsafe.Sizeof(*temp) != SCTPDefaultPRInfoSize {
			fmt.Println(unsafe.Sizeof(*temp))
			fmt.Println(SCTPDefaultPRInfoSize)
			t.Error("SCTPDefaultPRInfo sizes don't match")
		} else {
			fmt.Println("SCTPDefaultPRInfoSize: ", SCTPDefaultPRInfoSize)
		}
	}
	{
		temp := &SCTPEvent{}
		if unsafe.Sizeof(*temp) != SCTPEventSize {
			fmt.Println(unsafe.Sizeof(*temp))
			fmt.Println(SCTPEventSize)
			t.Error("SCTPEvent sizes don't match")
		} else {
			fmt.Println("SCTPEventSize: ", SCTPEventSize)
		}
	}
	{
		temp := &SCTPInfo{}
		if unsafe.Sizeof(*temp) != SCTPInfoSize {
			fmt.Println(unsafe.Sizeof(*temp))
			fmt.Println(SCTPInfoSize)
			t.Error("SCTPInfo sizes don't match")
		} else {
			fmt.Println("SCTPInfoSize: ", SCTPInfoSize)
		}
	}
	{
		fmt.Println("CMSG_SPACE(sizeof(SCTPSndRcvInfo)): ", syscall.CmsgSpace(SCTPSndRcvInfoSize))
	}
}

func TestPacking(t *testing.T) {
	var buffer []byte
	{
		hdr := &syscall.Cmsghdr{
			Level: syscall.IPPROTO_SCTP,
			Type:  SCTP_SNDRCV,
		}
		hdr.SetLen(syscall.CmsgSpace(SCTPSndRcvInfoSize))
		info := &SCTPSndRcvInfo{
			AssocId: 100,
		}
		buffer = append(buffer, Pack(hdr)...)
		buffer = append(buffer, Pack(info)...)
		fmt.Println(len(Pack(hdr)))
		fmt.Println(len(Pack(info)))
	}
	fmt.Printf("%s", hex.Dump(buffer))
	fmt.Println(len(buffer))
	fmt.Println(CMsgHeaderSize)
}

func TestMakeSockaddr(t *testing.T) {
	addr, err := MakeSCTPAddr("sctp4", "127.0.0.1:12345")
	if nil != err {
		fmt.Println("Error: ", err)
		t.FailNow()
	}
	{
		buffer := MakeSockaddr(addr)
		fmt.Printf(hex.Dump(buffer))
		fmt.Println(len(buffer))
	}
}
