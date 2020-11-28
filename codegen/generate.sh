#!/bin/sh
echo "Generating sctp constants from /usr/include/linux/sctp.h"
go tool cgo -godefs -srcdir . constants_sctp.go > ../sctp_constants.go
gofmt -w ../sctp_constants.go

echo "Generating sctp types from /usr/include/linux/sctp.h"
go tool cgo -godefs -srcdir . types_sctp.go > ../sctp_types.go
gofmt -w ../sctp_types.go