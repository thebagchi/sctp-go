package per

import (
	"encoding/asn1"
	"fmt"
	"math"
	"math/bits"
	"reflect"
	"strings"
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

// WriteSequence encodes any struct as an ASN.1 SEQUENCE using reflection
func (e *Encoder) WriteSequence(value any) error {
	return e.encodeSequence(reflect.ValueOf(value))
}

// WriteSequenceOf encodes a slice or array as SEQUENCE OF
func (e *Encoder) WriteSequenceOf(elements any, lb, ub *int64, extensible bool) error {
	v := reflect.ValueOf(elements)
	if v.Kind() != reflect.Slice && v.Kind() != reflect.Array {
		return fmt.Errorf("WriteSequenceOf requires slice or array")
	}

	length := v.Len()

	// Encode length (same logic as other constrained/extensible lengths)
	if lb != nil && ub != nil {
		rangeSize := *ub - *lb + 1
		if int64(length) < *lb || int64(length) > *ub {
			return fmt.Errorf("sequence-of length %d outside [%d,%d]", length, *lb, *ub)
		}
		extended := extensible && (int64(length) < *lb || int64(length) > *ub)
		if extensible {
			if err := e.WriteExtensionBit(extended); err != nil {
				return err
			}
			if extended {
				if err := e.writeNormallySmallNumber(uint64(length)); err != nil {
					return err
				}
			} else {
				if err := e.writeConstrainedWholeNumber(rangeSize, uint64(length)-uint64(*lb)); err != nil {
					return err
				}
			}
		} else {
			if err := e.writeConstrainedWholeNumber(rangeSize, uint64(length)-uint64(*lb)); err != nil {
				return err
			}
		}
	} else {
		// Unconstrained
		if e.aligned {
			if err := e.Align(); err != nil {
				return err
			}
		}
		if err := e.writeLengthDeterminant(uint64(length)); err != nil {
			return err
		}
	}

	// Encode each element
	for i := 0; i < length; i++ {
		elem := v.Index(i)
		if err := e.encodeValue(elem.Interface()); err != nil {
			return err
		}
	}
	return nil
}

// WriteChoice encodes a CHOICE using the _ struct{} `per:"ext"` marker pattern
// Extension alternatives are encoded as open types (length + raw PER bits)
func (e *Encoder) WriteChoice(value any) error {
	v := reflect.ValueOf(value)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return fmt.Errorf("WriteChoice requires a struct")
	}

	t := v.Type()
	numFields := t.NumField()

	// Count root alternatives and detect extensibility
	rootCount := 0
	extensible := false
	for i := 0; i < numFields; i++ {
		field := t.Field(i)
		if field.Name == "_" && field.Type.Kind() == reflect.Struct && field.Type.NumField() == 0 {
			tag := field.Tag.Get("per")
			if strings.Contains(tag, "ext") {
				extensible = true
				// Root count is number of choice= fields before this marker
				break
			}
		}
		if tag := field.Tag.Get("per"); strings.Contains(tag, "choice=") {
			rootCount++
		}
	}

	// Find selected alternative
	var selectedIndex int = -1
	var selectedValue any

	for i := 0; i < numFields; i++ {
		fv := v.Field(i)
		if isZero(fv) {
			continue
		}

		field := t.Field(i)
		opts := GetFieldTag(field, t, i)

		if opts.choice == nil {
			continue
		}

		if selectedIndex != -1 {
			return fmt.Errorf("multiple choice alternatives set")
		}

		selectedIndex = *opts.choice
		selectedValue = fv.Interface()
	}

	if selectedIndex == -1 {
		return fmt.Errorf("no choice alternative selected")
	}

	// Determine if this is an extension
	extended := extensible && selectedIndex >= rootCount

	// Encode using low-level WriteChoiceIndex (handles extension bit and normally small index)
	if err := e.WriteChoiceIndex(selectedIndex, rootCount, extensible); err != nil {
		return err
	}

	// If it's an extension alternative, encode value as open type
	if extended {
		// Encode the value into a temporary encoder to get its raw PER bits
		tempEnc := NewEncoder(e.aligned)
		if err := tempEnc.encodeValue(selectedValue); err != nil {
			return err
		}
		extensionBits := tempEnc.Bytes()

		// Align before open type (X.691 11.5.3)
		if e.aligned {
			if err := e.Align(); err != nil {
				return err
			}
		}

		// Open type = unconstrained octet string
		return e.WriteOctetString(extensionBits, nil, nil, false)
	}

	// Root alternative - encode value directly
	return e.encodeValue(selectedValue)
}

// Helper: encodeValue dispatches to correct method for known kinds
// Note: SEQUENCE OF (slices/arrays) must be encoded via encodeField with proper field tags
func (e *Encoder) encodeValue(value any) error {
	if value == nil {
		return nil
	}

	rv := reflect.ValueOf(value)
	if rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return nil
		}
		rv = rv.Elem()
	}

	switch rv.Kind() {
	case reflect.Struct:
		return e.WriteSequence(value)
	case reflect.Slice, reflect.Array:
		// SEQUENCE OF must be encoded via encodeField with proper field tags
		// that specify constraints (lb, ub, extensible)
		return fmt.Errorf("SEQUENCE OF must be encoded via tagged field, not direct value")
	default:
		return fmt.Errorf("unsupported value type %v in encodeValue", rv.Kind())
	}
}

// encodeSequence encodes a struct as SEQUENCE
// Per X.691: optional field preamble (bitmap) precedes all field encodings
// Extended fields are encoded as open types with extension bitmap
func (e *Encoder) encodeSequence(v reflect.Value) error {
	if v.Kind() == reflect.Pointer {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return fmt.Errorf("encodeSequence requires struct")
	}

	t := v.Type()

	// Find extension marker
	var extensionStart int = -1
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.Name == "_" && field.Type.Kind() == reflect.Struct && field.Type.NumField() == 0 {
			tag := field.Tag.Get("per")
			if strings.Contains(tag, "ext") {
				extensionStart = i
				break
			}
		}
	}

	// First pass: collect optional fields in root and build presence bitmap (X.691 11.5.1)
	var rootOptionalIndices []int
	var rootOptionalPresent uint64

	for i := 0; i < t.NumField(); i++ {
		if extensionStart != -1 && i >= extensionStart {
			break // Skip extension fields in this pass
		}

		field := t.Field(i)
		if !v.Field(i).CanInterface() {
			continue
		}

		opts := GetFieldTag(field, t, i)

		// Pointer fields are automatically optional, even without opt tag
		isOptional := opts.opt || field.Type.Kind() == reflect.Pointer
		if isOptional {
			rootOptionalIndices = append(rootOptionalIndices, i)
			if !isZero(v.Field(i)) {
				bitPos := len(rootOptionalIndices) - 1
				rootOptionalPresent |= (1 << (len(rootOptionalIndices) - 1 - bitPos))
			}
		}
	}

	// Encode root optional presence bitmap upfront (X.691 SEQUENCE preamble)
	// Per X.691 11.5.1: N bits where N = number of optional fields
	if len(rootOptionalIndices) > 0 {
		// Build bitmap bytes (MSB = first optional field, LSB = last)
		bitmapBytes := make([]byte, (len(rootOptionalIndices)+7)/8)
		for i := 0; i < len(rootOptionalIndices); i++ {
			byteIdx := i / 8
			bitIdx := uint(7 - (i % 8))
			if (rootOptionalPresent & (1 << (len(rootOptionalIndices) - 1 - i))) != 0 {
				bitmapBytes[byteIdx] |= (1 << bitIdx)
			}
		}
		// Write bitmap bits directly (no length encoding for preamble)
		if err := e.writeBits(bitmapBytes, uint(len(rootOptionalIndices))); err != nil {
			return err
		}
	}

	// Encode extension bit if extensible
	var hasExtensions bool
	if extensionStart != -1 {
		// Check if any extension fields are present
		for i := extensionStart + 1; i < t.NumField(); i++ {
			if v.Field(i).CanInterface() && !isZero(v.Field(i)) {
				hasExtensions = true
				break
			}
		}
		if err := e.WriteExtensionBit(hasExtensions); err != nil {
			return err
		}
	}

	// Second pass: encode all root fields (skip absent optional fields)
	rootOptIdx := 0
	for i := 0; i < t.NumField(); i++ {
		if extensionStart != -1 && i >= extensionStart {
			break // Skip extension fields in this pass
		}

		field := t.Field(i)
		fv := v.Field(i)

		// Skip unexported fields
		if !fv.CanInterface() {
			continue
		}

		opts := GetFieldTag(field, t, i)

		// Pointer fields are automatically optional
		isOptional := opts.opt || field.Type.Kind() == reflect.Ptr

		// Skip absent optional fields
		if isOptional {
			if rootOptIdx < len(rootOptionalIndices) && rootOptionalIndices[rootOptIdx] == i {
				if isZero(fv) {
					rootOptIdx++
					continue
				}
				rootOptIdx++
			}
		}

		if err := e.encodeField(field, fv); err != nil {
			return err
		}
	}

	// Third pass: encode extension fields as open types with extension bitmap (X.691 11.5.3)
	if hasExtensions && extensionStart != -1 {
		// Collect extension fields
		var extensionIndices []int
		var extensionPresent uint64

		for i := extensionStart + 1; i < t.NumField(); i++ {
			field := t.Field(i)
			fv := v.Field(i)

			if !fv.CanInterface() {
				continue
			}

			tag := field.Tag.Get("per")
			if tag == "" {
				continue
			}

			extensionIndices = append(extensionIndices, i)
			if !isZero(fv) {
				idx := len(extensionIndices) - 1
				bitPos := len(extensionIndices) - 1 - idx
				extensionPresent |= (1 << bitPos)
			}
		}

		// Encode extension bitmap
		if len(extensionIndices) > 0 {
			// Build bitmap bytes (MSB = first extension field, LSB = last)
			bitmapBytes := make([]byte, (len(extensionIndices)+7)/8)
			for i := 0; i < len(extensionIndices); i++ {
				byteIdx := i / 8
				bitIdx := uint(7 - (i % 8))
				if (extensionPresent & (1 << (len(extensionIndices) - 1 - i))) != 0 {
					bitmapBytes[byteIdx] |= (1 << bitIdx)
				}
			}
			// Write bitmap bits directly (no length encoding for preamble)
			if err := e.writeBits(bitmapBytes, uint(len(extensionIndices))); err != nil {
				return err
			}

			// Encode each present extension as open type
			for idx, fieldIdx := range extensionIndices {
				bitPos := len(extensionIndices) - 1 - idx
				if (extensionPresent & (1 << bitPos)) != 0 {
					fv := v.Field(fieldIdx)
					field := t.Field(fieldIdx)

					// Encode value as open type (temporary encoder)
					tempEnc := NewEncoder(e.aligned)
					if err := tempEnc.encodeField(field, fv); err != nil {
						return err
					}
					extensionBits := tempEnc.Bytes()

					// Align before open type (X.691 11.5.3)
					if e.aligned {
						if err := e.Align(); err != nil {
							return err
						}
					}

					// Open type = unconstrained octet string
					if err := e.WriteOctetString(extensionBits, nil, nil, false); err != nil {
						return err
					}
				}
			}
		}
	}

	return nil
}

// encodeField encodes a single field based on its type and tag constraints
func (e *Encoder) encodeField(field reflect.StructField, value reflect.Value) error {
	t := field.Type
	v := value

	if v.Kind() == reflect.Pointer {
		if v.IsNil() {
			return nil
		}
		v = v.Elem()
		t = t.Elem()
	}

	// Handle Null type
	if t.Name() == "Null" && t.Kind() == reflect.Struct {
		return e.WriteNull()
	}

	// Parse field tags for constraints (using cache)
	opts := GetFieldTag(field, field.Type, 0)

	switch t.Kind() {
	case reflect.Bool:
		return e.WriteBool(v.Bool())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return e.WriteInt(v.Int(), nil, nil, false)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return e.WriteInt(int64(v.Uint()), nil, nil, false)
	case reflect.Float32, reflect.Float64:
		return e.WriteFloat(v.Float())
	case reflect.String:
		return e.WriteCharacterString(v.String(), nil, nil, false)
	case reflect.Slice:
		if t.Elem().Kind() == reflect.Uint8 {
			// []byte - OCTET STRING with field tag constraints
			return e.WriteOctetString(v.Bytes(), opts.lb, opts.ub, opts.ext)
		}
		// Other slices as SEQUENCE OF with field tag constraints
		return e.WriteSequenceOf(v.Interface(), opts.lb, opts.ub, opts.ext)
	case reflect.Array:
		if t.Elem().Kind() == reflect.Uint8 {
			// [n]byte - OCTET STRING with field tag constraints
			return e.WriteOctetString(v.Slice(0, v.Len()).Bytes(), opts.lb, opts.ub, opts.ext)
		}
		// Other arrays as SEQUENCE OF with field tag constraints
		return e.WriteSequenceOf(v.Interface(), opts.lb, opts.ub, opts.ext)
	case reflect.Struct:
		return e.WriteSequence(v.Interface())
	default:
		return fmt.Errorf("unsupported type: %s", t.String())
	}
}

// isZero reports whether v is the zero value for its type.
// It handles all common kinds: primitives, pointers, slices, maps, structs, etc.
func isZero(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Complex64, reflect.Complex128:
		return v.Complex() == 0
	case reflect.String:
		return v.Len() == 0
	case reflect.Array:
		return v.Len() == 0 // empty array is considered zero
	case reflect.Slice, reflect.Map, reflect.Chan, reflect.Func, reflect.Interface, reflect.Ptr:
		return v.IsNil()
	case reflect.Struct:
		// Recursively check all fields
		for i := 0; i < v.NumField(); i++ {
			if !isZero(v.Field(i)) {
				return false
			}
		}
		return true
	case reflect.UnsafePointer:
		return v.Pointer() == 0
	default:
		// For any unsupported kind, fall back to reflect's IsZero (available from Go 1.13+)
		return v.IsZero()
	}
}

// bitsRequired calculates the minimum bits needed to represent values in range [0, maxVal]
func bitsRequired(value uint64) uint {
	if value == 0 {
		return 0 // single value → no bits needed
	}
	return uint(bits.Len64(value))
}
