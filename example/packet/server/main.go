package main

import (
	"encoding/hex"
	"fmt"
	"os"
	"syscall"

	sctp "github.com/thebagchi/sctp-go"
)

func main() {
	addr, err := sctp.MakeSCTPAddr("sctp4", "127.0.0.1:12345")
	if nil != err {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}

	server, err := sctp.ListenSCTP(
		"sctp4",
		syscall.SOCK_SEQPACKET,
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

	events, err := server.GetEventSubscribe()
	if nil != err {
		fmt.Println("Error: ", err)
	}
	if nil != events {
		events.DataIoEvent = 1
		events.AssociationEvent = 1
		events.ShutdownEvent = 1

		err = server.SetEventSubscribe(events)
		if nil != err {
			fmt.Println("Error: ", err)
		}
	}
	{
		var (
			data = make([]byte, 8192)
			flag = 0
		)
		for {
			info := &sctp.SCTPSndRcvInfo{}
			len, err := server.RecvMsg(data, info, &flag)
			if nil != err {
				fmt.Println("Error: ", err)
				break
			}
			if len == 0 {
				fmt.Println("Connection terminated!!!")
				break
			} else {
				if flag&sctp.SCTP_MSG_NOTIFICATION > 0 {
					notification, err := sctp.ParseNotification(data[:len])
					if nil != err {
						fmt.Println("Error: ", err)
					} else {
						fmt.Println(
							"Notification received: ",
							sctp.NotificationName(notification.GetType()),
							notification.GetType(),
						)
						if notification.GetType() == sctp.SCTP_ASSOC_CHANGE {
							if event, ok := notification.(*sctp.SCTPAssocChange); ok {
								fmt.Println(event.State, event.Flags, event.AssocId)
							}
						}
					}
				} else {
					fmt.Println("Rcvd bytes: ", len)
					buffer := data[:len]
					fmt.Println(hex.Dump(buffer))
					fmt.Println(hex.Dump(sctp.Pack(info)))
				}
			}
		}
	}
}
