package sctp_go

import (
	"fmt"
	"sync"
	"syscall"
)

const (
	MaxNumberOfEvents = 64
)

type Callback func()

type Poller struct {
	descriptor int
	callbacks  sync.Map
}

var (
	poller *Poller = nil
)

func init() {
	poller = &Poller{
		descriptor: -1,
	}
}

func GetPoller() *Poller {
	return poller
}

func (p *Poller) Init() error {
	var err error = nil
	descriptor, err := syscall.EpollCreate1(0)
	if nil == err {
		p.descriptor = descriptor
		return nil
	}
	return err
}

func (p *Poller) Finalize() {
	if p.IsInitialized() {
		syscall.Close(p.descriptor)
	}
}

func (p *Poller) IsInitialized() bool {
	return p.descriptor != -1
}

func (p *Poller) Add(fd int, cb Callback) error {
	err := syscall.EpollCtl(p.descriptor, syscall.EPOLL_CTL_ADD, fd,
		&syscall.EpollEvent{
			Fd:     int32(fd),
			Events: syscall.EPOLLIN,
		},
	)
	if nil == err {
		p.callbacks.Store(fd, cb)
	}
	return err
}

func (p *Poller) Del(fd int) error {
	p.callbacks.Delete(fd)
	err := syscall.EpollCtl(p.descriptor, syscall.EPOLL_CTL_DEL, fd, nil)
	return err
}

func (p *Poller) Loop() {
	events := make([]syscall.EpollEvent, MaxNumberOfEvents)
	for {
		n, err := syscall.EpollWait(p.descriptor, events, 100)
		if err != nil && err != syscall.EINTR {
			break
		}
		for i := 0; i < n; i++ {
			fd := events[i].Fd
			if callback, ok := p.callbacks.Load(int(fd)); ok {
				if handle, ok := callback.(Callback); ok {
					fmt.Println("Calling for fd: ", fd)
					handle()
				}
			}
		}
	}
}
