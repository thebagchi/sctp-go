package main

import (
	"fmt"
	sctp_go "github.com/thebagchi/sctp-go"
	"os"
)

func main() {
	local, err := sctp_go.MakeSCTPAddr("sctp4", "127.0.0.1:54321")
	if nil != err {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}
	remote, err := sctp_go.MakeSCTPAddr("sctp4", "127.0.0.1:12345")
	if nil != err {
		fmt.Println("Error: ", err)
		os.Exit(2)
	}
	conn, err := sctp_go.DialSCTP(
		"sctp4",
		local,
		remote,
		&sctp_go.SCTPInitMsg{
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
		fmt.Println("Local: ", local)
	} else {
		fmt.Println("Error: local addr not received")
	}

	len, err := conn.SendMsg([]byte("HELLO WORLD"), nil)

	if nil != err {
		fmt.Println("Error: ", err)
	} else {
		fmt.Println(fmt.Sprintf("Sent %d bytes", len))
	}

}
