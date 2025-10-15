package sctp_go

import (
	"errors"
	"sync"
	"sync/atomic"
	"syscall"
)

const (
	MaxNumberOfEvents = 64
)

// Callback represents a function that handles epoll events for a file descriptor.
type Callback func()

// Poller manages epoll-based event polling for file descriptors with associated callbacks.
type Poller struct {
	descriptor atomic.Int32
	callbacks  sync.Map
}

var (
	poller *Poller = nil
)

func init() {
	poller = &Poller{}
	poller.descriptor.Store(-1)
}

// GetPoller returns the global singleton poller instance.
func GetPoller() *Poller {
	return poller
}

// Init initializes the epoll instance by creating an epoll file descriptor.
// Returns an error if already initialized.
func (p *Poller) Init() error {
	if p.IsInitialized() {
		return errors.New("poller already initialized")
	}
	descriptor, err := syscall.EpollCreate1(0)
	if err != nil {
		return err
	}
	p.descriptor.Store(int32(descriptor))
	return nil
}

// Finalize closes the epoll file descriptor if initialized.
func (p *Poller) Finalize() {
	descriptor := p.descriptor.Load()
	if descriptor != -1 {
		syscall.Close(int(descriptor))
		p.descriptor.Store(-1)
	}
}

// IsInitialized checks if the poller has been initialized with a valid epoll descriptor.
func (p *Poller) IsInitialized() bool {
	return p.descriptor.Load() != -1
}

// Add registers a file descriptor with the poller and associates it with a callback function.
// Returns an error if the file descriptor is already registered.
func (p *Poller) Add(fd int, cb Callback) error {
	// Check if already registered
	if _, exists := p.callbacks.Load(fd); exists {
		return errors.New("file descriptor already registered")
	}

	descriptor := p.descriptor.Load()
	err := syscall.EpollCtl(int(descriptor), syscall.EPOLL_CTL_ADD, fd,
		&syscall.EpollEvent{
			Fd:     int32(fd),
			Events: syscall.EPOLLIN,
		},
	)
	if err != nil {
		return err
	}
	p.callbacks.Store(fd, cb)
	return nil
}

// Del removes a file descriptor from the poller and deletes its associated callback.
// Returns an error if the file descriptor is not registered.
func (p *Poller) Del(fd int) error {
	// Check if registered
	if _, exists := p.callbacks.Load(fd); !exists {
		return errors.New("file descriptor not registered")
	}

	p.callbacks.Delete(fd)
	descriptor := p.descriptor.Load()
	return syscall.EpollCtl(int(descriptor), syscall.EPOLL_CTL_DEL, fd, nil)
}

// Loop runs the event loop, waiting for epoll events and executing associated callbacks.
// The loop continues until an error occurs or the poller is not initialized.
func (p *Poller) Loop() {
	if !p.IsInitialized() {
		return
	}

	events := make([]syscall.EpollEvent, MaxNumberOfEvents)
	descriptor := p.descriptor.Load()

	for {
		n, err := syscall.EpollWait(int(descriptor), events, 100)
		if err != nil {
			if err == syscall.EINTR {
				continue
			}
			// break on any other error
			break
		}
		for i := 0; i < n; i++ {
			fd := events[i].Fd
			if callback, ok := p.callbacks.Load(int(fd)); ok {
				if handle, ok := callback.(Callback); ok {
					handle()
				}
			}
		}
	}
}
