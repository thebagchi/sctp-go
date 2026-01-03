package per

import (
	"encoding/asn1"
	"fmt"
	"math"
	"math/bits"
	"unsafe"

	"github.com/thebagchi/sctp-go/gnbsim/per/bitbuffer"
)

// Encoder represents a PER encoder for bit-level encoding
type Encoder struct {
	codec   *bitbuffer.Codec
	aligned bool
}

// NewEncoder creates a new PER encoder
// aligned: true for APER (Aligned PER), false for UPER (Unaligned PER)
func NewEncoder(aligned bool) *Encoder {
	return &Encoder{
		codec:   bitbuffer.CreateWriter(),
		aligned: aligned,
	}
}

// Bytes returns the encoded bytes, correctly trimmed to the exact bit length
func (e *Encoder) Bytes() []byte {
	return e.codec.Bytes()
}

// Align aligns to the next byte boundary in aligned mode only
func (e *Encoder) Align() error {
	if e.aligned {
		return e.codec.Align()
	}
	return nil
}

// WriteInt encodes an INTEGER type according to X.691
func (e *Encoder) WriteInt(value int64, lb, ub *int64, extensible bool) error {
	if lb != nil && ub != nil {
		valueRange := *ub - *lb + 1

		if extensible {
			extended := value < *lb || value > *ub
			if err := e.WriteExtensionBit(extended); err != nil {
				return err
			}
			if extended {
				return e.writeUnconstrainedInt(value)
			}
		}

		if value < *lb || value > *ub {
			return fmt.Errorf("value %d outside constrained range [%d, %d]", value, *lb, *ub)
		}

		return e.writeConstrainedWholeNumber(valueRange, uint64(value-*lb))

	} else if lb != nil {
		if value < *lb {
			return fmt.Errorf("value %d below lower bound %d", value, *lb)
		}
		return e.writeSemiConstrainedInt(*lb, value)
	}

	return e.writeUnconstrainedInt(value)
}

// WriteBool encodes BOOLEAN (1 bit)
func (e *Encoder) WriteBool(value bool) error {
	if value {
		return e.codec.Write(1, 1)
	}
	return e.codec.Write(1, 0)
}

// WriteExtensionBit writes the extension presence bit
func (e *Encoder) WriteExtensionBit(present bool) error {
	if present {
		return e.codec.Write(1, 1)
	}
	return e.codec.Write(1, 0)
}

// WriteFloat encodes REAL (IEEE 754 double)
func (e *Encoder) WriteFloat(value float64) error {
	bits := math.Float64bits(value)
	return e.codec.Write(64, bits)
}

// WriteNull encodes NULL (no encoding required)
// ITU-T X.691 Section 10.0: Null
func (e *Encoder) WriteNull() error {
	return nil
}

// WriteObjectIdentifier encodes OBJECT IDENTIFIER
// ITU-T X.691 Section 10.2: Object Identifier
// ObjectIdentifier is encoded as an OCTET STRING containing the BER encoding
func (e *Encoder) WriteObjectIdentifier(oid asn1.ObjectIdentifier) error {
	// Marshal the OID using asn1.Marshal to get BER encoding
	encoded, err := asn1.Marshal(oid)
	if err != nil {
		return fmt.Errorf("failed to marshal ObjectIdentifier: %w", err)
	}
	return e.WriteOctetString(encoded, nil, nil, false)
}

// WriteChoiceIndex encodes CHOICE index
func (e *Encoder) WriteChoiceIndex(index int, choiceCount int, extensible bool) error {
	if extensible {
		extended := index >= choiceCount
		if err := e.WriteExtensionBit(extended); err != nil {
			return err
		}
		if extended {
			return e.writeNormallySmallNumber(uint64(index))
		}
	}

	if index < 0 || index >= choiceCount {
		return fmt.Errorf("choice index %d out of root range [0, %d)", index, choiceCount)
	}
	if choiceCount <= 1 {
		return nil
	}

	return e.writeConstrainedWholeNumber(int64(choiceCount), uint64(index))
}

// WriteBitString encodes BIT STRING
func (e *Encoder) WriteBitString(bs *asn1.BitString, lb, ub *int64, extensible bool) error {
	bitLen := uint(bs.BitLength)

	if extensible && lb != nil && ub != nil {
		extended := int64(bitLen) < *lb || int64(bitLen) > *ub
		if err := e.WriteExtensionBit(extended); err != nil {
			return err
		}
		if extended {
			if err := e.writeNormallySmallNumber(uint64(bitLen)); err != nil {
				return err
			}
			if e.aligned && bitLen >= 16 {
				if err := e.Align(); err != nil {
					return err
				}
			}
			return e.writeBits(bs.Bytes, bitLen)
		}
	}

	// Length encoding
	if lb != nil && ub != nil {
		sizeRange := *ub - *lb + 1
		if int64(bitLen) < *lb || int64(bitLen) > *ub {
			return fmt.Errorf("bit string length %d outside [%d, %d]", bitLen, *lb, *ub)
		}
		if err := e.writeConstrainedWholeNumber(sizeRange, uint64(bitLen)-uint64(*lb)); err != nil {
			return err
		}
	} else if lb != nil {
		if int64(bitLen) < *lb {
			return fmt.Errorf("bit string length %d below lower bound %d", bitLen, *lb)
		}
		if e.aligned {
			if err := e.Align(); err != nil {
				return err
			}
		}
		if err := e.writeLengthDeterminant(uint64(int64(bitLen) - *lb)); err != nil {
			return err
		}
	} else {
		if e.aligned {
			if err := e.Align(); err != nil {
				return err
			}
		}
		if err := e.writeLengthDeterminant(uint64(bitLen)); err != nil {
			return err
		}
	}

	// Align contents if large in APER
	if e.aligned && bitLen >= 16 {
		if err := e.Align(); err != nil {
			return err
		}
	}

	return e.writeBits(bs.Bytes, bitLen)
}

// WriteOctetString encodes OCTET STRING
func (e *Encoder) WriteOctetString(data []byte, lb, ub *int64, extensible bool) error {
	length := uint64(len(data))

	if extensible && lb != nil && ub != nil {
		extended := int64(length) < *lb || int64(length) > *ub
		if err := e.WriteExtensionBit(extended); err != nil {
			return err
		}
		if extended {
			if err := e.writeNormallySmallNumber(length); err != nil {
				return err
			}
			if e.aligned {
				if err := e.Align(); err != nil {
					return err
				}
			}
			return e.codec.WriteBytes(data)
		}
	}

	// Length encoding
	if lb != nil && ub != nil {
		sizeRange := *ub - *lb + 1
		if int64(length) < *lb || int64(length) > *ub {
			return fmt.Errorf("octet string length %d outside [%d, %d]", length, *lb, *ub)
		}
		if err := e.writeConstrainedWholeNumber(sizeRange, length-uint64(*lb)); err != nil {
			return err
		}
	} else if lb != nil {
		if int64(length) < *lb {
			return fmt.Errorf("octet string length %d below lower bound %d", length, *lb)
		}
		if e.aligned {
			if err := e.Align(); err != nil {
				return err
			}
		}
		if err := e.writeLengthDeterminant(uint64(int64(length) - *lb)); err != nil {
			return err
		}
	} else {
		if e.aligned {
			if err := e.Align(); err != nil {
				return err
			}
		}
		if err := e.writeLengthDeterminant(length); err != nil {
			return err
		}
	}

	// Contents must be octet-aligned in APER
	if e.aligned {
		if err := e.Align(); err != nil {
			return err
		}
	}
	return e.codec.WriteBytes(data)
}

// WriteCharacterString is a convenience for IA5String, NumericString, etc.
func (e *Encoder) WriteCharacterString(value string, lb, ub *int64, extensible bool) error {
	if len(value) == 0 {
		return e.WriteOctetString(nil, lb, ub, extensible)
	}
	// Use unsafe to create a []byte view of the string without allocation
	data := unsafe.Slice(unsafe.StringData(value), len(value))
	return e.WriteOctetString(data, lb, ub, extensible)
}

// WriteEnumeration encodes ENUMERATED
func (e *Encoder) WriteEnumeration(value uint64, enumCount uint64, extensible bool) error {
	if extensible {
		extended := value >= enumCount
		if err := e.WriteExtensionBit(extended); err != nil {
			return err
		}
		if extended {
			return e.writeNormallySmallNumber(value)
		}
	}

	if value >= enumCount {
		return fmt.Errorf("enumeration value %d >= root count %d", value, enumCount)
	}
	return e.writeConstrainedWholeNumber(int64(enumCount), value)
}

// writeConstrainedWholeNumber encodes constrained non-negative integer (X.691 16)
func (e *Encoder) writeConstrainedWholeNumber(valueRange int64, value uint64) error {
	if valueRange <= 0 {
		return fmt.Errorf("invalid range %d", valueRange)
	}

	if valueRange <= 255 {
		bits := bitsRequired(uint64(valueRange - 1))
		if bits == 0 {
			bits = 1
		}
		return e.codec.Write(uint8(bits), value)
	}

	if e.aligned {
		if err := e.Align(); err != nil {
			return err
		}
	}

	if valueRange == 256 {
		return e.codec.Write(8, value)
	}
	if valueRange <= 65536 {
		return e.codec.Write(16, value)
	}

	// Large constrained: length + value (X.691 10.2.5)
	// Length is encoded as constrained whole number (range 1-8 bytes)
	nbytes := (bits.Len64(value) + 7) / 8
	if nbytes == 0 {
		nbytes = 1
	}
	// Encode length as constrained whole number with range 1-8
	if err := e.writeConstrainedWholeNumber(8, uint64(nbytes-1)); err != nil {
		return err
	}
	return e.codec.Write(uint8(nbytes*8), value)
}

// writeSemiConstrainedInt encodes semi-constrained integer
func (e *Encoder) writeSemiConstrainedInt(lb int64, value int64) error {
	if e.aligned {
		if err := e.Align(); err != nil {
			return err
		}
	}
	return e.writeLengthDeterminant(uint64(value - lb))
}

// writeUnconstrainedInt encodes unconstrained INTEGER
func (e *Encoder) writeUnconstrainedInt(value int64) error {
	if e.aligned {
		if err := e.Align(); err != nil {
			return err
		}
	}

	var nbytes uint64
	if value < 0 {
		nbytes = uint64((bits.Len64(uint64(-value-1)) + 7) / 8)
	} else {
		nbytes = uint64((bits.Len64(uint64(value)) + 7) / 8)
	}
	if nbytes == 0 {
		nbytes = 1
	}

	if err := e.writeLengthDeterminant(nbytes); err != nil {
		return err
	}
	return e.codec.Write(uint8(nbytes*8), uint64(value))
}

// writeNormallySmallNumber encodes normally small non-negative number (extension values)
func (e *Encoder) writeNormallySmallNumber(value uint64) error {
	if value < 64 {
		return e.codec.Write(7, value)
	}
	if err := e.codec.Write(1, 1); err != nil {
		return err
	}
	return e.writeLengthDeterminant(value)
}

// writeLengthDeterminant encodes length per X.691 with proper fragmentation
func (e *Encoder) writeLengthDeterminant(length uint64) error {
	if e.aligned {
		if err := e.Align(); err != nil {
			return err
		}
	}

	if length < 128 {
		return e.codec.Write(8, length)
	}
	if length < 16384 {
		return e.codec.Write(16, 0x8000|length)
	}

	remaining := length
	for remaining >= 16384 {
		if e.aligned {
			if err := e.Align(); err != nil {
				return err
			}
		}

		var indicator uint64
		if remaining >= 65536 {
			indicator = 0xC3
			remaining -= 65536
		} else if remaining >= 49152 {
			indicator = 0xC2
			remaining -= 49152
		} else if remaining >= 32768 {
			indicator = 0xC1
			remaining -= 32768
		} else {
			indicator = 0xC0
			remaining -= 16384
		}
		if err := e.codec.Write(8, indicator); err != nil {
			return err
		}
	}

	// Final fragment — align before it
	if e.aligned {
		if err := e.Align(); err != nil {
			return err
		}
	}
	if remaining < 128 {
		return e.codec.Write(8, remaining)
	}
	return e.codec.Write(16, 0x8000|remaining)
}

// writeBits writes bit string contents (without length)
func (e *Encoder) writeBits(data []byte, count uint) error {
	if count == 0 {
		return nil
	}

	num := count / 8
	if num > 0 {
		if err := e.codec.WriteBytes(data[:num]); err != nil {
			return err
		}
	}

	remaining := count % 8
	if remaining > 0 {
		var (
			last  = data[num]
			value = uint64(last >> (8 - remaining))
		)
		return e.codec.Write(uint8(remaining), value)
	}
	return nil
}

// bitsRequired calculates the minimum bits needed to represent values in range [0, maxVal]
func bitsRequired(value uint64) uint {
	if value == 0 {
		return 0 // single value → no bits needed
	}
	return uint(bits.Len64(value))
}
