# Using sctp in golang.

## See example/sctp_server.go for server implementation.
## See example/sctp_client.go for client implementation.

### Running Server

#### Stream Server

```sh
$ go run example/sctp_server.go
```

#### Sequential Packet Server

```sh
$ go run example/packet_server.go
```

### Running Client

#### Stream Client 

```sh
$ go run example/sctp_client.go 
```

#### Sequential Packet Client

```sh
$ go run example/packet_client.go
```

#### Getting

```sh
$ go get github.com/thebagchi/sctp-go
```
