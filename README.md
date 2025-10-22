# sctp-go

[![Go Version](https://img.shields.io/badge/go-%3E%3D1.17-blue.svg)](https://golang.org/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A Go library for implementing the Stream Control Transmission Protocol (SCTP) in Go applications. SCTP is a transport-layer protocol that provides reliable, message-oriented data transfer with features like multi-streaming, multi-homing, and congestion control.

## Features

- Full SCTP protocol implementation
- Support for both IPv4 and IPv6
- Multi-streaming and multi-homing capabilities
- Connection-oriented and message-oriented communication
- Integration with Go's net package interfaces
- Comprehensive test suite

## Installation

```bash
go get github.com/thebagchi/sctp-go
```

## Usage

### Basic Client

```go
package main

import (
    "fmt"
    "os"

    sctp "github.com/thebagchi/sctp-go"
)

func main() {
    local, err := sctp.MakeSCTPAddr("sctp4", "127.0.0.1:54321")
    if err != nil {
        fmt.Println("Error:", err)
        os.Exit(1)
    }

    remote, err := sctp.MakeSCTPAddr("sctp4", "127.0.0.1:12345")
    if err != nil {
        fmt.Println("Error:", err)
        os.Exit(1)
    }

    conn, err := sctp.DialSCTP(
        "sctp4",
        local,
        remote,
        &sctp.SCTPInitMsg{
            NumOutStreams:  0xFFFF,
            MaxInStreams:   0,
            MaxAttempts:    0,
            MaxInitTimeout: 0,
        },
    )
    if err != nil {
        fmt.Println("Error:", err)
        os.Exit(1)
    }
    defer conn.Close()

    // Use the connection...
}
```

### Basic Server

```go
package main

import (
    "fmt"
    "os"
    "syscall"

    sctp "github.com/thebagchi/sctp-go"
)

func main() {
    addr, err := sctp.MakeSCTPAddr("sctp4", "127.0.0.1:12345")
    if err != nil {
        fmt.Println("Error:", err)
        os.Exit(1)
    }

    server, err := sctp.ListenSCTP("sctp4", syscall.SOCK_STREAM, addr)
    if err != nil {
        fmt.Println("Error:", err)
        os.Exit(1)
    }
    defer server.Close()

    for {
        conn, err := server.Accept()
        if err != nil {
            fmt.Println("Error:", err)
            continue
        }

        go handleConnection(conn.(*sctp.SCTPConn))
    }
}

func handleConnection(conn *sctp.SCTPConn) {
    defer conn.Close()
    // Handle the connection...
}
```

## Examples

The `example/` directory contains several example implementations:

- **Simple**: Basic client-server communication (`example/simple/`)
- **Packet**: Sequential packet examples (`example/packet/`)
- **Epoll**: Epoll-based implementation (`example/epoll/`)

To run the simple server:

```bash
go run example/simple/server/main.go
```

To run the simple client:

```bash
go run example/simple/client/main.go
```

## API Overview

### Connection Management
- `DialSCTP()` - Establish an SCTP connection
- `ListenSCTP()` - Create an SCTP listener
- `Accept()` - Accept incoming connections

### Address Handling
- `MakeSCTPAddr()` - Create SCTP addresses
- `ResolveSCTPAddr()` - Resolve hostnames to SCTP addresses

### Data Transfer
- `SendMsg()` - Send messages with stream information
- `RecvMsg()` - Receive messages with stream information

### Connection Information
- `GetInitMsg()` - Get initialization message
- `GetPrimaryPeerAddr()` - Get primary peer address
- `RemoteAddr()` / `LocalAddr()` - Get connection addresses

## Testing

Run the test suite:

```bash
go test
```

Run with race detection:

```bash
go test -race
```

## Requirements

- Go 1.17 or later
- Linux kernel with SCTP support (most modern distributions)

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
