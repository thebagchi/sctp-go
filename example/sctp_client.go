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
	conn, err := sctpgo.DialSCTP("sctp4", local, remote, &sctpgo.SCTPInitMsg{
		NumOutStreams:  0xFFFF,
		MaxInStreams:   0,
		MaxAttempts:    0,
		MaxInitTimeout: 0,
	})
	if nil != err {
		fmt.Println("Error: ", err)
		os.Exit(3)
	}
	_ = conn
}
