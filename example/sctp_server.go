package main

import (
	"fmt"
	sctpgo "github.com/thebagchi/sctp-go"
	"os"
)

func main() {
	if true {
		addr, err := sctpgo.MakeSCTPAddr("sctp4", "127.0.0.1:12345")
		if nil != err {
			fmt.Println("Error: ", err)
			os.Exit(1)
		}

		server, err := sctpgo.ListenSCTP(
			"sctp4",
			addr,
			&sctpgo.SCTPInitMsg{
				NumOutStreams:  0xffff,
				MaxInStreams:   0,
				MaxAttempts:    0,
				MaxInitTimeout: 0,
			},
		)
		if nil != err {
			fmt.Println("Error: ", err)
			os.Exit(2)
		}

		defer server.Close()

		for {
			conn, err := server.AcceptSCTP()
			if nil != err {
				fmt.Println("Error: ", err)
				continue
			}
			_ = conn
		}
	}
}
