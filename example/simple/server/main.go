package main

import (
	"encoding/hex"
	"fmt"
	"os"
	"syscall"

	sctp "github.com/thebagchi/sctp-go"
)

func HandleClient(conn *sctp.SCTPConn) {
	if nil == conn {
		return
	}
	var (
		data = make([]byte, 8192)
		flag = 0
	)
	for {
		info := &sctp.SCTPSndRcvInfo{}
		len, err := conn.RecvMsg(data, info, &flag)
		if nil != err {
			fmt.Println("Error: ", err)
			break
		}
		if len == 0 {
			fmt.Println("Connection terminated!!!")
			break
		} else {
			fmt.Println("Rcvd bytes: ", len)
			buffer := data[:len]
			fmt.Println(hex.Dump(buffer))
			fmt.Println(hex.Dump(sctp.Pack(info)))
		}
	}
}

func main() {

	addr, err := sctp.MakeSCTPAddr("sctp4", "127.0.0.1:12345")
	if nil != err {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}

	server, err := sctp.ListenSCTP(
		"sctp4",
		syscall.SOCK_STREAM,
		addr,
		&sctp.SCTPInitMsg{
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

	if local := server.Addr(); nil != local {
		fmt.Println("Addr: ", local)
	} else {
		fmt.Println("Error: local addr not received")
	}

	for {
		conn, err := server.AcceptSCTP()
		if nil != err {
			fmt.Println("Error: ", err)
			continue
		}
		if remote := conn.RemoteAddr(); nil != remote {
			fmt.Println("New connection from: ", remote)
		}
		go HandleClient(conn)
	}
}
