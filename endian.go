package sctp_go

import (
	"encoding/binary"
)

func HostToNetworkShort(number uint16) uint16 {
	if endian == binary.LittleEndian {
		return (number << 8 & 0xff00) | (number >> 8 & 0x00ff)
	}
	return number
}

func NetworkToHostShort(number uint16) uint16 {
	if endian == binary.LittleEndian {
		return (number << 8 & 0xff00) | (number >> 8 & 0x00ff)
	}
	return number
}

func HostToNetwork(number uint32) uint32 {
	if endian == binary.LittleEndian {
		return uint32(HostToNetworkShort(uint16(number)))<<16 | uint32(HostToNetworkShort(uint16(number>>16)))
	}
	return number
}

func NetworkToHost(number uint32) uint32 {
	if endian == binary.LittleEndian {
		return uint32(NetworkToHostShort(uint16(number)))<<16 | uint32(NetworkToHostShort(uint16(number>>16)))
	}
	return number
}

func HostToNetworkLong(number uint64) uint64 {
	if endian == binary.LittleEndian {
		return uint64(HostToNetwork(uint32(number)))<<32 | uint64(HostToNetwork(uint32(number>>32)))
	}
	return number
}

func NetworkToHostLong(number uint64) uint64 {
	if endian == binary.LittleEndian {
		return uint64(NetworkToHost(uint32(number)))<<32 | uint64(NetworkToHost(uint32(number>>32)))
	}
	return number
}
