package main

import (
	"fmt"
	sctpgo "github.com/thebagchi/sctp-go"
	"os"
)

func main() {
	local, err := sctpgo.MakeSCTPAddr("sctp4", "127.0.0.1:54321")
	if nil != err {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}
	remote, err := sctpgo.MakeSCTPAddr("sctp4", "127.0.0.1:12345")
	if nil != err {
		fmt.Println("Error: ", err)
		os.Exit(2)
	}
	conn, err := sctpgo.DialSCTP(
		"sctp4",
		local,
		remote,
		&sctpgo.SCTPInitMsg{
			NumOutStreams:  0xFFFF,
			MaxInStreams:   0,
			MaxAttempts:    0,
			MaxInitTimeout: 0,
		},
	)
	if nil != err {
		fmt.Println("Error: ", err)
		os.Exit(3)
	}
	defer conn.Close()

	if peer, err := conn.GetPrimaryPeerAddr(); nil == err {
		fmt.Println("Peer: ", peer)
	} else {
		fmt.Println("Error: ", err)
	}

	if remote := conn.RemoteAddr(); nil != remote {
		fmt.Println("Peer: ", remote)
	} else {
		fmt.Println("Error: remote addr not received")
	}

	if local := conn.LocalAddr(); nil != local {
		fmt.Println("Peer: ", local)
	} else {
		fmt.Println("Error: local addr not received")
	}
}
