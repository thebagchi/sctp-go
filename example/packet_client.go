package main

import (
	"fmt"
	sctp "github.com/thebagchi/sctp-go"
	"os"
	"syscall"
)

func main() {
	local, err := sctp.MakeSCTPAddr("sctp4", "127.0.0.1:54321")
	if nil != err {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}
	remote, err := sctp.MakeSCTPAddr("sctp4", "127.0.0.1:12345")
	if nil != err {
		fmt.Println("Error: ", err)
		os.Exit(2)
	}

	client, err := sctp.ListenSCTP(
		"sctp4",
		syscall.SOCK_SEQPACKET,
		local,
		&sctp.SCTPInitMsg{
			NumOutStreams:  0xffff,
			MaxInStreams:   0,
			MaxAttempts:    0,
			MaxInitTimeout: 0,
		},
	)
	if nil != err {
		fmt.Println("Error: ", err)
		os.Exit(3)
	}

	defer client.Close()
	assoc, err := client.Connect(remote)
	if nil != err {
		fmt.Println("Error: ", err)
		os.Exit(4)
	}

	length, err := client.SendMsg(
		[]byte("HELLO WORLD"),
		&sctp.SCTPSndRcvInfo{
			AssocId: int32(assoc),
			Ppid:    0,
			Stream:  0,
			Flags:   0,
		},
	)
	if nil != err {
		fmt.Println("Error: ", err)
	} else {
		fmt.Println(fmt.Sprintf("Sent %d bytes", length))
	}

	err = client.Disconnect(assoc)
	if nil != err {
		fmt.Println("Error: ", err)
	}
}
