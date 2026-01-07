package per

import (
	"encoding/asn1"
	"fmt"
)

// Empty represents an empty type with no value
// Used in PER encoding for placeholder or marker fields
type Empty struct{}

// Null represents the ASN.1 NULL type (deprecated, use NULL instead)
// NULL is a simple type that has no value (encoding produces no bits)
// ITU-T X.691 Section 10.0: Null
type Null struct{}

// NewNull creates a new Null instance
func NewNull() *Null {
	return &Null{}
}

// NewEmpty creates a new Empty instance
func NewEmpty() *Empty {
	return &Empty{}
}

// BitStringToBase2 converts an asn1.BitString to a base-2 string representation
// The resulting string contains only '0' and '1' characters representing the bits
// Example: BitString with 4 bits "1010" returns "1010"
func BitStringToBase2(bs *asn1.BitString) string {
	if bs == nil {
		return ""
	}

	if bs.BitLength == 0 {
		return ""
	}

	result := ""

	// Iterate through each bit in the BitString
	for i := 0; i < bs.BitLength; i++ {
		byteIdx := i / 8
		bitIdx := 7 - (i % 8) // MSB first (most significant bit in leftmost position)

		if byteIdx >= len(bs.Bytes) {
			break
		}

		// Extract the bit at position bitIdx
		bit := (bs.Bytes[byteIdx] >> uint(bitIdx)) & 1
		if bit == 1 {
			result += "1"
		} else {
			result += "0"
		}
	}

	return result
}

// BitStringFromBase2 creates an asn1.BitString from a base-2 string
// The input string should contain only '0' and '1' characters
// Example: "1010" creates a BitString with 4 bits
// Supports arbitrary length bit strings (not limited to 64 bits)
func BitStringFromBase2(binaryStr string) (*asn1.BitString, error) {
	if binaryStr == "" {
		return &asn1.BitString{}, nil
	}

	// Validate that the string contains only '0' and '1'
	for _, ch := range binaryStr {
		if ch != '0' && ch != '1' {
			return nil, fmt.Errorf("invalid binary string: contains non-binary character '%c'", ch)
		}
	}

	bitLen := len(binaryStr)
	numBytes := (bitLen + 7) / 8
	bytes := make([]byte, numBytes)

	// Convert the binary string to bytes directly (works for any length)
	for i, ch := range binaryStr {
		byteIdx := i / 8
		bitIdx := 7 - (i % 8) // MSB first

		if ch == '1' {
			bytes[byteIdx] |= (1 << uint(bitIdx))
		}
	}

	return &asn1.BitString{
		Bytes:     bytes,
		BitLength: bitLen,
	}, nil
}
