package main

import (
	"fmt"
	"syscall"
	"time"

	sctp "github.com/thebagchi/sctp-go"
)

func main() {

	if err := sctp.GetPoller().Init(); nil != err {
		fmt.Println(err.Error())
		return
	}

	defer sctp.GetPoller().Finalize()

	address := "127.0.0.1:12345"
	addr, err := sctp.MakeSCTPAddr("sctp", address)
	if nil != err {
		fmt.Println(err.Error())
		return
	}
	client, err := sctp.ListenSCTP(
		"sctp",
		syscall.SOCK_SEQPACKET,
		addr,
		&sctp.SCTPInitMsg{
			NumOutStreams:  10,
			MaxInStreams:   10,
			MaxAttempts:    10,
			MaxInitTimeout: 0,
		},
	)
	if nil != err {
		fmt.Println(err.Error())
		return
	}
	events, err := client.GetEventSubscribe()
	if nil != err {
		fmt.Println("Error: ", err)
	}
	if nil != events {
		events.DataIoEvent = 1
		events.AssociationEvent = 1
		events.ShutdownEvent = 1
		err = client.SetEventSubscribe(events)
		if nil != err {
			fmt.Println("Error: ", err)
		}
	}
	go func() {
		fmt.Println("Registering fd: ", client.FD())
		sctp.GetPoller().Add(client.FD(), func() {
			info := &sctp.SCTPSndRcvInfo{}
			flag := 0
			data := make([]byte, 8192)
			_, err := client.RecvMsg(data, info, &flag)
			if err != nil {
				fmt.Println("Error: ", err)
			}
		})
		sctp.GetPoller().Loop()
	}()
	fmt.Println("sctp client started: ", time.Now())
	time.Sleep(30 * time.Second)
	fmt.Println("sctp client closing: ", time.Now())
	sctp.GetPoller().Del(client.FD())
	err = client.Close()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("sctp client closed")
	time.Sleep(60 * time.Second)
}
