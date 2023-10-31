package sctp_go

import (
	"fmt"
	"syscall"
	"testing"
)

//go:generate go test -count=1 -race -timeout 30s -run ^TestPtr$ github.com/thebagchi/sctp-go
//go:generate go test -count=1 -gcflags=all=-d=checkptr -timeout 30s -run ^TestPtr$ github.com/thebagchi/sctp-go

func TestPtr(t *testing.T) {

	addr, err := MakeSCTPAddr("sctp", "::1/127.0.0.1:12345")
	if nil != err {
		fmt.Println("Error: ", err)
		t.FailNow()
	}

	server, err := ListenSCTP(
		"sctp",
		syscall.SOCK_STREAM,
		addr,
		&SCTPInitMsg{
			NumOutStreams:  0xffff,
			MaxInStreams:   0,
			MaxAttempts:    0,
			MaxInitTimeout: 0,
		},
	)

	if nil != err {
		fmt.Println("Error: ", err)
		t.FailNow()
	}

	defer server.Close()

	if local := server.Addr(); nil != local {
		fmt.Println("Addr: ", local)
	} else {
		fmt.Println("Error: ", err)
		t.FailNow()
	}

}
